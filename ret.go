package common

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
)

// Result represents a common-used result struct.
type Result struct {
	Code int         `json:"code"` // return code
	Msg  string      `json:"msg"`  // message
	Data interface{} `json:"data"` // data object
}

// NewResult creates a result with Code=0, Msg="", Data=nil.
func (*DuduRet) NewResult() *Result {
	return &Result{
		Code: 0,
		Msg:  "",
		Data: nil,
	}
}

// RetResult writes HTTP response with "Content-Type, application/json".
func (*DuduRet) RetResult(w http.ResponseWriter, r *http.Request, res *Result) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(res)
	if err != nil {
		return
	}

	w.Write(data)
}

// RetGzResult writes HTTP response with "Content-Type, application/json" and "Content-Encoding, gzip".
func (*DuduRet) RetGzResult(w http.ResponseWriter, r *http.Request, res *Result) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	err := json.NewEncoder(gz).Encode(res)
	if nil != err {

		return
	}

	err = gz.Close()
	if nil != err {

		return
	}
}

// RetJSON writes HTTP response with "Content-Type, application/json".
func (*DuduRet) RetJSON(w http.ResponseWriter, r *http.Request, res map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(res)
	if err != nil {
		return
	}

	w.Write(data)
}

// RetGzJSON writes HTTP response with "Content-Type, application/json" and "Content-Encoding, gzip".
func (*DuduRet) RetGzJSON(w http.ResponseWriter, r *http.Request, res map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	err := json.NewEncoder(gz).Encode(res)
	if nil != err {
		return
	}

	err = gz.Close()
	if nil != err {
		return
	}
}
