package parser

import (
	"bytes"
	"log"
)

var CRLF = []byte("\r\n")

type Type byte

const (
	Array        Type = '*'
	Integer      Type = ':'
	SimpleString Type = '+'
	BulkString   Type = '$'
	Error        Type = '-'
)

type Value struct {
	Typ     Type
	Integer int
	Str     string
	Array   []Value
	Nil     bool
}

func Parse(request []byte) (result Value, rest []byte) {
	data, rest, found := bytes.Cut(request, CRLF)
	if !found {
		log.Fatal("Invalid RESP expression.")
	}
	if data[0] == '*' { // Array
		n := parseInt(data[1:])
		var arr Value
		for i := 0; i < n; i++ {
			var res Value
			res, rest = Parse(rest)
			arr.Array = append(arr.Array, res)
		}
		arr.Typ = Array
		return arr, nil
	} else if data[0] == '+' { // Simple String
		result.Typ = SimpleString
		result.Str = string(data[1:])
		return result, rest
	} else if data[0] == ':' { // Integer
		result.Typ = Integer
		result.Integer = parseInt(data[1:])
		return result, rest
	} else if data[0] == '$' { // Bulk String
		if data[1] == '-' {
			result.Nil = true
			return result, rest
		}
		str, rest, found := bytes.Cut(rest, CRLF)
		if !found {
			log.Fatal("Invalid RESP expression.")
		}
		result.Typ = BulkString
		result.Str = string(str)
		return result, rest
	} else if data[0] == '-' { // Error
		result.Typ = Error
		result.Str = string(data[1:])
		return result, rest
	} else {
		log.Fatal("Invalid RESP expression.")
	}
	return Value{}, nil
}

func parseInt(data []byte) (res int) {
	for _, c := range data {
		res = 10*res + int(c-'0')
	}
	return
}
