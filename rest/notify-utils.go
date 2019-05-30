package rest

import (
	"fmt"
)

func _AppendAppName(uri string, appName string) string {
	l := len(uri)
	if uri[l-1] == '/' {
		return fmt.Sprintf("%s%s", uri, appName)
	}
	return fmt.Sprintf("%s/%s", uri, appName)
}
