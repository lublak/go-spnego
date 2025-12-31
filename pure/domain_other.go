//go:build !windows
// +build !windows

package pure

func getDomain() (string, error) {
	return "", nil
}
