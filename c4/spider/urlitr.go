package spider

import (
	"github.com/mg/i"
	"github.com/mg/i/hoi"
	"strings"
)

// Normalize Iterator
// Removes empty urls and urls starting with tel: and mailto:
// Removes hash strings, trailing slashes
// Returns lower case string
func removeUrl(itr i.Iterator) bool {
	url, _ := itr.Value().(string)
	url = strings.TrimSpace(url)
	if url == "" {
		return false
	}
	if strings.HasPrefix(url, "tel:") || strings.HasPrefix(url, "mailto:") {
		return false
	}
	return true
}

func remHash(u string) string {
	idx := strings.Index(u, "#")
	if idx == -1 {
		return u
	}
	if idx == 0 {
		return ""
	}
	return u[0 : idx-1]
}

func normalizeUrl(itr i.Iterator) interface{} {
	url, _ := itr.Value().(string)
	url = strings.TrimSuffix(remHash(url), "/")
	return strings.ToLower(url)
}

func NormalizeItr(itr i.Forward) i.Forward {
	return hoi.Filter(removeUrl, hoi.Map(normalizeUrl, itr))
}

// Host Mapper
// Maps host name to urls that don't have it
// eg /a/b -> http://www.domain.com/a/b
func hostFromBase(base string) string {
	slashslashidx := strings.Index(base, "//")
	idx := strings.Index(base[slashslashidx+2:], "/")
	if idx > 0 {
		return base[0 : slashslashidx+2+idx]
	}
	return base
}

func hostMapper(base string) hoi.MapFunc {
	host := hostFromBase(base)
	if !strings.HasSuffix(base, "/") {
		base = base + "/"
	}
	return func(itr i.Iterator) interface{} {
		url, _ := itr.Value().(string)
		if strings.Contains(url, "://") {
			return url
		}
		if strings.HasPrefix(url, "/") {
			return host + url
		}
		return base + url
	}
}

func HostMapper(base string, itr i.Forward) i.Forward {
	return hoi.Map(hostMapper(base), itr)
}

// Referer
// Returns url with referer
func referer(ref string) hoi.MapFunc {
	return func(itr i.Iterator) interface{} {
		url, _ := itr.Value().(string)
		return []string{url, ref}
	}
}

func Referer(ref string, itr i.Forward) i.Forward {
	return hoi.Map(referer(ref), itr)
}

// Remove url if it is not refered by ref
func removeIfNotReferedBy(ref string) hoi.FilterFunc {
	host := hostFromBase(ref)
	return func(itr i.Iterator) bool {
		ref, _ := itr.Value().([]string)
		return strings.HasPrefix(ref[1], host)
	}
}

func BindByRef(ref string, itr i.Forward) i.Forward {
	return hoi.Filter(removeIfNotReferedBy(ref), itr)
}

// Remove url if it is not on host
func removeIfNotOnHost(url string) hoi.FilterFunc {
	host := hostFromBase(url)
	return func(itr i.Iterator) bool {
		u, _ := itr.Value().([]string)
		return strings.HasPrefix(u[0], host)
	}
}

func BindByHost(host string, itr i.Forward) i.Forward {
	return hoi.Filter(removeIfNotOnHost(host), itr)
}
