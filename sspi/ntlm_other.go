//go:build !windows
// +build !windows

package sspi

import (
	"net/http"

	"github.com/lublak/go-spnego/option"
)

func NewNtlmRoundTripper(base http.RoundTripper, options option.AuthOptions) http.RoundTripper {
	return nil
}
