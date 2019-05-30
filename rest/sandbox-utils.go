package rest

const (
	SANDBOX_SUFFIX = "-dev" // pay-app名称有该后缀，表示使用sandbox方式
)

func _IsSandbox(payApp string) (isSandbox bool) {
	sl := len(SANDBOX_SUFFIX)
	al := len(payApp)
	if al <= sl {
		return false
	}
	if payApp[al-sl:] == SANDBOX_SUFFIX {
		return true
	}
	return false
}
