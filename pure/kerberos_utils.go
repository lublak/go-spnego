package pure

import (
	"os"
	"runtime"

	"github.com/jcmturner/gokrb5/v8/config"
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

func kerberosConfig() (*config.Config, error) {
	configPath := os.Getenv("KRB5_CONFIG")

	if len(configPath) == 0 {
		return kerberosConfigFromPath(getDefaultConfigPath())
	}

	config, err := kerberosConfigFromPath(configPath)

	if err != nil {
		if os.IsNotExist(err) {
			return kerberosConfigFromPath(getDefaultConfigPath())
		}
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
