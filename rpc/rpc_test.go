package rpc_test

import (
	"basiclsp/rpc"
	"testing"
)

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"

	actual := rpc.EncodeMessage(EncodingExample{Testing: true})

	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	message := "Content-Length: 17\r\n\r\n{\"method\":\"test\"}"

	decoded, _, err := rpc.DecodeMessage([]byte(message))
	if err != nil {
		t.Fatal(err)
	}

	expectedDecoded := "test"
	if decoded != expectedDecoded {
		t.Errorf("Expected method %s but got %s", expectedDecoded, decoded)
	}
}
