package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"

	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

func ValidateDatabaseName(name string) bool {
	if len(name) < 3 {
		return false
	}

	parts := strings.Split(name, ".")

	if len(parts) < 2 {
		return false
	}

	// invalid extension
	if len(parts[len(parts)-1]) < 2 {
		return false
	}

	for i := 0; i < len(parts); i++ {
		rxp := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
		if !rxp.MatchString(parts[i]) {
			return false
		}
	}

	return true
}

// MergeSet add value to a ordered set of integers stored in key. If value
// is already on the key, than the set will be skipped.
func MergeSet(writer KVWriter, key []byte, value uint64, debug bool) error {
	var (
		buf      *bytes.Buffer
		err      error
		v        uint64
		i        uint64
		inserted bool
	)

	data, err := writer.Get(key)
	if err != nil {
		return err
	}

	if debug {
		fmt.Printf("[INFO] %d ids == %d GB of ids\n", len(data)/8, len(data)/(1024*1024*1024))
	}

	buf = new(bytes.Buffer)
	lenBytes := uint64(len(data))

	// O(n)
	for i = 0; i < lenBytes; i += 8 {
		v = utils.BytesToUint64(data[i : i+8])

		// returns if value is already stored
		if v == value {
			return nil
		}

		if value < v {
			err = binary.Write(buf, binary.BigEndian, value)
			if err != nil {
				return err
			}
			inserted = true
		}

		err = binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			return err
		}
	}

	if lenBytes == 0 || !inserted {
		err = binary.Write(buf, binary.BigEndian, value)
		if err != nil {
			return err
		}
	}

	return writer.Set(key, buf.Bytes())
}
