package spnego

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/lublak/go-spnego/options"
)

func NewClient(base *http.Client, api options.Api, spnegoOptions options.Options) *http.Client {
	if base == nil {
		base = &http.Client{}
	}
	if base.Jar == nil {
		base.Jar, _ = cookiejar.New(nil)
	}

	base.Transport = NewRoundTripper(base.Transport, api, spnegoOptions)

	return base
}
