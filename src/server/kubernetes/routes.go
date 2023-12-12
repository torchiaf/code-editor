package kubernetes

import (
	"fmt"
	"os"
	"strings"
)

func GetUserHost(user string) string {
	return os.Getenv(fmt.Sprintf("CODE_EDITOR_%s_SERVICE_HOST", strings.ToUpper(user)))
}

func GetUserPort(user string) string {
	return os.Getenv(fmt.Sprintf("CODE_EDITOR_%s_SERVICE_PORT", strings.ToUpper(user)))
}
