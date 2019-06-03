package rest

import (
	"strings"
)

const (
	SANDBOX_SUFFIX = "-dev" // pay-app名称有该后缀，表示使用sandbox方式
)

func _IsSandbox(payApp string) (isSandbox bool) {
	return strings.HasSuffix(payApp, SANDBOX_SUFFIX)
}
