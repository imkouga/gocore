package rpcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"sync"

	"github.com/imkouga/gocore/loger"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type pbServerCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	m                *sync.Mutex
	req              uint64
	currentReqHeader *RequestHeader
	handles          map[uint64]*RequestHeader
}

func NewPbServerCode(conn io.ReadWriteCloser) *pbServerCodec {

	codec := &pbServerCodec{
		r:       conn,
		w:       conn,
		c:       conn,
		m:       new(sync.Mutex),
		handles: make(map[uint64]*RequestHeader),
	}
	return codec
}

func (codec *pbServerCodec) ReadRequestHeader(req *Request) error {

	var (
		rh  RequestHeader
		rhb []byte
		err error
	)

	if rhb, err = recvFrame(codec.r); nil != err {
		if io.EOF != err {
			loger.Errorf("rpcd.pbClientCodec.ReadRequestHeader: read data failed. err[%s]", err.Error())
		}
		return err
	}

	loger.Tracef("rpcd.pbClientCodec.ReadRequestHeader: header len[%d]", len(rhb))
	if err = proto.Unmarshal(rhb, &rh); nil != err {
		loger.Errorf("rpcd.pbServerCodec.ReadRequestHeader: proto unmarshal data failed. err[%s]", err.Error())
		return err
	}
	req.Seq = codec.handling(&rh)
	req.ServiceMethod = rh.Method
	codec.currentReqHeader = &rh

	return nil
}

func (codec *pbServerCodec) ReadRequestBody(x interface{}) error {

	var (
		body proto.Message
		rhb  []byte
		err  error
		ok   bool
	)
	if nil == codec.currentReqHeader {
		loger.Errorf("rpcd.pbServerCodec.ReadRequestBody: not exist req header.")
		return errors.New("not exist req header")
	}

	if rhb, err = recvFrame(codec.r); nil != err {
		if io.EOF != err {
			loger.Errorf("rpcd.pbClientCodec.ReadResponseBody: read data failed. err[%s]", err.Error())
		}
		return err
	}
	loger.Tracef("rpcd.pbClientCodec.ReadRequestBody: request body len[%d]", len(rhb))

	if nil != x {
		if body, ok = x.(proto.Message); !ok {
			err := fmt.Errorf("rpcd.pbClientCodec.ReadResponseBody: %T does not implement proto.Message", x)
			loger.Error(err)
			return err
		}
		if err = codec.currentReqHeader.Check(rhb); nil != err {
			loger.Errorf("rpcd.pbServerCodec.ReadRequestBody: err[%s]", err.Error())
			return err
		}

		if RpcdSourceGrowerGateway == codec.currentReqHeader.RpcdSource {
			loger.Tracef("rpcd.pbServerCodec.ReadRequestBody: sugar request, method[%s], source[%s]", codec.currentReqHeader.Method, codec.currentReqHeader.RpcdSource)
			sugarCoating := new(SugarCoatingRequest)
			if err = proto.Unmarshal(rhb, sugarCoating); nil != err {
				loger.Errorf("rpcd.pbServerCodec.ReadRequestBody: proto unmarshal data failed. err[%s]", err.Error())
				return err
			}
			if err = jsonpb.UnmarshalString(string(sugarCoating.Data), body); nil != err {
				loger.Errorf("rpcd.pbServerCodec.ReadRequestBody: jsonpb unmarshal data failed. err[%s]", err.Error())
				return err
			}
			return nil
		}

		if err = proto.Unmarshal(rhb, body); nil != err {
			loger.Errorf("rpcd.pbServerCodec.ReadRequestBody: proto unmarshal data failed. err[%s]", err.Error())
			return err
		}
	}

	return nil
}

func (codec *pbServerCodec) WriteResponse(resp *Response, x interface{}) error {

	var (
		body []byte
		err  error
	)
	req, err := codec.handled(resp)
	if nil != err {
		return err
	}
	if len(resp.Error) <= 0 {
		if body, err = codec.transferResponseBody(req, x); nil != err {
			loger.Errorf("rpcd.pbServerCodec.WriteResponse: %s", err.Error())
			return err
		}
	}

	if err := codec.writeResponseHeader(req, resp, body); nil != err {
		loger.Errorf("rpcd.pbServerCodec.WriteResponse: %s", err.Error())
		return err
	}
	if err := codec.writeResponseBody(body); nil != err {
		loger.Errorf("rpcd.pbServerCodec.WriteResponse: %s", err.Error())
		return err
	}

	loger.Tracef("write response finish. seq[%d] method[%s]", resp.Seq, resp.ServiceMethod)
	return nil
}

func (codec *pbServerCodec) writeResponseHeader(req *RequestHeader, resp *Response, x []byte) error {

	respHeader := &ResponseHeader{
		Id:       req.Id,
		Method:   resp.ServiceMethod,
		Error:    resp.Error,
		Checksum: crc32.ChecksumIEEE(x),
	}

	headers, err := proto.Marshal(respHeader)
	if nil != err {
		return fmt.Errorf("proto marshal data failed. err[%s]", err.Error())
	}

	if err := sendFrame(codec.w, headers); nil != err {
		return fmt.Errorf("write data to I/O failed. err[%s]", err.Error())
	}
	return nil
}

func (codec *pbServerCodec) writeResponseBody(body []byte) error {
	if err := sendFrame(codec.w, body); nil != err {
		return fmt.Errorf("write data to I/O failed. err[%s]", err.Error())
	}
	return nil
}

func (codec *pbServerCodec) transferResponseBody(req *RequestHeader, x interface{}) ([]byte, error) {

	var (
		sugarData []byte
		err       error
	)
	if req.RpcdSource == RpcdSourceGrowerGateway {
		sugarCoatingResponse := new(SugarCoatingResponse)
		if s, ok := x.(SugarCoatingResponseInterface); ok {
			if sugarData, err = s.GetSugarData(); nil != err {
				return nil, err
			}
		} else if sugarData, err = json.Marshal(x); nil != err {
			return nil, err
		}
		sugarCoatingResponse.Data = sugarData

		x = sugarCoatingResponse
	}

	if resp, ok := x.(proto.Message); ok {
		return proto.Marshal(resp)
	}
	return nil, fmt.Errorf("%T does not implement proto.Message", x)
}

// Close can be called multiple times and must be idempotent.
func (codec *pbServerCodec) Close() error {
	return codec.c.Close()
}

func (codec *pbServerCodec) handling(rh *RequestHeader) uint64 {

	codec.m.Lock()
	defer codec.m.Unlock()

	codec.req++
	codec.handles[codec.req] = rh
	return codec.req
}

func (codec *pbServerCodec) handled(resp *Response) (*RequestHeader, error) {

	codec.m.Lock()
	defer codec.m.Unlock()

	if reqHeader, ok := codec.handles[resp.Seq]; ok {
		delete(codec.handles, resp.Seq)
		return reqHeader, nil
	}
	return nil, errors.New("not exist resp header.")
}
