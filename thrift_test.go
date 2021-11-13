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

func readAndParseFile(filePath string) (res *Thrift, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	parser := NewParser(file, true)
	res = NewThrift(nil, "")
	if err = res.parse(parser); err != nil {
		return nil, err
	}
	return
}

// https://raw.githubusercontent.com/apache/thrift/master/test/ThriftTest.thrift
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

	if got, want := res.String(), src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

// https://raw.githubusercontent.com/apache/thrift/master/contrib/fb303/if/fb303.thrift
func TestThrift_toStringFB303(t *testing.T) {
	src, err := readFile("./examples/fb303.thrift")
	if err != nil {
		t.Errorf("readFile error: %v", err)
		return
	}
	res, err := readAndParseFile("./examples/fb303.thrift")
	if err != nil {
		t.Errorf("readAndParseFile error: %v", err)
		return
	}

	if got, want := res.String(), src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

// https://gitbox.apache.org/repos/asf?p=cassandra.git;a=blob_plain;f=interface/cassandra.thrift;hb=refs/heads/cassandra-3.0
func TestThrift_toStringCassandra(t *testing.T) {
	src, err := readFile("./examples/cassandra.thrift")
	if err != nil {
		t.Errorf("readFile error: %v", err)
		return
	}
	res, err := readAndParseFile("./examples/cassandra.thrift")
	if err != nil {
		t.Errorf("readAndParseFile error: %v", err)
		return
	}

	if got, want := res.String(), src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

// https://git-wip-us.apache.org/repos/asf?p=thrift.git;a=blob_plain;f=test/AnnotationTest.thrift;h=06bf57194fc06a6b9eae7138e196bdefca260618;hb=HEAD
func TestThrift_toStringAnnotationTest(t *testing.T) {
	src, err := readFile("./examples/annotation_test.thrift")
	if err != nil {
		t.Errorf("readFile error: %v", err)
		return
	}
	res, err := readAndParseFile("./examples/annotation_test.thrift")
	if err != nil {
		t.Errorf("readAndParseFile error: %v", err)
		return
	}

	if got, want := res.String(), src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestThrift_basicAnnotationTest(t *testing.T) {
	res, err := readAndParseFile("./examples/annotation_test.thrift")
	if err != nil {
		t.Errorf("readAndParseFile error: %v", err)
		return
	}
	firstTypeDef, _ := res.Nodes[0].NodeValue().(TypeDef)

	if got, want := len(res.Nodes), 7; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := res.Nodes[0].NodeType(), "TypeDef"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := res.Nodes[1].NodeType(), "Struct"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := res.Nodes[2].NodeType(), "Exception"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := res.Nodes[3].NodeType(), "TypeDef"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := res.Nodes[4].NodeType(), "TypeDef"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := res.Nodes[5].NodeType(), "Enum"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := res.Nodes[6].NodeType(), "Service"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := firstTypeDef.Type.Type, FIELD_TYPE_LIST; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
