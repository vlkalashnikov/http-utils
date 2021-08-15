package httputils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ResourceError struct {
	URL      string
	HTTPCode int
	Message  string
	Body     interface{}
	Err      error `json:"-"`
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

func HttpReqXML(method, urlString string, body []byte, headers map[string]string, cookie *http.Cookie, resultStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	if headers == nil {
		headers = map[string]string{"Content-Type": "text/xml"}
	} else {
		headers["Content-Type"] = "text/xml"
	}

	httpStatus, responseBody, err = sendHttpReq(method, urlString, body, headers, cookie)
	if err != nil {
		return
	}

	if resultStruct != nil && len(responseBody) > 0 {
		err = xml.Unmarshal(responseBody, resultStruct)
	}

	return
}

func HttpReqJSON(method, urlString string, body []byte, headers map[string]string, cookie *http.Cookie, resultStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	method = strings.TrimSpace(strings.ToUpper(method))

	if headers == nil {
		headers = map[string]string{"Content-Type": "application/json"}
	} else {
		headers["Content-Type"] = "application/json"
	}

	httpStatus, responseBody, err = sendHttpReq(method, urlString, body, headers, cookie)
	if err != nil {
		return
	}

	if resultStruct != nil && len(responseBody) > 0 {
		err = json.Unmarshal(responseBody, resultStruct)
	}

	return
}

func sendHttpReq(method, urlString string, data []byte, headers map[string]string, cookie *http.Cookie) (httpStatus int, buf []byte, err error) {
	client := &http.Client{}

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