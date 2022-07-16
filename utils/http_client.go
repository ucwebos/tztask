package utils

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	// HTTPNoKeepAliveClient is http client without keep alive
	HTTPNoKeepAliveClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	defaultHTTPClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 2048,
			IdleConnTimeout:     time.Minute * 5,
		},
	}
	defaultTimeout    = 500
	defaultRetryCount = 2
)

// PostRaw PostRaw
func PostRaw(client *http.Client, url string, header http.Header, reqBody interface{}, params ...int) ([]byte, *http.Response, error) {
	var (
		data []byte
		resp *http.Response
		err  error
	)
	timeout, retryCount := genDefaultParams(params...)
	for i := 0; i < retryCount; i++ {
		data, resp, err = do(client, http.MethodPost, url, header, reqBody, timeout)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Printf("PostRaw err: %v", err)
	}
	return data, resp, err
}

// PostWithUnmarshal do http get with unmarshal
func PostWithUnmarshal(client *http.Client, url string, header http.Header, reqBody interface{}, resp interface{}, params ...int) error {
	data, _, err := PostRaw(client, url, header, reqBody, params...)
	if err != nil {
		return err
	}
	// for no resp needed request.
	if resp == nil {
		return nil
	}
	// for big int
	decoder := jsoniter.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	err = decoder.Decode(resp)
	if err != nil {
		log.Printf("PostWithUnmarshal.Decode err: %v", err)
	}
	return err
}

// GetRaw get http raw
func GetRaw(client *http.Client, url string, header http.Header, reqBody interface{}, params ...int) ([]byte, *http.Response, error) {
	var (
		data []byte
		resp *http.Response
		err  error
	)
	timeout, retryCount := genDefaultParams(params...)
	for i := 0; i < retryCount; i++ {
		data, resp, err = do(client, http.MethodGet, url, header, reqBody, timeout)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Printf("GetRaw err: %v", err)
	}
	return data, resp, err
}

// GetWithUnmarshal do http get with unmarshal
func GetWithUnmarshal(client *http.Client, url string, header http.Header, reqBody interface{}, resp interface{}, params ...int) error {
	data, _, err := GetRaw(client, url, header, reqBody, params...)
	if err != nil {
		return err
	}
	// for no resp needed request.
	if resp == nil {
		return nil
	}
	decoder := jsoniter.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	err = decoder.Decode(resp)
	if err != nil {
		log.Printf("GetWithUnmarshal.Decode err: %v", err)
	}
	return err
}

func genDefaultParams(params ...int) (int, int) {
	timeout, retryCount := defaultTimeout, defaultRetryCount
	switch {
	case len(params) >= 2:
		timeout, retryCount = params[0], params[1]
	case len(params) >= 1:
		timeout = params[0]
	}
	return timeout, retryCount
}

func do(client *http.Client, method string, url string, header http.Header, reqBody interface{}, timeout int) ([]byte, *http.Response, error) {
	if client == nil {
		client = defaultHTTPClient
	}
	var reader io.Reader
	switch v := reqBody.(type) {
	case nurl.Values:
		reader = strings.NewReader(v.Encode())
	case []byte:
		reader = bytes.NewBuffer(v)
	case string:
		reader = strings.NewReader(v)
	case io.Reader:
		reader = v
	default:
		buff := &bytes.Buffer{}
		err := jsoniter.NewEncoder(buff).Encode(v)
		if err != nil {
			return nil, nil, err
		}
		reader = buff
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, nil, err
	}
	if header != nil {
		req.Header = header
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeout))
	defer cancelFunc()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	return data, resp, nil
}
