package handler

import (
	"errors"
	"redis-server/parser"
	"strconv"
	"strings"
)

var ldb = make(map[string][]string)

var ListCommands = map[string]bool{
	"lpush":  true,
	"lrange": true,
	"llen":   true,
	"lindex": true,
	"lset":   true,
	"lrem":   true,
	"lpos":   true,
	"rpush":  true,
}

func HandleList(request parser.Value) []byte {
	if request.Typ != parser.Array {
		return []byte("-Invalid Request\r\n")
	}
	switch strings.ToLower(request.Array[0].Str) {
	case "lpush":
		return handleLpush(request.Array[1].Str, request.Array[2:])
	case "lrange":
		if len(request.Array) != 4 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleLrange(request.Array[1].Str, request.Array[2].Str, request.Array[3].Str)
	case "llen":
		if len(request.Array) != 2 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleLlen(request.Array[1].Str)
	case "lindex":
		if len(request.Array) != 3 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleLindex(request.Array[1].Str, request.Array[2].Str)
	case "lset":
		if len(request.Array) != 4 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleLset(request.Array[1].Str, request.Array[2].Str, request.Array[3].Str)
	case "lrem":
		if len(request.Array) != 4 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleLrem(request.Array[1].Str, request.Array[2].Str, request.Array[3].Str)
	case "lpos":
		if len(request.Array) != 3 {
			return ErrResp(errors.New("invalid number of arguments"))
		}
		return handleLpos(request.Array[1].Str, request.Array[2].Str)
	case "rpush":
		return handleRpush(request.Array[1].Str, request.Array[2:])
	default:
		return ErrResp(errors.New("unsupported command"))
	}
}

func handleLpush(name string, values []parser.Value) []byte {
	if _, ok := ldb[name]; !ok {
		ldb[name] = make([]string, 0)
	}
	for _, val := range values {
		ldb[name] = append([]string{val.Str}, ldb[name]...)
	}
	return IntResp(len(ldb[name]))
}

func handleRpush(name string, values []parser.Value) []byte {
	if _, ok := ldb[name]; !ok {
		ldb[name] = make([]string, 0)
	}
	for _, val := range values {
		ldb[name] = append(ldb[name], val.Str)
	}
	return IntResp(len(ldb[name]))
}

func handleLrange(name, _start, _end string) []byte {
	start, _ := strconv.Atoi(_start)
	end, _ := strconv.Atoi(_end)
	if val, ok := ldb[name]; ok {
		n := len(val)
		if n == 0 {
			return ArrayResp([]string{})
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
			return ArrayResp([]string{})
		}
		return ArrayResp(val[start : end+1])
	}
	return ArrayResp([]string{})
}

func handleLlen(name string) []byte {
	if _, ok := ldb[name]; ok {
		return IntResp(len(ldb[name]))
	}
	return IntResp(0)
}

func handleLindex(name, _index string) []byte {
	index, err := strconv.Atoi(_index)
	if err != nil {
		return ErrResp(errors.New("invalid index"))
	}
	if _, ok := ldb[name]; ok {
		n := len(ldb[name])
		if index < 0 {
			index += n
		}
		if index < 0 || index >= n {
			return NilResp()
		}
		return StringResp(ldb[name][index])
	}
	return NilResp()
}

func handleLset(name, _index, element string) []byte {
	index, err := strconv.Atoi(_index)
	if err != nil {
		return ErrResp(errors.New("invalid index"))
	}
	if _, ok := ldb[name]; ok {
		n := len(ldb[name])
		if index < 0 {
			index += n
		}
		if index < 0 || index >= n {
			return ErrResp(errors.New("invalid index"))
		}
		ldb[name][index] = element
		return StringResp("OK")
	}
	return ErrResp(errors.New("list does not exist"))
}

func handleLpos(name, want string) []byte {
	if _, ok := ldb[name]; ok {
		for i, str := range ldb[name] {
			if str == want {
				return IntResp(i)
			}
		}
	}
	return NilResp()
}

func handleLrem(name, _count, element string) []byte {
	count, err := strconv.Atoi(_count)
	if err != nil {
		return ErrResp(errors.New("invalid count"))
	}
	if _, ok := ldb[name]; ok {
		n := len(ldb[name])
		if count == 0 {
			count = n
		}
		filtered := make([]string, 0)
		if count < 0 {
			for i := n - 1; i >= 0; i-- {
				if ldb[name][i] == element && count < 0 {
					count++
					continue
				}
				filtered = append([]string{ldb[name][i]}, filtered...)
			}
		} else {
			for _, str := range ldb[name] {
				if str == element && count > 0 {
					count--
					continue
				}
				filtered = append(filtered, str)
			}
		}
		ldb[name] = filtered
		return IntResp(n - len(ldb[name]))
	}
	return ErrResp(errors.New("list does not exist"))
}
