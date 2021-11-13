package thrifter

import (
	"io"
	"os"
	"testing"
)

func readFile(filePath string) (res string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	res = string(bytes)
	return
}

func readAndParseFile(filePath string) (res string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	parser := NewParser(file, true)
	n := NewThrift(nil, "")
	if err = n.parse(parser); err != nil {
		return "", err
	}
	res = n.String()
	return
}

// Test thrift file is a copy from https://raw.githubusercontent.com/apache/thrift/master/test/ThriftTest.thrift.
func TestThrift_toStringThriftTest(t *testing.T) {
	src, err := readFile("./examples/ThriftTest.thrift")
	if err != nil {
		t.Errorf("readFile error: %v", err)
		return
	}
	res, err := readAndParseFile("./examples/ThriftTest.thrift")
	if err != nil {
		t.Errorf("readAndParseFile error: %v", err)
		return
	}

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
