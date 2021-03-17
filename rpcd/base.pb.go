package rpcd

import (
	"errors"
	"hash/crc32"

	"github.com/golang/protobuf/proto"
)

const (
	RpcdSourceGrowerGateway = "grower-gateway"
)

type SugarCoatingRequest struct {
	RpcdSource string `protobuf:"bytes,1,opt,name=RpcdSource" json:"RpcdSource,omitempty"`
	Data       []byte `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
}

func (m *SugarCoatingRequest) Reset()         { *m = SugarCoatingRequest{} }
func (m *SugarCoatingRequest) String() string { return proto.CompactTextString(m) }
func (*SugarCoatingRequest) ProtoMessage()    {}

type SugarCoatingResponseInterface interface {
	GetSugarData() ([]byte, error)
}

type SugarCoatingResponse struct {
	Data []byte `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *SugarCoatingResponse) Reset()                        { *m = SugarCoatingResponse{} }
func (m *SugarCoatingResponse) String() string                { return proto.CompactTextString(m) }
func (*SugarCoatingResponse) ProtoMessage()                   {}
func (m *SugarCoatingResponse) GetSugarData() ([]byte, error) { return m.Data, nil }

type RequestHeader struct {
	Id         uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Method     string `protobuf:"bytes,2,opt,name=method" json:"method,omitempty"`
	RpcdSource string `protobuf:"bytes,3,opt,name=rpcdSource" json:"rpcdSource,omitempty"`
	Checksum   uint32 `protobuf:"varint,5,opt,name=checksum" json:"checksum,omitempty"`
}

func (m *RequestHeader) Reset()         { *m = RequestHeader{} }
func (m *RequestHeader) String() string { return proto.CompactTextString(m) }
func (*RequestHeader) ProtoMessage()    {}
func (m *RequestHeader) Check(body []byte) error {
	if m.Checksum != crc32.ChecksumIEEE(body) {
		return errors.New("crc32 check sum failed.")
	}
	return nil
}

type ResponseHeader struct {
	Id       uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Error    string `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
	Method   string `protobuf:"bytes,2,opt,name=method" json:"method,omitempty"`
	Checksum uint32 `protobuf:"varint,5,opt,name=checksum" json:"checksum,omitempty"`
}

func (m *ResponseHeader) Reset()         { *m = ResponseHeader{} }
func (m *ResponseHeader) String() string { return proto.CompactTextString(m) }
func (*ResponseHeader) ProtoMessage()    {}
func (m *ResponseHeader) Check(body []byte) error {
	if m.Checksum != crc32.ChecksumIEEE(body) {
		return errors.New("crc32 check sum failed.")
	}
	return nil
}

func GenSugarCoatingRequest(source string, data []byte) *SugarCoatingRequest {
	return &SugarCoatingRequest{RpcdSource: source, Data: data}
}
