package spnego

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/lublak/go-spnego/option"
)

func NewClient(base *http.Client, api option.ApiType, options option.AuthOptions) *http.Client {
	if base == nil {
		base = &http.Client{}
	}
	if base.Jar == nil {
		base.Jar, _ = cookiejar.New(nil)
	}

	base.Transport = NewRoundTripper(base.Transport, api, options)

	return base
}
