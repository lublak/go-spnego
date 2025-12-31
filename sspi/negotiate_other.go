//go:build !windows
// +build !windows

package sspi

import (
	"net/http"

	"github.com/lublak/go-spnego/option"
)

func NewNegotiateRoundTripper(base http.RoundTripper, options option.AuthOptions) http.RoundTripper {
	return nil
}
