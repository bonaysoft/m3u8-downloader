package httputil

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

var defaultClient = retryablehttp.NewClient()

func init() {
	defaultClient.Logger = nil
}

func Get(url string) (*http.Response, error) {
	return defaultClient.Get(url)
}
