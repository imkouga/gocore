package rpcd

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"sync"

	"github.com/imkouga/gocore/loger"

	"github.com/golang/protobuf/proto"
)

type pbClientCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	m                 *sync.Mutex
	currentRespHeader *ResponseHeader
	pends             map[uint64]*ResponseHeader
}

func NewPbClientCodec(conn io.ReadWriteCloser) *pbClientCodec {

	codec := &pbClientCodec{
		r:     conn,
		w:     conn,
		c:     conn,
		m:     new(sync.Mutex),
		pends: make(map[uint64]*ResponseHeader),
	}
	return codec
}

func (codec *pbClientCodec) WriteRequest(req *Request, data interface{}) error {

	pbRequest, err := codec.transferRequestBody(data)
	if nil != err {
		loger.Errorf("rpcd.pbClientCodec.WriteRequest: %s", err.Error())
		return err
	}

	requestHeader := &RequestHeader{
		Id:       req.Seq,
		Method:   req.ServiceMethod,
		Checksum: crc32.ChecksumIEEE(pbRequest),
	}
	if x, ok := data.(*SugarCoatingRequest); ok {
		requestHeader.RpcdSource = x.RpcdSource
	}
	if err := codec.writeRequestHeader(requestHeader); nil != err {
		loger.Errorf("rpcd.pbClientCodec.WriteRequest.writeRequestHeader: %s", err.Error())
		return err
	}

	if err := codec.writeRequestBody(pbRequest); nil != err {
		loger.Errorf("rpcd.pbClientCodec.WriteRequest.writeRequestBody: %s", err.Error())
		return err
	}

	loger.Tracef("rpcd.pbClientCodec.WriteRequest: write request finish. req[%d], method[%s]", req.Seq, req.ServiceMethod)
	return nil
}

func (codec *pbClientCodec) ReadResponseHeader(resp *Response) error {

	var (
		respBody []byte
		err      error
	)

	if respBody, err = recvFrame(codec.r); nil != err {
		if io.EOF != err {
			loger.Errorf("rpcd.pbClientCodec.ReadResponseHeader: read data failed. err[%s]", err.Error())
		}
		return err
	}
	if len(respBody) <= 0 {
		loger.Trace("rpcd.pbClientCodec.ReadResponseHeader: read resp header len is 0")
		return nil
	}

	responseHeader := new(ResponseHeader)
	if err = proto.Unmarshal(respBody, responseHeader); nil != err {
		loger.Errorf("rpcd.pbClientCodec.ReadResponseHeader: proto unmarshal data failed. err[%s]", err.Error())
		return err
	}
	if codec.currentRespHeader, err = codec.pended(responseHeader); nil != err {
		loger.Errorf("rpcd.pbClientCodec.ReadResponseHeader: %s", err.Error())
		return err
	}

	resp.Seq, resp.Error, resp.ServiceMethod = codec.currentRespHeader.Id, codec.currentRespHeader.Error, codec.currentRespHeader.Method
	return nil
}

func (codec *pbClientCodec) ReadResponseBody(x interface{}) error {

	if nil == codec.currentRespHeader {
		loger.Errorf("rpcd.pbClientCodec.ReadResponseBody: not exist resp header.")
		return errors.New("not exist resp header")
	}
	respBody, err := recvFrame(codec.r)
	if nil != err {
		if io.EOF != err {
			loger.Errorf("rpcd.pbClientCodec.ReadResponseBody: read data failed. err[%s]", err.Error())
		}
		return err
	}

	if nil != x {
		if err := codec.currentRespHeader.Check(respBody); nil != err {
			loger.Errorf("rpcd.pbClientCodec.ReadResponseBody: err[%s]", err.Error())
			return err
		}
		if respMessage, ok := x.(proto.Message); ok {
			if err := proto.Unmarshal(respBody, respMessage); nil != err {
				loger.Errorf("rpcd.pbClientCodec.ReadResponseBody: proto unmarshal data failed. err[%s]", err.Error())
				return err
			}
			return nil
		}
		err := fmt.Errorf("rpcd.pbClientCodec.ReadResponseBody: %T does not implement proto.Message", x)
		loger.Error(err)
		return err
	}

	return nil
}

func (codec *pbClientCodec) Close() error {
	return codec.c.Close()
}

func (codec *pbClientCodec) writeRequestHeader(req *RequestHeader) error {

	codec.pending(req)

	loger.Trace("rpcd.pbClientCodec.writeRequestHeader:", req)
	headers, err := proto.Marshal(req)
	if nil != err {
		return fmt.Errorf("proto marshal data failed, err is %s", err.Error())
	}
	if err := sendFrame(codec.w, headers); nil != err {
		return fmt.Errorf("write data failed, err is %s", err.Error())
	}
	return nil
}

func (codec *pbClientCodec) writeRequestBody(data []byte) error {
	loger.Tracef("rpcd.pbClientCodec.writeRequestBody: write request body len[%d]", len(data))
	if err := sendFrame(codec.w, data); nil != err {
		return fmt.Errorf("write data failed, err is %s", err.Error())
	}
	return nil
}

func (codec *pbClientCodec) transferRequestBody(x interface{}) ([]byte, error) {
	if pbRequest, ok := x.(proto.Message); ok {
		return proto.Marshal(pbRequest)
	}
	return nil, fmt.Errorf("%T does not implement proto.Message", x)
}

func (codec *pbClientCodec) pending(header *RequestHeader) {
	codec.m.Lock()
	defer codec.m.Unlock()
	respHeader := &ResponseHeader{
		Id:     header.Id,
		Method: header.Method,
	}
	codec.pends[header.Id] = respHeader
}

func (codec *pbClientCodec) pended(header *ResponseHeader) (*ResponseHeader, error) {

	codec.m.Lock()
	defer codec.m.Unlock()

	if respHeader, exist := codec.pends[header.Id]; exist {
		header.Method = respHeader.Method
		delete(codec.pends, header.Id)
		return header, nil
	}
	return nil, errors.New("not exist seq in pends")
}
