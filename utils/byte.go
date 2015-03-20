package utils

import (
	"bytes"
	"encoding/binary"
)

func Uint64ToBytes(i uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)

	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func BytesToUint64(b []byte) uint64 {
	var i uint64

	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		panic(err)
	}

	return i
}

func Float64ToBytes(f float64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, f)

	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
