package spnego

import (
	"net/http"
	"net/http/cookiejar"

	spnego_options "github.com/lublak/go-spnego/options"
)

func NewClient(base *http.Client, api spnego_options.Api, options spnego_options.Options) *http.Client {
	if base == nil {
		base = &http.Client{}
	}
	if base.Jar == nil {
		base.Jar, _ = cookiejar.New(nil)
	}

	base.Transport = NewRoundTripper(base.Transport, api, options)

	return base
}
