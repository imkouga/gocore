package rpcd

import (
	"encoding/binary"
	"errors"
	"io"
	"net"

	"github.com/imkouga/gocore/loger"
)

func sendFrame(w io.Writer, data []byte) error {

	// Allocate enough space for the biggest uvarint
	var size [binary.MaxVarintLen64]byte

	if len(data) <= 0 {
		n := binary.PutUvarint(size[:], uint64(0))
		return write(w, size[:n], false)
	}

	// Write the size and data
	n := binary.PutUvarint(size[:], uint64(len(data)))
	if err := write(w, size[:n], false); err != nil {
		return err
	}
	if err := write(w, data, false); err != nil {
		return err
	}
	return nil
}

func recvFrame(r io.Reader) ([]byte, error) {

	var (
		size uint64
		data []byte
		err  error
	)

	size, err = readUvarint(r)
	loger.Trace("rpcd.recvFrame size:", size, ", err: ", err)
	if err != nil || size <= 0 {
		return nil, err
	}

	data = make([]byte, size)
	err = read(r, data)
	return data, err
}

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func readUvarint(r io.Reader) (uint64, error) {

	var (
		x   uint64
		s   uint
		b   byte
		err error
	)

	for i := 0; ; i++ {
		if b, err = readByte(r); nil != err {
			return 0, err
		}
		if i > 9 || (i == 9 && b > 1) {
			return x, errors.New("rpcd: varint overflows a 64-bit integer")
		}
		if b < 0x80 {
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

func read(r io.Reader, data []byte) error {

	for index := 0; index < len(data); {
		n, err := r.Read(data[index:])
		if err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				return err
			}
		}
		index += n
	}
	return nil
}

func readByte(r io.Reader) (byte, error) {
	data := make([]byte, 1)
	if err := read(r, data); err != nil {
		return 0, err
	}
	return data[0], nil
}

func write(w io.Writer, data []byte, onePacket bool) error {

	if onePacket {
		_, err := w.Write(data)
		return err
	}
	for index := 0; index < len(data); {
		n, err := w.Write(data[index:])
		if err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				return err
			}
		}
		index += n
	}
	return nil
}
