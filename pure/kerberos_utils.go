package pure

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/credentials"
)

func getDefaultConfigPath() string {
	switch runtime.GOOS {
	case "linux":
		return "/etc/krb5.conf"
	case "darwin":
		return "/opt/local/etc/krb5.conf"
	case "windows":
		return "c:\\windows\\krb5.ini"
	default:
		return "/etc/krb5/krb5.conf"
	}
}

func defaultKerberosConfig() (*config.Config, error) {
	configPath := os.Getenv("KRB5_CONFIG")

	if len(configPath) == 0 {
		return kerberosConfigFromPath(getDefaultConfigPath())
	}

	config, err := kerberosConfigFromPath(configPath)

	if err != nil {
		if os.IsNotExist(err) {
			return kerberosConfigFromPath(getDefaultConfigPath())
		}
		return nil, err
	}

	return config, nil
}

func kerberosConfigFromPath(path string) (*config.Config, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return config.NewFromReader(file)
}

func kerberosCCacheFromName(ccname string) (*credentials.CCache, error) {
	var ccpath string

	if strings.HasPrefix(ccname, "FILE:") {
		ccpath = ccname[len("FILE:"):]
	} else if strings.HasPrefix(ccname, "DIR:") {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}

		dir := ccname[len("DIR:"):]

		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}

		ccpath = filepath.Join(dir, "krb5cc_"+u.Uid)
	} else {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}

		ccpath = "/tmp/krb5cc_" + u.Uid
	}

	return credentials.LoadCCache(ccpath)
}

func defaultKerberosCCache() (*credentials.CCache, error) {
	return kerberosCCacheFromName(os.Getenv("KRB5CCNAME"))
}
