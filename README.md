# http-utils

Lightweight HTTP helpers on top of `net/http` for JSON/XML requests and multipart file uploads.

## Installation

```bash
go get github.com/vlkalashnikov/http-utils
```

```go
import utils "github.com/vlkalashnikov/http-utils"
```

## API Overview

The package provides two API styles:

- Legacy methods without context, for backward compatibility.
- Context-aware methods with `Ctx` suffix, recommended for new code.

Available pairs:

- `HttpReqAuthXML` / `HttpReqAuthXMLCtx`
- `HttpReqAuthJSON` / `HttpReqAuthJSONCtx`
- `HttpReqXML` / `HttpReqXMLCtx`
- `HttpReqJSON` / `HttpReqJSONCtx`
- `HttpReqPostFormJSON` / `HttpReqPostFormJSONCtx`
- `HttpReqPostFormXML` / `HttpReqPostFormXMLCtx`
- `HttpReqPostFile` / `HttpReqPostFileCtx`
- `HttpReqAuthPutFile` / `HttpReqAuthPutFileCtx`
- `HttpReqAuthPostFile` / `HttpReqAuthPostFileCtx`

## Common Arguments

- `headers map[string]string`: additional request headers.
- `cookie *http.Cookie`: optional cookie, pass `nil` if not needed.
- `transport *http.Transport`: custom transport, pass `nil` for default behavior.
- `timeout int`: timeout in seconds. `0` means default `30s`.
- `responseStruct interface{}`: target for JSON/XML decode. Pass `nil` to skip decode.

All methods return:

- `httpStatus int`
- `responseBody []byte` (raw response body)
- `err error`

## Example: JSON Request

```go
package main

import (
	"log"

	utils "github.com/vlkalashnikov/http-utils"
)

type PingResponse struct {
	Message string `json:"message"`
}

func main() {
	var out PingResponse

	status, raw, err := utils.HttpReqJSON(
		"GET",
		"https://example.com/api/ping",
		nil,
		nil,
		nil,
		nil,
		10,
		&out,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%d body=%s parsed=%+v", status, string(raw), out)
}
```

## Example: Context-Aware Request

```go
package main

import (
	"context"
	"log"
	"time"

	utils "github.com/vlkalashnikov/http-utils"
)

type CreateResponse struct {
	ID string `json:"id"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	payload := []byte(`{"name":"demo"}`)
	token := "Bearer your-token"
	var out CreateResponse

	status, raw, err := utils.HttpReqAuthJSONCtx(
		ctx,
		"POST",
		"https://example.com/api/items",
		token,
		payload,
		nil,
		nil,
		nil,
		0,
		&out,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%d body=%s id=%s", status, string(raw), out.ID)
}
```

## Example: File Upload

```go
package main

import (
	"context"
	"log"
	"time"

	utils "github.com/vlkalashnikov/http-utils"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file := utils.FileItem{
		Key:      "file",
		FileName: "report.txt",
		Content:  []byte("hello"),
	}

	status, raw, err := utils.HttpReqPostFileCtx(
		ctx,
		"https://example.com/upload",
		map[string]string{"folder": "daily"},
		file,
		nil,
		nil,
		nil,
		30,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%d body=%s", status, string(raw))
}
```

## Error Handling

For HTTP status `>= 400`, methods return `*ResourceError`.

```go
var re *utils.ResourceError
if errors.As(err, &re) {
	log.Printf("url=%s status=%d err=%v", re.URL, re.HTTPCode, re.Err)
}
```

## Behavior Notes

- JSON/XML helpers set `Content-Type` automatically.
- Header maps are cloned internally and are not mutated in-place.
- `Authorization` header is set when `token` is not empty.
