package httpclient

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	KDefaultTimeout  = 10
	KDefaultTryCount = 1
	KErrCode         = -1

	MinTimeoutLimited  = 1
	MaxTryCountLimited = 10
	MinTryCountLimited = 1
)

type filterType int

const (
	NormalType filterType = 0
	GzipType   filterType = 1
)

const (
	gzipStr          = "gzip"
	authBeeLibHeader = "beelib"
	defaultBeeLibKey = "7aa766df8c3f03822ca0077e26bfc983"
)

type httpParams struct {
	timeout    int
	tryCount   int
	headerInfo map[string]string
	otherInfo  interface{}
}

func NewHttpParams() *httpParams {
	return &httpParams{timeout: KDefaultTimeout, tryCount: KDefaultTryCount}
}

func (hp *httpParams) AddHeader(key, value string) *httpParams {
	if hp.headerInfo == nil {
		hp.headerInfo = make(map[string]string)
	}

	hp.headerInfo[key] = value
	return hp
}

func (hp *httpParams) AddMultiHeaders(headers map[string]string) *httpParams {
	for key, value := range headers {
		hp.AddHeader(key, value)
	}
	return hp
}

func (hp *httpParams) SetTimeout(timeout int) *httpParams {

	if timeout < MinTimeoutLimited {
		timeout = KDefaultTimeout
	}

	hp.timeout = timeout
	return hp
}

func (hp *httpParams) SetMaxTryCount(count int) *httpParams {

	switch {
	case count > MaxTryCountLimited:
		hp.tryCount = MaxTryCountLimited
	case count < MinTimeoutLimited:
		hp.tryCount = MinTimeoutLimited
	default:
		hp.tryCount = count
	}
	return hp
}

func buildClient(timeout int) *http.Client {

	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		},
	}

	return c
}

func buildHttpRequest(method, url string, hp *httpParams, body []byte) (*http.Request, error) {

	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if hp != nil {
		for key, value := range hp.headerInfo {
			request.Header.Add(key, value)
		}
	}

	return request, nil
}

func XGet(url string, hp *httpParams) (int, []byte, error) {
	return xRequestDo(http.MethodGet, url, hp, nil)
}

func XGetAndCarryBody(url string, hp *httpParams, body []byte) (int, []byte, error) {
	return xRequestDo(http.MethodGet, url, hp, body)
}

func XPost(url string, hp *httpParams, body []byte) (int, []byte, error) {
	return xRequestDo(http.MethodPost, url, hp, body)
}

func xRequestDo(requestType, url string, hp *httpParams, body []byte) (int, []byte, error) {

	var (
		code int
		data []byte
		err  error
	)

	if hp == nil {
		hp = NewHttpParams()
	}

	hp.AddHeader(authBeeLibHeader, defaultBeeLibKey)
	for hp.tryCount > 0 {

		code, data, err = httpRequestImpl(requestType, url, hp, body)
		if err != nil {
			hp.tryCount--
			continue
		}

		return code, data, err

	}

	return KErrCode, nil, errors.New(fmt.Sprintf("retry all fail ,error:%v", err))
}

func httpRequestImpl(requestType, url string, hp *httpParams, body []byte) (int, []byte, error) {

	request, err := buildHttpRequest(requestType, url, hp, body)
	if err != nil {
		return KErrCode, nil, err
	}

	client := buildClient(hp.timeout)
	res, err := client.Do(request)
	if err != nil {
		return KErrCode, nil, err
	}
	defer res.Body.Close()

	cEncoding := res.Header.Get("Content-Encoding")
	ftype := scanReponseType(cEncoding)

	switch ftype {

	case GzipType:
		compressedReader, err := gzip.NewReader(res.Body)
		if err != nil {
			return res.StatusCode, nil, err
		}
		defer compressedReader.Close()

		data, err := ioutil.ReadAll(compressedReader)
		if err != nil {
			return KErrCode, nil, err
		}

		return res.StatusCode, data, nil
	default:
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return KErrCode, nil, err
	}

	return res.StatusCode, data, nil
}

func scanReponseType(contentEncoding string) filterType {

	switch contentEncoding {
	case gzipStr:
		return GzipType
	default:
		return NormalType
	}

}
