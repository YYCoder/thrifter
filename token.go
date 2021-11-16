package thrifter

import (
	"bytes"
	"regexp"
	"strings"
)

type token int

var baseTypeTokens = []string{"bool", "byte", "i8", "i16", "i32", "i64", "double", "string", "binary", "slist"}

func isBaseTypeToken(str string) bool {
	for _, s := range baseTypeTokens {
		if s == str {
			return true
		}
	}
	return false
}

const (
	// special tokens
	T_ILLEGAL token = iota
	T_EOF
	T_IDENT
	T_STRING // string literal
	T_NUMBER // integer or float

	// white space
	T_SPACE
	T_LINEBREAK // \n
	T_RETURN    // \r
	T_TAB       // \t

	// punctuator
	T_SEMICOLON   // ;
	T_COLON       // :
	T_EQUALS      // =
	T_QUOTE       // "
	T_SINGLEQUOTE // '
	T_LEFTPAREN   // (
	T_RIGHTPAREN  // )
	T_LEFTCURLY   // {
	T_RIGHTCURLY  // }
	T_LEFTSQUARE  // [
	T_RIGHTSQUARE // ]
	T_COMMENT     // /
	T_LESS        // <
	T_GREATER     // >
	T_COMMA       // ,
	T_DOT         // .
	T_PLUS        // +
	T_MINUS       // -

	// declaration keywords
	keywordStart
	T_NAMESPACE
	T_ENUM
	T_SENUM // currently not supported
	T_CONST
	T_SERVICE
	T_STRUCT
	T_INCLUDE
	T_CPP_INCLUDE
	T_TYPEDEF
	T_UNION
	T_EXCEPTION

	// field keywords
	T_OPTIONAL
	T_REQUIRED

	// type keywords
	T_MAP
	T_SET
	T_LIST

	// function keywords
	T_ONEWAY
	T_VOID
	T_THROWS
	keywordEnd
)

// Get corresponding token from string literal, mostly used for generate token.
func GetToken(literal string) token {
	return toToken(literal)
}

func toToken(literal string) token {
	switch literal {
	// white space
	case "\n":
		return T_LINEBREAK
	case "\r":
		return T_RETURN
	case " ":
		return T_SPACE
	case "\t":
		return T_TAB
	// punctuator
	case ";":
		return T_SEMICOLON
	case ":":
		return T_COLON
	case "=":
		return T_EQUALS
	case "\"":
		return T_QUOTE
	case "'":
		return T_SINGLEQUOTE
	case "(":
		return T_LEFTPAREN
	case ")":
		return T_RIGHTPAREN
	case "{":
		return T_LEFTCURLY
	case "}":
		return T_RIGHTCURLY
	case "[":
		return T_LEFTSQUARE
	case "]":
		return T_RIGHTSQUARE
	case "<":
		return T_LESS
	case ">":
		return T_GREATER
	case ",":
		return T_COMMA
	case ".":
		return T_DOT
	case "+":
		return T_PLUS
	case "-":
		return T_MINUS

	// declaration keywords
	case "namespace":
		return T_NAMESPACE
	case "enum":
		return T_ENUM
	case "senum":
		return T_SENUM
	case "const":
		return T_CONST
	case "service":
		return T_SERVICE
	case "struct":
		return T_STRUCT
	case "include":
		return T_INCLUDE
	case "cpp_include":
		return T_CPP_INCLUDE
	case "typedef":
		return T_TYPEDEF
	case "union":
		return T_UNION
	case "exception":
		return T_EXCEPTION

	// field keywords
	case "optional":
		return T_OPTIONAL
	case "required":
		return T_REQUIRED

	// type keywords
	case "map":
		return T_MAP
	case "set":
		return T_SET
	case "list":
		return T_LIST

	// function keywords
	case "oneway":
		return T_ONEWAY
	case "void":
		return T_VOID
	case "throws":
		return T_THROWS
	default:
		return T_IDENT
	}
}

// comment type
const (
	SINGLE_LINE_COMMENT = iota + 1 // like this
	MULTI_LINE_COMMENT             /* like this */
	BASH_LIKE_COMMENT              // # like this
)

// isDigit returns true if the rune is a digit.
func IsDigit(lit rune) bool {
	return (lit >= '0' && lit <= '9')
}

// determine whether it is an integer or a float number
func IsNumber(str string) (isFloat bool, isInt bool) {
	isFloat, _ = regexp.MatchString("^\\d+\\.\\d+$", str)
	isInt, _ = regexp.MatchString("^\\d+$", str)
	return
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
func IsKeyword(tok token) bool {
	return keywordStart < tok && tok < keywordEnd
}

func IsWhitespace(tok token) bool {
	return tok == T_SPACE || tok == T_LINEBREAK || tok == T_RETURN || tok == T_TAB
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
const quoteString = "\""
const quoteRune = '"'

// UnQuote removes one matching leading and trailing single or double quote.
// cannot use strconv.Unquote as this unescapes quotes.
func unQuote(lit string) (string, rune) {
	if len(lit) < 2 {
		return lit, quoteRune
	}
	chars := []rune(lit)
	first, last := chars[0], chars[len(chars)-1]
	if first != last {
		return lit, quoteRune
	}
	if s := string(chars[0]); s == quoteString || s == singleQuoteString {
		return string(chars[1 : len(chars)-1]), chars[0]
	}
	return lit, quoteRune
}
