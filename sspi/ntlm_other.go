//go:build !windows
// +build !windows

package sspi

import (
	"net/http"

	"github.com/lublak/go-spnego/options"
)

func NewNtlmRoundTripper(base http.RoundTripper, spnegoOptions options.Options) http.RoundTripper {
	return nil
}
