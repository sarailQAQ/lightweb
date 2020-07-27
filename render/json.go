package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type JSON struct {
	Data interface{}
}

type AsciiJSON struct {
	Data interface{}
}

var jsonAsciiContentType = []string{"application/json"}
var jsonContentType = []string{"application/json; charset=utf-8"}


func (r JSON) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (r AsciiJSON) Render(w  http.ResponseWriter) (err error) {
	writeContentType(w,jsonAsciiContentType)
	ret,err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for _,r := range string(ret) {
		cvt := string(r)
		if r >= 128 {
			cvt = fmt.Sprintf("\\u%04x", int64(r))
		}
		buffer.WriteString(cvt)
	}

	_,err = w.Write(buffer.Bytes())
	return
}

func (r AsciiJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonAsciiContentType)
}
