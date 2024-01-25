package handler

import (
	"bytes"
	"strconv"
)

func StringResp(data string) []byte {
	response := []byte("+\"")
	response = append(response, data...)
	response = append(response, "\"\r\n"...)
	return response
}

func NilResp() []byte {
	return []byte("$-1\r\n")
}

func IntResp(n int) []byte {
	response := []byte(":")
	response = append(response, strconv.Itoa(n)...)
	response = append(response, "\r\n"...)
	return response
}

func ErrResp(err error) []byte {
	response := []byte("-")
	response = append(response, err.Error()...)
	response = append(response, "\r\n"...)
	return response
}

func ArrayResp(arr []string) []byte {
	var buf bytes.Buffer
	buf.WriteByte('*')
	buf.WriteString(strconv.Itoa(len(arr)))
	buf.WriteString("\r\n")
	for _, elem := range arr {
		buf.WriteByte('$')
		buf.WriteString(strconv.Itoa(len(elem)))
		buf.WriteString("\r\n")
		buf.WriteString(elem)
		buf.WriteString("\r\n")
	}
	return buf.Bytes()
}
