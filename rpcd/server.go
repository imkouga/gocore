package rpcd

import (
	"errors"
	"fmt"
	"go/token"
	"io"
	"net"
	"reflect"
	"strings"
	"sync"

	"github.com/imkouga/gocore/loger"
)

// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var (
	typeOfError    = reflect.TypeOf((*error)(nil)).Elem()
	invalidRequest = struct{}{}
)

type ServerCodec interface {
	ReadRequestHeader(*Request) error
	ReadRequestBody(interface{}) error
	WriteResponse(*Response, interface{}) error

	// Close can be called multiple times and must be idempotent.
	Close() error
}

type methodType struct {
	sync.Mutex // protects counters
	method     reflect.Method
	ArgType    reflect.Type
	ReplyType  reflect.Type
	numCalls   uint
}

type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

// Request is a header written before every RPC call. It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Request struct {
	ServiceMethod string   // format: "Service.Method"
	Seq           uint64   // sequence number chosen by client
	next          *Request // for free list in Server
}

// Response is a header written before every RPC return. It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Response struct {
	ServiceMethod string    // echoes that of the Request
	Seq           uint64    // echoes that of the request
	Error         string    // error, if any.
	next          *Response // for free list in Server
}

// HandleHTTP registers an HTTP handler for RPC messages to DefaultServer
// on DefaultRPCPath and a debugging handler on DefaultDebugPath.
// It is still necessary to invoke http.Serve(), typically in a go statement.

type RPCServer struct {
	serviceMap sync.Map   // map[string]*service
	reqLock    sync.Mutex // protects freeReq
	freeReq    *Request
	respLock   sync.Mutex // protects freeResp
	freeResp   *Response
}

func NewRPCServer() *RPCServer {
	return &RPCServer{}
}

func (rs *RPCServer) ServeWithPbCodec(addr string) error {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				loger.Errorf("rpcd.RPCServer.ServeWithPbCodec: Accept failed. err[%s]", err.Error())
			}

			codec := NewPbServerCode(conn)

			go rs.ServeCodec(codec)
		}
	}()

	return nil
}

// ServeCodec is like ServeConn but uses the specified codec to
// decode requests and encode responses.
func (server *RPCServer) ServeCodec(codec ServerCodec) {

	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		service, mtype, req, argv, replyv, keepReading, err := server.readRequest(codec)
		if err != nil {
			if err != io.EOF {
				loger.Errorf("rpcd.RPCServer.ServeCode: rpc err[%s]", err.Error())
			}
			if !keepReading {
				break
			}
			// send a response if we actually managed to read a header.
			if req != nil {
				server.sendResponse(sending, req, invalidRequest, codec, err.Error())
				server.freeRequest(req)
			}
			continue
		}
		wg.Add(1)
		go service.call(server, sending, wg, mtype, req, argv, replyv, codec)
	}
	// We've seen that there are no more requests.
	// Wait for responses to be sent before closing codec.
	wg.Wait()
	codec.Close()
}

