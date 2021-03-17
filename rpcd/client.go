package rpcd

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/imkouga/gocore/loger"
)

const (
	doneChannelUnbufferedMsg = "rpcd.Client.Go: done channel is unbuffered"
)

// ServerError represents an error that has been returned from
// the remote side of the RPC connection.
type ServerError string

func (e ServerError) Error() string {
	return string(e)
}

var ErrShutdown = errors.New("connection is shut down")

type Call struct {
	ServiceMethod string      // The name of the service and method to call.
	Args          interface{} // The argument to the function (*struct).
	Reply         interface{} // The reply from the function (*struct).
	Error         error       // After completion, the error status.
	Done          chan *Call  // Receives *Call when Go is complete.
}

type Client struct {
	codec ClientCodec

	reqMutex sync.Mutex // protects following
	request  Request

	mutex    sync.Mutex // protects following
	seq      uint64
	pending  map[uint64]*Call
	closed   bool // user has called Close
	shutdown bool // server has told us to stop
}

type ClientCodec interface {
	WriteRequest(*Request, interface{}) error
	ReadResponseHeader(*Response) error
	ReadResponseBody(interface{}) error

	Close() error
}

func Dial(network, address string, opts ...DialOption) (*Client, error) {

	dialOpts := &dialOptions{
		timeout: time.Second * 30,
	}
	for _, opt := range opts {
		opt.apply(dialOpts)
	}

	conn, err := net.DialTimeout(network, address, dialOpts.timeout)
	if err != nil {
		return nil, err
	}
	codec := NewPbClientCodec(conn)
	return NewClientWithCodec(codec), nil
}

func NewClientWithCodec(codec ClientCodec) *Client {

	client := &Client{
		codec:   codec,
		pending: make(map[uint64]*Call),
	}
	go client.input()
	return client
}

func (client *Client) send(call *Call) {

	client.reqMutex.Lock()
	defer client.reqMutex.Unlock()

	if client.ClosedOrShutdowned() {
		call.Error = ErrShutdown
		call.done()
		return
	}

	// Register this call.
	client.mutex.Lock()
	seq := client.seq
	client.seq++
	client.pending[seq] = call
	client.mutex.Unlock()

	// Encode and send the request.
	client.request.Seq = seq
	client.request.ServiceMethod = call.ServiceMethod
	err := client.codec.WriteRequest(&client.request, call.Args)
	if err != nil {
		client.mutex.Lock()
		call = client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (client *Client) input() {

	var err error
	var response Response
	for err == nil {
		response = Response{}
		err = client.codec.ReadResponseHeader(&response)
		if err != nil {
			break
		}
		seq := response.Seq
		client.mutex.Lock()
		call := client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()

		switch {
		case call == nil:
			err = client.codec.ReadResponseBody(nil)
			if err != nil {
				err = fmt.Errorf("rpcd.Client.input: reading error body[%s]", err.Error())
			}
		case response.Error != "":
			call.Error = ServerError(response.Error)
			err = client.codec.ReadResponseBody(nil)
			if err != nil {
				err = fmt.Errorf("rpcd.Client.input: reading error body[%s]", err.Error())
			}
			call.done()
		default:
			err = client.codec.ReadResponseBody(call.Reply)
			if err != nil {
				call.Error = fmt.Errorf("rpcd.Client.input: reading body[%s]", err.Error())
			}
			call.done()
		}
	}
	// Terminate pending calls.
	client.reqMutex.Lock()
	client.mutex.Lock()
	client.shutdown = true
	closed := client.closed
	if err == io.EOF {
		if err = ErrShutdown; !closed {
			err = io.ErrUnexpectedEOF
		}
	}
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
	client.mutex.Unlock()
	client.reqMutex.Unlock()
	if err != io.EOF && !closed {
		loger.Error("rpcd.Client.input: client protocol error:", err)
	}
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		loger.Warn("rpcd.Call.done: discarding Call reply due to insufficient Done chan capacity")
	}
}

func (client *Client) Close() error {

	client.mutex.Lock()
	if client.closed {
		client.mutex.Unlock()
		return ErrShutdown
	}
	client.closed = true
	client.mutex.Unlock()
	return client.codec.Close()
}

func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {

	call := new(Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply

	if done == nil {
		done = make(chan *Call, 10) // buffered.
	} else if cap(done) <= 0 {
		call.Done = done
		call.Error = errors.New(doneChannelUnbufferedMsg)
		loger.Error(doneChannelUnbufferedMsg)
		return call
	}

	call.Done = done
	client.send(call)
	return call
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

func (client *Client) ClosedOrShutdowned() bool {

	if nil == client {
		return true
	}
	client.mutex.Lock()
	defer client.mutex.Unlock()

	return client.closed || client.shutdown
}
