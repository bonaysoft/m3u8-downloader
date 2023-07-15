package usecase

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bonaysoft/m3u8-downloader/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewDecrypterRemote(t *testing.T) {
	http.HandleFunc(pathM3u8RULDecrypt, func(w http.ResponseWriter, r *http.Request) {
		encryptedURL := r.PostFormValue("encrypted")
		if encryptedURL == "" {
			return
		}

		rr, err := base64.URLEncoding.DecodeString(encryptedURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := json.Marshal(map[string]string{"plain_url": string(rr)})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(result)
	})
	http.HandleFunc(pathKeyURLWrapper, func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.PostFormValue("uri"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		q := u.Query()
		q.Add("uid", "test123")
		u.RawQuery = q.Encode()

		_, _ = w.Write([]byte(u.String()))

	})
	http.HandleFunc(pathKeyDecrypt, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte{115, 67, 132, 22, 228, 243, 158, 78, 219, 160, 176, 214, 103, 96, 143, 168})
	})

	s := httptest.NewServer(http.DefaultServeMux)

	d, err := NewDecrypterRemote(entity.Properties{BaseURL: s.URL, Token: "test123"})
	assert.NoError(t, err)

	mu := entity.NewM3u8URL(base64.URLEncoding.EncodeToString([]byte("https://example.com/test.m3u8")))
	assert.NoError(t, d.M3u8URLDecrypt(mu))
	assert.Equal(t, "https://example.com/test.m3u8", mu.PlainURL)

	assert.Equal(t, "https://example.com/test-key-uri?uid=test123", d.KeyURLWrapper("https://example.com/test-key-uri"))

	decrypted, err := d.KeyDecrypt([]byte{6, 28, 178, 34, 211, 151, 172, 45, 236, 198, 131, 230, 5, 85, 186, 247})
	assert.NoError(t, err)
	assert.Equal(t, []byte{115, 67, 132, 22, 228, 243, 158, 78, 219, 160, 176, 214, 103, 96, 143, 168}, decrypted)

}