// ServeRequest is like ServeCodec but synchronously serves a single request.
// It does not close the codec upon completion.
func (server *RPCServer) ServeRequest(codec ServerCodec) error {

	sending := new(sync.Mutex)
	service, mtype, req, argv, replyv, keepReading, err := server.readRequest(codec)
	if err != nil {
		if !keepReading {
			return err
		}
		// send a response if we actually managed to read a header.
		if req != nil {
			server.sendResponse(sending, req, invalidRequest, codec, err.Error())
			server.freeRequest(req)
		}
		return err
	}
	service.call(server, sending, nil, mtype, req, argv, replyv, codec)
	return nil
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (server *RPCServer) RegisterName(name string, rcvr interface{}, routes map[string]string) error {
	return server.register(rcvr, name, routes)
}

func (server *RPCServer) register(rcvr interface{}, sname string, routes map[string]string) error {

	var err error

	s := new(service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	if sname == "" {
		s := "rpcd.RPCServer.register: no service name for type " + s.typ.String()
		loger.Error(s)
		return errors.New(s)
	}
	s.name = sname

	// Install the methods
	s.method, err = suitableMethods(s.typ, routes)
	if nil != err {
		return err
	}
	if len(s.method) <= 0 {
		method, _ := suitableMethods(reflect.PtrTo(s.typ), routes)
		str := "rpcd.RPCServer.register: type " + sname + " has no exported methods of suitable type"
		if len(method) != 0 {
			str = "rpcd.RPCServer.register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		}
		return errors.New(str)
	}

	if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
		return errors.New("rpcd.RPCServer.register: service already defined: " + sname)
	}
	return nil
}

// A value sent as a placeholder for the server's response value when the server
// receives an invalid request. It is never decoded by the client since the Response
// contains an error when it is used.
func (server *RPCServer) sendResponse(sending *sync.Mutex, req *Request, reply interface{}, codec ServerCodec, errmsg string) {

	resp := server.getResponse()
	// Encode the response header
	resp.ServiceMethod = req.ServiceMethod
	if errmsg != "" {
		resp.Error = errmsg
		reply = invalidRequest
	}
	resp.Seq = req.Seq
	sending.Lock()
	err := codec.WriteResponse(resp, reply)
	if err != nil {
		loger.Errorf("rpcd.RPCServer.sendResponse: writing response err[%s]", err.Error())
	}
	sending.Unlock()
	server.freeResponse(resp)
}

func (s *service) call(server *RPCServer, sending *sync.Mutex, wg *sync.WaitGroup, mtype *methodType, req *Request, argv, replyv reflect.Value, codec ServerCodec) {

	if wg != nil {
		defer wg.Done()
	}
	mtype.Lock()
	mtype.numCalls++
	mtype.Unlock()
	function := mtype.method.Func
	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{s.rcvr, argv, replyv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	errmsg := ""
	if errInter != nil {
		errmsg = errInter.(error).Error()
	}
	server.sendResponse(sending, req, replyv.Interface(), codec, errmsg)
	server.freeRequest(req)
}

func (server *RPCServer) getRequest() *Request {

	server.reqLock.Lock()
	req := server.freeReq
	if req == nil {
		req = new(Request)
	} else {
		server.freeReq = req.next
		*req = Request{}
	}
	server.reqLock.Unlock()
	return req
}

func (server *RPCServer) freeRequest(req *Request) {
	server.reqLock.Lock()
	req.next = server.freeReq
	server.freeReq = req
	server.reqLock.Unlock()
}

func (server *RPCServer) getResponse() *Response {

	server.respLock.Lock()

	resp := server.freeResp
	if resp == nil {
		resp = new(Response)
		server.respLock.Unlock()
		return resp
	}

	server.freeResp = resp.next
	*resp = Response{}
	server.respLock.Unlock()
	return resp
}

func (server *RPCServer) freeResponse(resp *Response) {
	server.respLock.Lock()
	resp.next = server.freeResp
	server.freeResp = resp
	server.respLock.Unlock()
}

func (server *RPCServer) readRequest(codec ServerCodec) (service *service, mtype *methodType, req *Request, argv, replyv reflect.Value, keepReading bool, err error) {

	service, mtype, req, keepReading, err = server.readRequestHeader(codec)
	if err != nil {
		if !keepReading {
			return
		}
		// discard body
		codec.ReadRequestBody(nil)
		return
	}

	// Decode the argument value.
	argIsValue := false // if true, need to indirect before calling.
	if mtype.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(mtype.ArgType.Elem())
	} else {
		argv = reflect.New(mtype.ArgType)
		argIsValue = true
	}
	// argv guaranteed to be a pointer now.
	if err = codec.ReadRequestBody(argv.Interface()); err != nil {
		return
	}
	if argIsValue {
		argv = argv.Elem()
	}

	replyv = reflect.New(mtype.ReplyType.Elem())

	switch mtype.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(mtype.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(mtype.ReplyType.Elem(), 0, 0))
	}
	return
}

func (server *RPCServer) readRequestHeader(codec ServerCodec) (svc *service, mtype *methodType, req *Request, keepReading bool, err error) {

	// Grab the request header.
	req = server.getRequest()
	err = codec.ReadRequestHeader(req)
	if err != nil {
		req = nil
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		err = fmt.Errorf("rpcd.RPCServer.readRequestHeader: server cannot decode request. err[%s]", err.Error())
		return
	}

	// We read the header successfully. If we see an error now,
	// we can still recover and move on to the next request.
	keepReading = true

	dot := strings.Index(req.ServiceMethod, ".")
	if dot < 0 {
		err = fmt.Errorf("rpcd.RPCServer.readRequestHeader: service/method request ill-formed. method[%s]", req.ServiceMethod)
		return
	}
	serviceName := req.ServiceMethod[:dot]
	methodName := req.ServiceMethod[dot+1:]

	// Look up the request.
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = fmt.Errorf("rpcd.RPCServer.readRequestHeader: can't find service. method[%s]", req.ServiceMethod)
		return
	}
	svc = svci.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = fmt.Errorf("rpcd.RPCServer.readRequestHeader: can't find method[%s]", req.ServiceMethod)
	}
	return
}

func (m *methodType) NumCalls() (n uint) {
	m.Lock()
	n = m.numCalls
	m.Unlock()
	return n
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return token.IsExported(t.Name()) || t.PkgPath() == ""
}

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using log if reportErr is true.
func suitableMethods(typ reflect.Type, routes map[string]string) (map[string]*methodType, error) {

	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() != 3 {
			loger.Errorf("rpc.suitableMethods: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			loger.Errorf("rpc.suitableMethods: argument type of method %q is not exported: %q\n", mname, argType)
			continue
		}
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			loger.Errorf("rpc.suitableMethods: reply type of method %q is not a pointer: %q\n", mname, replyType)
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			loger.Errorf("rpc.suitableMethods: reply type of method %q is not exported: %q\n", mname, replyType)
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			loger.Errorf("rpc.suitableMethods: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			loger.Errorf("rpc.suitableMethods: return type of method %q is %q, must be error\n", mname, returnType)
			continue
		}
		if _, exist := routes[mname]; exist {
			mname = routes[mname]
		}
		methods[mname] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
	}
	return methods, nil
}
