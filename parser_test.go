package thrifter

import (
	"fmt"
	"strings"
	"testing"
)

func Test_tokens(t *testing.T) {
	rd := strings.NewReader(`// first comment
/* block comment */ // inline comment
aa
ab	11
  bb
/* 123123
	1231231
	aaaaaa
*/
`)
	parser := NewParser(rd)
	res, _ := parser.Parse()
	curTok := res.StartToken

	fmt.Printf("%+v\n", curTok)
	for tok := curTok.Next; tok != nil; tok = tok.Next {
		fmt.Printf("%+v\n", tok)
	}
}
