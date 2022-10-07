package links

import "net/url"

const jsScheme = "javascript"

func cleanURL(base *url.URL, link string) (rv string, ok bool) {
	u, err := url.Parse(link)
	if err != nil {
		return
	}

	if u.Host == "" {
		if u = base.ResolveReference(u); u.Host == "" {
			return
		}
	}

	switch u.Scheme {
	case jsScheme:
		return
	case "":
		u.Scheme = base.Scheme
	}

	if u.Path == "" {
		u.Path = "/"
	}

	u.Fragment = ""

	return u.String(), true
}
