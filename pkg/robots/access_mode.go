package robots

type accessMode byte

const (
	allowAll accessMode = 0
	gotRules accessMode = 1
	denyAll  accessMode = 2
)
