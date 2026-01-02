//go:build !windows
// +build !windows

package internal

func GetDomain() (string, error) {
	return "", nil
}
