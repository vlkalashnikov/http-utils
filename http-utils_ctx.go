package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func cloneHeaders(headers map[string]string) map[string]string {
	if len(headers) == 0 {
		return map[string]string{}
	}

	cloned := make(map[string]string, len(headers))
	for key, value := range headers {
		cloned[key] = value
	}

	return cloned
}

func headersWithContentType(headers map[string]string, contentType string, overwrite bool) map[string]string {
	cloned := cloneHeaders(headers)

	if overwrite {
		cloned["Content-Type"] = contentType
		return cloned
	}

	if _, ok := cloned["Content-Type"]; !ok {
		cloned["Content-Type"] = contentType
	}

	return cloned
}

func HttpReqAuthXMLCtx(ctx context.Context, method, urlString, token string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	method = strings.TrimSpace(strings.ToUpper(method))
	headers = headersWithContentType(headers, "text/xml", true)

	httpStatus, responseBody, err = sendHttpReqCtx(ctx, method, urlString, token, body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) != 0 {
		err = xml.Unmarshal(responseBody, responseStruct)
	}
	return
}

func HttpReqAuthJSONCtx(ctx context.Context, method, urlString, token string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	method = strings.TrimSpace(strings.ToUpper(method))
	headers = headersWithContentType(headers, "application/json", true)

	httpStatus, responseBody, err = sendHttpReqCtx(ctx, method, urlString, token, body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) != 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}
	return
}

func HttpReqXMLCtx(ctx context.Context, method, urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	method = strings.TrimSpace(strings.ToUpper(method))
	headers = headersWithContentType(headers, "text/xml", true)

	httpStatus, responseBody, err = sendHttpReqCtx(ctx, method, urlString, "", body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) > 0 {
		err = xml.Unmarshal(responseBody, responseStruct)
	}

	return
}

func HttpReqJSONCtx(ctx context.Context, method, urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	method = strings.TrimSpace(strings.ToUpper(method))
	headers = headersWithContentType(headers, "application/json", true)

	httpStatus, responseBody, err = sendHttpReqCtx(ctx, method, urlString, "", body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) > 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}

	return
}

func HttpReqPostFormJSONCtx(ctx context.Context, urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	headers = headersWithContentType(headers, "application/x-www-form-urlencoded", true)

	httpStatus, responseBody, err = sendHttpReqCtx(ctx, "POST", urlString, "", body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}
	if responseStruct != nil && len(responseBody) != 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}
	return
}

func HttpReqPostFormXMLCtx(ctx context.Context, urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	headers = headersWithContentType(headers, "application/x-www-form-urlencoded", false)

	httpStatus, responseBody, err = sendHttpReqCtx(ctx, "POST", urlString, "", body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}
	if responseStruct != nil && len(responseBody) != 0 {
		err = xml.Unmarshal(responseBody, responseStruct)
	}
	return
}

func HttpReqPostFileCtx(ctx context.Context, urlString string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range paramTexts {
		if err = writer.WriteField(k, v); err != nil {
			return
		}
	}

	fileWriter, err := writer.CreateFormFile(paramFile.Key, paramFile.FileName)
	if err != nil {
		return
	}

	if _, err = fileWriter.Write(paramFile.Content); err != nil {
		return
	}

	headers = headersWithContentType(headers, writer.FormDataContentType(), true)

	if err = writer.Close(); err != nil {
		return
	}

	httpStatus, responseBody, err = sendHttpReqCtx(ctx, "POST", urlString, "", body.Bytes(), headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) > 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}

	return
}

func HttpReqAuthPutFileCtx(ctx context.Context, urlString, token string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return httpReqAuthFileCtx(ctx, "PUT", urlString, token, paramTexts, paramFile, headers, cookie, transport, timeout, responseStruct)
}

func HttpReqAuthPostFileCtx(ctx context.Context, urlString, token string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return httpReqAuthFileCtx(ctx, "POST", urlString, token, paramTexts, paramFile, headers, cookie, transport, timeout, responseStruct)
}

func httpReqAuthFileCtx(ctx context.Context, method, urlString, token string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range paramTexts {
		if err = writer.WriteField(k, v); err != nil {
			return
		}
	}

	fileWriter, err := writer.CreateFormFile(paramFile.Key, paramFile.FileName)
	if err != nil {
		return
	}

	if _, err = fileWriter.Write(paramFile.Content); err != nil {
		return
	}

	headers = headersWithContentType(headers, writer.FormDataContentType(), true)

	if err = writer.Close(); err != nil {
		return
	}

	httpStatus, responseBody, err = sendHttpReqCtx(ctx, method, urlString, token, body.Bytes(), headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) > 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}

	return
}

func sendHttpReqCtx(ctx context.Context, method, urlString, token string, data []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int) (httpStatus int, buf []byte, err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	defaultTimeout := 30 * time.Second //default timeout

	if timeout > 0 {
		defaultTimeout = time.Duration(timeout) * time.Second
	}

	client := &http.Client{
		Timeout: defaultTimeout,
	}

	if transport != nil {
		client.Transport = transport
	}

	request, err := http.NewRequestWithContext(ctx, method, urlString, bytes.NewBuffer(data))

	if err != nil {
		return httpStatus, nil, &ResourceError{URL: urlString, Err: err}
	}

	if cookie != nil {
		request.AddCookie(cookie)
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	if token != "" {
		request.Header.Set("Authorization", token)
	}

	if strings.ContainsAny(urlString, "?") {
		urlTemp, err := url.Parse(urlString)
		if err != nil {
			return httpStatus, nil, &ResourceError{URL: urlString, Err: err}
		}

		urlQuery := urlTemp.Query()
		urlTemp.RawQuery = urlQuery.Encode()
		urlString = urlTemp.String()
	}

	response, err := client.Do(request)
	if err != nil {
		return httpStatus, nil, &ResourceError{URL: urlString, Err: err}
	}
	defer response.Body.Close()

	buf, err = io.ReadAll(response.Body)
	if err != nil {
		return httpStatus, nil, &ResourceError{URL: urlString, Err: err, HTTPCode: response.StatusCode}
	}

	httpStatus = response.StatusCode
	if response.StatusCode > 399 {
		return httpStatus, buf, &ResourceError{
			URL:      urlString,
			Err:      fmt.Errorf("incorrect status code"),
			HTTPCode: response.StatusCode,
			Message:  "incorrect response.StatusCode",
			Body:     string(data),
		}
	}

	return
}
