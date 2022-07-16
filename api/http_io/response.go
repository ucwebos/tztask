package http_io

import (
	jsoniter "github.com/json-iterator/go"
	"tztask/utils/simple_server"
)

type HTTPResponse struct {
	Code  int32       `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data,omitempty"`
}

// JSONSuccess 成功返回
func JSONSuccess(ctx *simple_server.Context, data interface{}) {
	res := HTTPResponse{
		Data: data,
	}
	buf, _ := jsoniter.Marshal(res)
	_, _ = ctx.W.Write(buf)
}

// JSONError 失败返回
func JSONError(ctx *simple_server.Context, err Error) {
	res := HTTPResponse{
		Code:  err.Code,
		Error: err.Error(),
		Data:  struct{}{},
	}
	buf, _ := jsoniter.Marshal(res)
	_, _ = ctx.W.Write(buf)
}
