// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package codec

import (
	"fmt"
	"io"
	"math"

	"strconv"
)

func readByte(reader io.Reader) (byte, error) {
	var data bytes1
	_, err := reader.Read(data[0:])
	if nil != err {
		return 0, err
	}

	return data[0], nil
}

func readUint16(reader io.Reader) (uint16, error) {
	var data bytes2
	_, err := reader.Read(data[0:])
	if nil != err {
		return 0, err
	}

	return (uint16(data[0]) << 8) | uint16(data[1]), nil
}

func readUint32(reader io.Reader) (uint32, error) {
	var data bytes4
	_, err := reader.Read(data[0:])
	if nil != err {
		return 0, err
	}

	return (uint32(data[0]) << 24) | (uint32(data[1]) << 16) | (uint32(data[2]) << 8) | uint32(data[3]), nil
}

func readUint64(reader io.Reader) (uint64, error) {
	var data bytes8
	_, err := reader.Read(data[0:])
	if err != nil {
		return 0, err
	}

	return (uint64(data[0]) << 56) | (uint64(data[1]) << 48) | (uint64(data[2]) << 40) | (uint64(data[3]) << 32) | (uint64(data[4]) << 24) | (uint64(data[5]) << 16) | (uint64(data[6]) << 8) | uint64(data[7]), nil
}

func readInt16(reader io.Reader) (int16, error) {
	var data bytes2
	_, err := reader.Read(data[0:])
	if nil != err {
		return 0, err
	}

	return (int16(data[0]) << 8) | int16(data[1]), nil
}

func readInt32(reader io.Reader) (int32, error) {
	var data bytes4
	_, err := reader.Read(data[0:])
	if nil != err {
		return 0, err
	}

	return (int32(data[0]) << 24) | (int32(data[1]) << 16) | (int32(data[2]) << 8) | int32(data[3]), nil
}

func readInt64(reader io.Reader) (int64, error) {
	var data bytes8
	_, err := reader.Read(data[0:])
	if nil != err {
		return 0, err
	}

	return (int64(data[0]) << 56) | (int64(data[1]) << 48) | (int64(data[2]) << 40) | (int64(data[3]) << 32) | (int64(data[4]) << 24) | (int64(data[5]) << 16) | (int64(data[6]) << 8) | int64(data[7]), nil
}

func decodeArray(reader io.Reader, count uint) ([]interface{}, error) {
	var i uint
	arr := make([]interface{}, count)
	for i = 0; i < count; i++ {
		v, err := decode(reader)
		if nil != err {
			return nil, err
		}

		arr[i] = v
	}

	return arr, nil
}

func decodeMap(reader io.Reader, count uint) (map[interface{}]interface{}, error) {
	var i uint
	dict := make(map[interface{}]interface{})

	for i = 0; i < count; i++ {
		k, err := decode(reader)
		if nil != err {
			return nil, err
		}

		v, err := decode(reader)
		if nil != err {
			return nil, err
		}

		dict[k] = v
	}

	return dict, nil
}

func decode(reader io.Reader) (v interface{}, err error) {
	c, err := readByte(reader)
	if nil != err {
		return nil, err
	}

	switch c {
	case cNil:
		return nil, nil
	case cFalse:
		return false, nil
	case cTrue:
		return true, nil
	case cFloat32:
		{
			data, err := decode(reader)
			if nil != err {
				return nil, err
			}

			return math.Float32frombits(data.(uint32)), nil
		}
	case cFloat64:
		{
			data, err := decode(reader)
			if nil != err {
				return nil, err
			}

			return math.Float64frombits(data.(uint64)), nil
		}
	case cUint8:
		{
			data, err := readByte(reader)
			if nil != err {
				return nil, err
			}

			return data, nil
		}
	case cUint16:
		{
			data, err := readUint16(reader)
			if nil != err {
				return nil, err
			}

			return data, nil
		}
	case cUint32:
		{
			data, err := readUint32(reader)
			if nil != err {
				return nil, err
			}

			return data, nil
		}
	case cUint64:
		{
			data, err := readUint64(reader)
			if nil != err {
				return nil, err
			}

			return data, nil
		}
	case cInt8:
		{
			data, err := readByte(reader)
			if nil != err {
				return nil, err
			}

			return int8(data), nil
		}
	case cInt16:
		{
			data, err := readInt16(reader)
			if nil != err {
				return nil, err
			}

			return data, nil
		}
	case cInt32:
		{
			data, err := readInt32(reader)
			if nil != err {
				return nil, err
			}

			return data, nil
		}
	case cInt64:
		{
			data, err := readInt64(reader)
			if nil != err {
				return nil, err
			}

			return data, nil
		}
	case cStr32:
		{
			n, err := readUint32(reader)
			if nil != err {
				return nil, err
			}

			if n == 0 {
				return "", nil
			}

			data := make(bytes, n)
			_, err = reader.Read(data)
			if nil != err {
				return nil, err
			}

			return string(data), nil
		}
	case cArr32:
		{
			n, err := readUint32(reader)
			if nil != err {
				return nil, err
			}

			return decodeArray(reader, uint(n))
		}
	case cMap32:
		{
			n, err := readUint32(reader)
			if nil != err {
				return nil, err
			}

			return decodeMap(reader, uint(n))
		}
	}

	return nil, fmt.Errorf("Unsupported code: %s", strconv.Itoa(int(c)))
}
