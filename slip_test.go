package slip

import (
	"testing"
)

func TestEncode(t *testing.T) {
	data := []byte{42}
	t.Log("Data", data)

	encoded := Encode(data)
	t.Log("Encoded", encoded)

	if encoded[0] != end {
		t.Fail()
	}

	if encoded[1] != data[0] {
		t.Fail()
	}

	if encoded[2] != end {
		t.Fail()
	}
}

func TestEncodeEscapeEnd(t *testing.T) {
	data := []byte{end}
	t.Log("Data", data)

	encoded := Encode(data)
	t.Log("Encoded", encoded)

	if encoded[0] != end {
		t.Fail()
	}

	if encoded[1] != esc {
		t.Fail()
	}

	if encoded[2] != escEnd {
		t.Fail()
	}

	if encoded[3] != end {
		t.Fail()
	}
}

func TestEncodeEscapeEscape(t *testing.T) {
	data := []byte{esc}
	t.Log("Data", data)

	encoded := Encode(data)
	t.Log("Encoded", encoded)

	if encoded[0] != end {
		t.Fail()
	}

	if encoded[1] != esc {
		t.Fail()
	}

	if encoded[2] != escEsc {
		t.Fail()
	}

	if encoded[3] != end {
		t.Fail()
	}
}

func TestDecode(t *testing.T) {
	data := []byte{42, 55, 111, 0}
	t.Log("Data:", data)

	encoded := Encode(data)
	t.Log("Encoded:", encoded)

	in := make(chan byte, 1024)
	for _, d := range encoded {
		in <- d
	}

	out := make(chan []byte)
	go Decode(in, out)

	//time.Sleep(100 * time.Millisecond)

	decoded := <-out

	t.Log("Decoded:", decoded)
	if len(data) != len(decoded) {
		t.Fail()
	}
	for i, d := range data {
		if d != data[i] {
			t.Fail()
		}
	}
}

func TestDecodeEscaped(t *testing.T) {
	data := []byte{42, 55, end, 111, esc, 0}
	t.Log("Data:", data)

	encoded := Encode(data)
	t.Log("Encoded:", encoded)

	in := make(chan byte, 1024)
	for _, d := range encoded {
		in <- d
	}

	out := make(chan []byte)
	go Decode(in, out)

	//time.Sleep(100 * time.Millisecond)

	decoded := <-out

	t.Log("Decoded:", decoded)
	if len(data) != len(decoded) {
		t.Fail()
	}
	for i, d := range data {
		if d != data[i] {
			t.Fail()
		}
	}
}
