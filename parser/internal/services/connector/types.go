package connector

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type PollParams struct {
	Url     string
	Path    string
	Params  url.Values
	Headers http.Header
}

type RequestParams struct {
	Url     string
	Path    string
	Body    json.RawMessage
	Headers http.Header
}
