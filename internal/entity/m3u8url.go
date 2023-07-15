package entity

import (
	"net/url"
	"strings"
)

type TsURLPart struct {
	Host  string `json:"host" form:"host"`
	Path  string `json:"path" form:"path"`
	Query string `json:"query" form:"query"`
}

type M3u8URL struct {
	Encrypted string    `json:"encrypted" form:"encrypted"`
	PlainURL  string    `json:"plain_url" form:"plain_url"`
	TsURLPart TsURLPart `json:"ts_url_part" form:"ts_url_part"`
}

func NewM3u8URL(v string) *M3u8URL {
	if u, err := url.Parse(v); err == nil && strings.HasPrefix(v, "http") {
		return &M3u8URL{
			PlainURL: u.String(),
		}
	}

	return &M3u8URL{
		Encrypted: v,
	}
}
