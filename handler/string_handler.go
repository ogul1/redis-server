package handler

import (
	"errors"
	"redis-server/parser"
	"strconv"
	"strings"
)

var sdb = make(map[string]string)

var StringCommands = map[string]bool{
	"command":  true,
	"ping":     true,
	"set":      true,
	"get":      true,
	"append":   true,
	"strlen":   true,
	"getrange": true,
	"decr":     true,
	"decrby":   true,
}

func HandleString(request parser.Value) []byte {
	if request.Typ != parser.Array {
		return []byte("-Invalid Request\r\n")
	}
	switch strings.ToLower(request.Array[0].Str) {
	case "command":
		return []byte("+OK\r\n")
	case "ping":
		return []byte("+PONG\r\n")
	case "set":
		if len(request.Array) != 3 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleSet(request.Array[1].Str, request.Array[2].Str)
	case "get":
		if len(request.Array) != 2 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleGet(request.Array[1].Str)
	case "append":
		if len(request.Array) != 3 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleAppend(request.Array[1].Str, request.Array[2].Str)
	case "strlen":
		if len(request.Array) != 2 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleStrLen(request.Array[1].Str)
	case "getrange":
		if len(request.Array) != 4 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleGetRange(request.Array[1].Str, request.Array[2].Str, request.Array[3].Str)
	case "decr":
		if len(request.Array) != 2 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleDecr(request.Array[1].Str)
	case "decrby":
		if len(request.Array) != 3 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleDecrBy(request.Array[1].Str, request.Array[2].Str)
	default:
		return ErrResp(errors.New("unsupported command"))
	}
}

func handleSet(key, val string) []byte {
	sdb[key] = val
	return []byte("+OK\r\n")
}

func handleGet(key string) []byte {
	if val, ok := sdb[key]; ok {
		return StringResp(val)
	}
	return NilResp()
}

func handleAppend(key, add string) []byte {
	if val, ok := sdb[key]; ok {
		sdb[key] = val + add
	} else {
		sdb[key] = add
	}
	return IntResp(len(sdb[key]))
}

func handleStrLen(key string) []byte {
	if val, ok := sdb[key]; ok {
		return IntResp(len(val))
	}
	return IntResp(0)
}

func handleGetRange(key, _start, _end string) []byte {
	start, _ := strconv.Atoi(_start)
	end, _ := strconv.Atoi(_end)
	if val, ok := sdb[key]; ok {
		n := len(val)
		if n == 0 {
			return StringResp("")
		}
		if start < 0 {
			start += n
		}
		if end < 0 {
			end += n
		}
		if start < 0 {
			start = 0
		}
		if end < 0 {
			end = -1
		}
		if end >= n {
			end = n - 1
		}
		if start > end {
			return StringResp("")
		}
		return StringResp(val[start : end+1])
	}
	return StringResp("")
}

func handleDecr(key string) []byte {
	if _, ok := sdb[key]; !ok {
		sdb[key] = "0"
	}
	val, err := strconv.Atoi(sdb[key])
	if err != nil {
		return ErrResp(err)
	}
	sdb[key] = strconv.Itoa(val - 1)
	return IntResp(val - 1)
}

func handleDecrBy(key, _amount string) []byte {
	if _, ok := sdb[key]; !ok {
		sdb[key] = "0"
	}
	amount, err := strconv.Atoi(_amount)
	if err != nil {
		return ErrResp(err)
	}
	val, err := strconv.Atoi(sdb[key])
	if err != nil {
		return ErrResp(err)
	}
	sdb[key] = strconv.Itoa(val - amount)
	return IntResp(val - amount)
}
