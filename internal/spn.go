package internal

import (
	"net"
	"net/http"
	"strings"
)

func GetSpnFromRequest(req *http.Request) string {
	spn := req.Host

	if len(spn) == 0 {
		spn = req.URL.Host
	}

	host, _, err := net.SplitHostPort(spn)
	if err == nil {
		spn = host
	}

	cname, err := net.LookupCNAME(spn)
	if err == nil && len(cname) > 0 {
		spn = strings.ToLower(cname)
	}

	spn = strings.TrimSuffix(spn, ".")

	return "HTTP/" + spn
}
