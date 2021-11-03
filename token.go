package thrifter

import (
	"bytes"
	"strings"
)

type token int

const baseTypeTokens = "bool byte i8 i16 i32 i64 double string binary slist"

const (
	// special tokens
	tILLEGAL token = iota
	tEOF
	tIDENT
	tSTRING // string literal
	tNUMBER // integer or float

	// white space
	tSPACE
	tLINEBREAK // \n
	tRETURN    // \r
	tTAB       // \t

	// punctuator
	tSEMICOLON   // ;
	tCOLON       // :
	tEQUALS      // =
	tQUOTE       // "
	tSINGLEQUOTE // '
	tLEFTPAREN   // (
	tRIGHTPAREN  // )
	tLEFTCURLY   // {
	tRIGHTCURLY  // }
	tLEFTSQUARE  // [
	tRIGHTSQUARE // ]
	tCOMMENT     // /
	tLESS        // <
	tGREATER     // >
	tCOMMA       // ,
	tDOT         // .
	tPLUS        // +
	tMINUS       // -

	// declaration keywords
	keywordStart
	tNAMESPACE
	tENUM
	tCONST
	tSERVICE
	tSTRUCT
	tINCLUDE
	tTYPEDEF
	tUNION
	tEXCEPTION

	// field keywords
	tOPTIONAL
	tREQUIRED

	// type keywords
	tMAP
	tSET
	tLIST

	// function keywords
	tONEWAY
	tVOID
	tTHROWS
	keywordEnd
)

func toToken(literal string) token {
	switch literal {
	// white space
	case "\n":
		return tLINEBREAK
	case "\r":
		return tRETURN
	case " ":
		return tSPACE
	case "\t":
		return tTAB
	// punctuator
	case ";":
		return tSEMICOLON
	case ":":
		return tCOLON
	case "=":
		return tEQUALS
	case "\"":
		return tQUOTE
	case "'":
		return tSINGLEQUOTE
	case "(":
		return tLEFTPAREN
	case ")":
		return tRIGHTPAREN
	case "{":
		return tLEFTCURLY
	case "}":
		return tRIGHTCURLY
	case "[":
		return tLEFTSQUARE
	case "]":
		return tRIGHTSQUARE
	case "<":
		return tLESS
	case ">":
		return tGREATER
	case ",":
		return tCOMMA
	case ".":
		return tDOT
	case "+":
		return tPLUS
	case "-":
		return tMINUS

	// declaration keywords
	case "namespace":
		return tNAMESPACE
	case "enum":
		return tENUM
	case "const":
		return tCONST
	case "service":
		return tSERVICE
	case "struct":
		return tSTRUCT
	case "include":
		return tINCLUDE
	case "typedef":
		return tTYPEDEF
	case "union":
		return tUNION
	case "exception":
		return tEXCEPTION

	// field keywords
	case "optional":
		return tOPTIONAL
	case "required":
		return tREQUIRED

	// type keywords
	case "map":
		return tMAP
	case "set":
		return tSET
	case "list":
		return tLIST

	// function keywords
	case "oneway":
		return tONEWAY
	case "void":
		return tVOID
	case "throws":
		return tTHROWS
	default:
		return tIDENT
	}
}

// comment type
const (
	SINGLE_LINE_COMMENT = iota + 1 // like this
	MULTI_LINE_COMMENT             /* like this */
	BASH_LIKE_COMMENT              // # like this
)

// isDigit returns true if the rune is a digit.
func isDigit(lit rune) bool {
	return (lit >= '0' && lit <= '9')
}

func getCommentValue(raw string, commentType int) (res string) {
	switch commentType {
	case SINGLE_LINE_COMMENT:
		res = strings.Replace(raw, "//", "", 1)
	case MULTI_LINE_COMMENT:
		res = strings.ReplaceAll(raw, "/*", "")
		res = strings.ReplaceAll(res, "*/", "")
	case BASH_LIKE_COMMENT:
		res = strings.Replace(raw, "#", "", 1)
	}
	return
}

// isKeyword returns if tok is in the keywords range
func isKeyword(tok token) bool {
	return keywordStart < tok && tok < keywordEnd
}

func isWhitespace(tok token) bool {
	return tok == tSPACE || tok == tLINEBREAK || tok == tRETURN || tok == tTAB
}

func toString(start *Token, end *Token) string {
	var res bytes.Buffer
	curr := start
	for curr != end {
		res.WriteString(curr.Raw)
		curr = curr.Next
	}
	res.WriteString(end.Raw)
	return res.String()
}

const singleQuoteString = "'"
const singleQuoteRune = '\''
const doubleQuoteString = "\""
const doubleQuoteRune = '"'

// UnQuote removes one matching leading and trailing single or double quote.
// cannot use strconv.Unquote as this unescapes quotes.
func unQuote(lit string) (string, rune) {
	if len(lit) < 2 {
		return lit, doubleQuoteRune
	}
	chars := []rune(lit)
	first, last := chars[0], chars[len(chars)-1]
	if first != last {
		return lit, doubleQuoteRune
	}
	if s := string(chars[0]); s == doubleQuoteString || s == singleQuoteString {
		return string(chars[1 : len(chars)-1]), chars[0]
	}
	return lit, doubleQuoteRune
}