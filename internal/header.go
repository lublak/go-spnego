package internal

import (
	"encoding/base64"
	"errors"
	"strings"
)

func HasBasic(headers []string) bool {
	for _, header := range headers {
		if strings.HasPrefix(header, "Basic") {
			return true
		}
	}
	return false
}

func HasNegotiate(headers []string) bool {
	for _, header := range headers {
		if strings.HasPrefix(header, "Negotiate") {
			return true
		}
	}
	return false
}

func HasNtlm(headers []string) bool {
	for _, header := range headers {
		if strings.HasPrefix(header, "NTLM") {
			return true
		}
	}
	return false
}

func FindAndDecodeNegotiateToken(headers []string) ([]byte, error) {
	for _, header := range headers {
		if strings.HasPrefix(header, "Negotiate ") {
			firstSpaceIndex := strings.IndexRune(header, ' ')
			if firstSpaceIndex < 0 {
				return nil, errors.New("invalid negotiate auth header")
			}

			token, err := base64.StdEncoding.DecodeString(header[len("Negotiate "):])
			if err != nil {
				return nil, err
			}

			return token, nil
		}
	}
	return nil, errors.New("empty negotiate auth header")
}

func EncodeNegotiateToken(token []byte) string {
	return "Negotiate " + base64.StdEncoding.EncodeToString(token)
}

func FindAndDecodeNtlmToken(headers []string) ([]byte, error) {
	for _, header := range headers {
		if strings.HasPrefix(header, "NTLM ") {
			firstSpaceIndex := strings.IndexRune(header, ' ')
			if firstSpaceIndex < 0 {
				return nil, errors.New("invalid NTLM auth header")
			}

			token, err := base64.StdEncoding.DecodeString(header[len("NTLM "):])
			if err != nil {
				return nil, err
			}

			return token, nil
		}
	}

	return nil, errors.New("empty negotiate auth header")
}

func EncodeNtlmToken(token []byte) string {
	return "NTLM " + base64.StdEncoding.EncodeToString(token)
}
