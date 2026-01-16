//go:build !windows
// +build !windows

package sspi

import (
	"net/http"

	spnego_options "github.com/lublak/go-spnego/options"
)

func NewNtlmRoundTripper(base http.RoundTripper, options spnego_options.Options) http.RoundTripper {
	return nil
}
