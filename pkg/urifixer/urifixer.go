package urifixer

import (
	"net/url"
	"path"
	"strings"
)

type FixerOpt struct {
	Host  string
	Path  string
	Query string
}

func (p *FixerOpt) apply(u *url.URL) {
	if p.Host != "" {
		u.Host = p.Host
	}

	if p.Path != "" {
		u.Path = path.Join(p.Path, path.Base(u.Path))
	}

	if p.Query != "" {
		query := u.Query()
		q, _ := url.ParseQuery(p.Query)
		for k, vs := range q {
			for _, v := range vs {
				query.Set(k, v)
			}
		}
		u.RawQuery = query.Encode()
	}
}

func MakeUp(uri string, pu *url.URL, opts ...FixerOpt) string {
	u, err := url.Parse(uri)
	if err == nil && !strings.HasPrefix(u.Path, "/") {
		u.Path = path.Join(path.Dir(pu.Path), u.Path)
	}

	if u.Scheme == "" {
		u.Scheme = pu.Scheme
	}
	if u.Host == "" {
		u.Host = pu.Host
	}

	query := u.Query()
	for k, vs := range pu.Query() {
		query[k] = append(query[k], vs...)
	}
	u.RawQuery = query.Encode()

	for _, opt := range opts {
		opt.apply(u)
	}

	return u.String()
}
