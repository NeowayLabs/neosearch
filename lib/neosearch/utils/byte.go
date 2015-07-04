package utils

import (
	"bytes"
	"encoding/binary"
)

func BoolToBytes(b bool) []byte {
	var (
		bs []byte = make([]byte, 1)
	)

	if b {
		bs[0] = byte(1)
	} else {
		bs[0] = byte(0)
	}

	return bs
}

func BytesToBool(b []byte) bool {
	if len(b) == 0 {
		panic("invalid boolean byte-array")
	}

	bv := b[0]

	if bv == 0 {
		return false
	}

	return true
}

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

func Int64ToBytes(i int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)

	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func BytesToInt64(b []byte) int64 {
	var i int64

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

func BytesToFloat64(b []byte) float64 {
	var f float64

	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &f)
	if err != nil {
		panic(err)
	}

	return f
}

func GetUint64Array(data []byte) []uint64 {
	var i, v uint64

	lenBytes := uint64(len(data))
	uints := make([]uint64, lenBytes/8)

	for i = 0; i < lenBytes; i += 8 {
		v = BytesToUint64(data[i : i+8])
		uints[i] = v
	}

	return uints
}
