package robots

type tokenKind byte

const (
	kindNone      tokenKind = 0
	kindUserAgent tokenKind = 1
	kindAllow     tokenKind = 2
	kindDisallow  tokenKind = 3
	kindSitemap   tokenKind = 4
)
