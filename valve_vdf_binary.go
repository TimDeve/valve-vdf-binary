package valve_vdf_binary

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func Parse(r io.Reader) (VdfValue, error) {
	buf := bufio.NewReader(r)

	byteArr, err := buf.Peek(1)

	if err == io.EOF {
		return vdfValue{}, errors.New("The vdf you are trying to parse appears empty")
	}

	if err != nil {
		return vdfValue{}, err
	}

	b := byteArr[0]

	if b != vdfMarkerMap &&
		b != vdfMarkerString &&
		b != vdfMarkerNumber &&
		b != vdfMarkerEndOfMap {
		return vdfValue{}, errors.New(
			"The vdf you are trying to parse appears not to be binary," +
				"are you sure it is not a text vdf?",
		)
	}

	p, err := parseMap(buf)

	if err == io.EOF {
		return vdfValue{}, errors.New("Reached the end of the file earlier than expected, your file might be corrupted")
	}

	return p, err
}

func parseMap(buf *bufio.Reader) (vdfValue, error) {
	m := make(VdfMap)

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
			value, err = parseMap(buf)
		case vdfMarkerNumber:
			value, err = parseNumber(buf)
		case vdfMarkerString:
			value, err = parseStringValue(buf)
		default:
			err = fmt.Errorf("Unexpected byte: 0x%02x, your file might be corrupted", b)
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
