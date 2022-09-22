package valve_vdf_binary_test

import (
	"bytes"
	"testing"

	. "github.com/TimDeve/valve-vdf-binary"
	test "github.com/stretchr/testify/require"
)

func TestEmptyFileError(t *testing.T) {
	r := bytes.NewReader([]byte{})

	_, err := Parse(r)

	test.ErrorContains(t, err, "The vdf you are trying to parse appears empty")
}

func TestFileIsNotBinaryError(t *testing.T) {
	r := bytes.NewReader(
		[]byte(`"vdf"
{
	"vdfvalue"	"1"
}`))

	_, err := Parse(r)

	test.ErrorContains(t, err, "The vdf you are trying to parse appears not to be binary,are you sure it is not a text vdf?")
}

func TestEofInMiddleOfFile(t *testing.T) {
	r := bytes.NewReader([]byte{0x02})

	_, err := Parse(r)

	test.ErrorContains(t, err, "Reached the end of the file earlier than expected, your file might be corrupted")
}

func TestMapWithNestedMap(t *testing.T) {
	r := bytes.NewReader([]byte{0x00, 'o', 'n', 'e', 0x00, 0x02, 't', 'w', 'o', 0x00, 0x39, 0x05, 0x00, 0x00, 0x08, 0x08})

	m, err := Parse(r)

	test.Nil(t, err)

	expected := MakeVdfValue(VdfMap{
		"one": MakeVdfValue(VdfMap{
			"two": MakeVdfValue(uint32(1337)),
		}),
	})

	test.Equal(t, expected, m)
}

func TestJustOneStringField(t *testing.T) {
	r := bytes.NewReader([]byte{0x01, 'o', 'n', 'e', 0x00, 't', 'w', 'o', 0x00, 0x08})

	m, err := Parse(r)

	test.Nil(t, err)

	expected := MakeVdfValue(VdfMap{
		"one": MakeVdfValue("two"),
	})

	test.Equal(t, expected, m)
}

func TestJustOneNumberField(t *testing.T) {
	r := bytes.NewReader([]byte{0x02, 'o', 'n', 'e', 0x00, 0x02, 0x00, 0x00, 0x00, 0x08})

	m, err := Parse(r)

	test.Nil(t, err)

	expected := MakeVdfValue(VdfMap{
		"one": MakeVdfValue(uint32(2)),
	})

	test.Equal(t, expected, m)
}

func TestEmptyMap(t *testing.T) {
	r := bytes.NewReader([]byte{0x08})

	m, err := Parse(r)

	test.Nil(t, err)

	expected := MakeVdfValue(VdfMap{})

	test.Equal(t, expected, m)
}
