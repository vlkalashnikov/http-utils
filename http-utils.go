package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	method = strings.TrimSpace(strings.ToUpper(method))

	if headers == nil {
		headers = map[string]string{"Content-Type": "text/xml"}
	} else {
		headers["Content-Type"] = "text/xml"
	}

	httpStatus, responseBody, err = sendHttpReq(method, urlString, token, body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) != 0 {
		err = xml.Unmarshal(responseBody, responseStruct)
	}
	return
}

func HttpReqAuthJSON(method, urlString, token string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	method = strings.TrimSpace(strings.ToUpper(method))

	if headers == nil {
		headers = map[string]string{"Content-Type": "application/json"}
	} else {
		headers["Content-Type"] = "application/json"
	}

	httpStatus, responseBody, err = sendHttpReq(method, urlString, token, body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) != 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}
	return
}

func HttpReqXML(method, urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	method = strings.TrimSpace(strings.ToUpper(method))

	if headers == nil {
		headers = map[string]string{"Content-Type": "text/xml"}
	} else {
		headers["Content-Type"] = "text/xml"
	}

	httpStatus, responseBody, err = sendHttpReq(method, urlString, "", body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) > 0 {
		err = xml.Unmarshal(responseBody, responseStruct)
	}

	return
}

func HttpReqJSON(method, urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	method = strings.TrimSpace(strings.ToUpper(method))

	if headers == nil {
		headers = map[string]string{"Content-Type": "application/json"}
	} else {
		headers["Content-Type"] = "application/json"
	}

	httpStatus, responseBody, err = sendHttpReq(method, urlString, "", body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) > 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}

	return
}

func HttpReqPostFormJSON(urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	if headers == nil {
		headers = map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	} else {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	}

	httpStatus, responseBody, err = sendHttpReq("POST", urlString, "", body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}
	if responseStruct != nil && len(responseBody) != 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}
	return
}

func HttpReqPostFormXML(urlString string, body []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	if headers == nil {
		headers = map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	} else {
		if _, ok := headers["Content-Type"]; !ok {
			headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
	}

	httpStatus, responseBody, err = sendHttpReq("POST", urlString, "", body, headers, cookie, transport, timeout)
	if err != nil {
		return
	}
	if responseStruct != nil && len(responseBody) != 0 {
		err = xml.Unmarshal(responseBody, responseStruct)
	}
	return
}

func HttpReqPostFile(urlString string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range paramTexts {
		writer.WriteField(k, v)
	}

	fileWriter, err := writer.CreateFormFile(paramFile.Key, paramFile.FileName)
	if err != nil {
		return
	}

	fileWriter.Write(paramFile.Content)

	if headers == nil {
		headers = map[string]string{"Content-Type": writer.FormDataContentType()}
	} else {
		headers["Content-Type"] = writer.FormDataContentType()
	}

	writer.Close()

	httpStatus, responseBody, err = sendHttpReq("POST", urlString, "", body.Bytes(), headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) > 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}

	return
}

func HttpReqAuthPutFile(urlString, token string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return httpReqAuthFile("PUT", urlString, token, paramTexts, paramFile, headers, cookie, transport, timeout, responseStruct)
}

func HttpReqAuthPostFile(urlString, token string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	return httpReqAuthFile("POST", urlString, token, paramTexts, paramFile, headers, cookie, transport, timeout, responseStruct)
}

func httpReqAuthFile(method, urlString, token string, paramTexts map[string]string, paramFile FileItem, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range paramTexts {
		writer.WriteField(k, v)
	}

	fileWriter, err := writer.CreateFormFile(paramFile.Key, paramFile.FileName)
	if err != nil {
		return
	}

	fileWriter.Write(paramFile.Content)

	if headers == nil {
		headers = map[string]string{"Content-Type": writer.FormDataContentType()}
	} else {
		headers["Content-Type"] = writer.FormDataContentType()
	}

	writer.Close()

	httpStatus, responseBody, err = sendHttpReq(method, urlString, token, body.Bytes(), headers, cookie, transport, timeout)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) > 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}

	return
}

func sendHttpReq(method, urlString, token string, data []byte, headers map[string]string, cookie *http.Cookie, transport *http.Transport, timeout int) (httpStatus int, buf []byte, err error) {
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

	request, err := http.NewRequest(method, urlString, bytes.NewBuffer(data))

	if err != nil {
		return httpStatus, nil, &ResourceError{URL: urlString, Err: err}
	}

	if cookie != nil {
		request.AddCookie(cookie)
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	if token != "" {
		request.Header.Add("Authorization", token)
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

	buf, err = ioutil.ReadAll(response.Body)
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
