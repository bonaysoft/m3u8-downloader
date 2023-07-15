package urifixer

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixURI(t *testing.T) {
	pu, _ := url.Parse("https://example.com/xxxx/xxx.m3u8")
	tests := []struct {
		uri      string
		expected string
	}{
		{uri: "https://example.com/abc.ts", expected: "https://example.com/abc.ts"},
		{uri: "/abc.ts", expected: "https://example.com/abc.ts"},
		{uri: "/abc.ts?signature=xxxxxxxx&t=16000123123", expected: "https://example.com/abc.ts?signature=xxxxxxxx&t=16000123123"},
		{uri: "abc.ts", expected: "https://example.com/xxxx/abc.ts"},
		{uri: "abc.ts?start=0&offset=100", expected: "https://example.com/xxxx/abc.ts?offset=100&start=0"},
	}

	for idx, test := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			u, err := url.Parse(MakeUp(test.uri, pu))
			assert.NoError(t, err)
			assert.Equal(t, test.expected, u.String())
		})
	}
}
