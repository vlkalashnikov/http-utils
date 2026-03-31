package utils

import (
	"context"
	"fmt"
	"net/http"
)

type ResourceError struct {
	URL      string
	HTTPCode int
	Message  string
	Body     interface{}
	Err      error `json:"-"`
}

type FileItem struct {
	Key      string
	FileName string
	Content  []byte
}

func (re *ResourceError) Error() string {
	return fmt.Sprintf(
		"Resource error: URL: %s, status code: %v,  err: %v, body: %v",
		re.URL,
		re.HTTPCode,
		re.Err,
		re.Body,
	)
}

func HttpReqAuthXML(method, urlString, token string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return HttpReqAuthXMLCtx(context.Background(), method, urlString, token, body, headers, cookie, transport, timeout, responseStruct)
}

func HttpReqAuthJSON(method, urlString, token string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return HttpReqAuthJSONCtx(context.Background(), method, urlString, token, body, headers, cookie, transport, timeout, responseStruct)
}

func HttpReqXML(method, urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return HttpReqXMLCtx(context.Background(), method, urlString, body, headers, cookie, transport, timeout, responseStruct)
}

func HttpReqJSON(method, urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return HttpReqJSONCtx(context.Background(), method, urlString, body, headers, cookie, transport, timeout, responseStruct)
}

func HttpReqPostFormJSON(urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return HttpReqPostFormJSONCtx(context.Background(), urlString, body, headers, cookie, transport, timeout, responseStruct)
}

func HttpReqPostFormXML(urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return HttpReqPostFormXMLCtx(context.Background(), urlString, body, headers, cookie, transport, timeout, responseStruct)
}

func HttpReqPostFile(urlString string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return HttpReqPostFileCtx(context.Background(), urlString, paramTexts, paramFile, headers, cookie, transport, timeout, responseStruct)
}

func HttpReqAuthPutFile(urlString, token string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return HttpReqAuthPutFileCtx(context.Background(), urlString, token, paramTexts, paramFile, headers, cookie, transport, timeout, responseStruct)
}

func HttpReqAuthPostFile(urlString, token string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return HttpReqAuthPostFileCtx(context.Background(), urlString, token, paramTexts, paramFile, headers, cookie, transport, timeout, responseStruct)
}
