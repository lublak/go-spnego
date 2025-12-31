package pure

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func getDomain() (string, error) {
	var domain *uint16
	var status uint32
	err := windows.NetGetJoinInformation(nil, &domain, &status)
	if err != nil {
		return "", err
	}
	defer windows.NetApiBufferFree((*byte)(unsafe.Pointer(domain)))

	if status == windows.NetSetupDomainName {
		return windows.UTF16PtrToString(domain), nil
	}

	return "", nil
}
