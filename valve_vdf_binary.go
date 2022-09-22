package valve_vdf_binary

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	vdfMarkerMap         byte = 0x00
	vdfMarkerString      byte = 0x01
	vdfMarkerNumber      byte = 0x02
	vdfMarkerEndOfMap    byte = 0x08
	vdfMarkerEndOfString byte = 0x00
)

type vdfMap map[string]vdfValue

type vdfValue struct {
	data any
}

func (sv *vdfValue) asString() (string, bool) {
	s, ok := sv.data.(string)
	return s, ok
}

func (sv *vdfValue) getString(key string) (string, bool) {
	m, ok := sv.asMap()
	if !ok {
		return "", ok
	}
	v, ok := m[key]
	if !ok {
		return "", ok
	}
	return v.asString()
}

func (sv *vdfValue) asUint() (uint32, bool) {
	i, ok := sv.data.(uint32)
	return i, ok
}

func (sv *vdfValue) getUint(key string) (uint32, bool) {
	m, ok := sv.asMap()
	if !ok {
		return 0, ok
	}
	v, ok := m[key]
	if !ok {
		return 0, ok
	}
	return v.asUint()
}

func (sv *vdfValue) asInt() (int, bool) {
	i, ok := sv.data.(uint32)
	return int(i), ok
}

func (sv *vdfValue) getInt(key string) (int, bool) {
	m, ok := sv.asMap()
	if !ok {
		return 0, ok
	}
	v, ok := m[key]
	if !ok {
		return 0, ok
	}
	return v.asInt()
}

func (sv *vdfValue) asFloat() (float32, bool) {
	f, ok := sv.data.(uint32)
	return float32(f), ok
}

func (sv *vdfValue) getFloat(key string) (float32, bool) {
	m, ok := sv.asMap()
	if !ok {
		return 0, ok
	}
	v, ok := m[key]
	if !ok {
		return 0, ok
	}
	return v.asFloat()
}

func (sv *vdfValue) asBool() (bool, bool) {
	i, ok := sv.data.(uint32)
	if i == 0 {
		return false, ok
	} else {
		return true, ok
	}
}

func (sv *vdfValue) getBool(key string) (bool, bool) {
	m, ok := sv.asMap()
	if !ok {
		return false, ok
	}
	v, ok := m[key]
	if !ok {
		return false, ok
	}
	return v.asBool()
}

func (sv *vdfValue) asMap() (vdfMap, bool) {
	m, ok := sv.data.(map[string]vdfValue)
	return m, ok
}

func (sv *vdfValue) getMap(key string) (vdfMap, bool) {
	m, ok := sv.asMap()
	if !ok {
		return nil, ok
	}
	v, ok := m[key]
	if !ok {
		return nil, ok
	}
	return v.asMap()
}

func parse(buf *bufio.Reader) (vdfValue, error) {
	m := make(map[string]vdfValue)

	for {
		b, err := buf.ReadByte()
		if err != nil {
			return vdfValue{}, err
		}

		if b == vdfMarkerEndOfMap {
			break
		}

		key, err := parseString(buf)
		if err != nil {
			return vdfValue{}, err
		}

		value, err := vdfValue{}, nil
		switch b {
		case vdfMarkerMap:
			value, err = parse(buf)
		case vdfMarkerNumber:
			value, err = parseNumber(buf)
		case vdfMarkerString:
			value, err = parseStringValue(buf)
		default:
			err = fmt.Errorf("Unexpected byte: 0x%02x", b)
		}

		if err != nil {
			return vdfValue{}, err
		}

		m[key] = value
	}

	return vdfValue{m}, nil
}

func parseNumber(buf *bufio.Reader) (vdfValue, error) {
	bf := make([]byte, 4)
	l, err := buf.Read(bf)
	if err != nil {
		return vdfValue{}, err
	}
	if l != len(bf) {
		return vdfValue{}, errors.New("Number did not have the required amound of bytes")
	}

	number := binary.LittleEndian.Uint32(bf)

	return vdfValue{number}, nil
}

func parseString(buf *bufio.Reader) (string, error) {
	s, err := buf.ReadString(vdfMarkerEndOfString)
	if err == nil {
		return s[:len(s)-1], nil // Removes the null terminator
	}
	return "", err
}

func parseStringValue(buf *bufio.Reader) (vdfValue, error) {
	s, err := parseString(buf)
	return vdfValue{s}, err
}
