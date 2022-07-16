package http_io

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"tztask/utils/simple_server"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func BindBody(ctx *simple_server.Context, obj interface{}) error {
	var (
		raw []byte
		err error
	)
	switch ctx.R.Method {
	case http.MethodGet:
		return errors.Errorf("method[%s] not allowd BindBody", http.MethodGet)
	}
	if ctx.R.Body != nil {
		raw, err = ioutil.ReadAll(ctx.R.Body)
		if err != nil {
			log.Printf("r.GetBody ReadAll err: %v", err)
			return err
		}
	}
	return decodeJSON(bytes.NewReader(raw), obj)
}

func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := JSON.NewDecoder(r)
	decoder.UseNumber()
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
