package connector

import "net/url"

type RequestParams struct {
	Url    string
	Path   string
	Params url.Values
}
