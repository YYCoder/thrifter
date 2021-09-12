package thrifter

import (
	"strings"
)

type token int

const baseTypeTokens = "bool byte i8 i16 i32 i64 double string binary slist"

const (
	// special tokens
	tILLEGAL token = iota
	tEOF
	tIDENT

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
		if isComment(literal) {
			return tCOMMENT
		}
		return tIDENT
	}
}

func isComment(lit string) bool {
	return strings.HasPrefix(lit, "//") || strings.HasPrefix(lit, "/*")
}

// isDigit returns true if the rune is a digit.
func isDigit(lit rune) bool {
	return (lit >= '0' && lit <= '9')
}

// isKeyword returns if tok is in the keywords range
func isKeyword(tok token) bool {
	return keywordStart < tok && tok < keywordEnd
}

func isWhitespace(tok token) bool {
	return tok == tSPACE || tok == tLINEBREAK || tok == tRETURN || tok == tTAB
}

const doubleQuoteRune = rune('"')
