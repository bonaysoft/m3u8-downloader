package usecase

import (
	"encoding/base64"
	"encoding/json"

	"github.com/bonaysoft/m3u8-downloader/internal/entity"
	"github.com/go-resty/resty/v2"
)

const (
	pathM3u8RULDecrypt = "/m3u8-url-decrypt"
	pathKeyURLWrapper  = "/key-url-wrapper"
	pathKeyDecrypt     = "/key-decrypt"
)

type DecrypterRemote struct {
	hc *resty.Client
}

func NewDecrypterRemote(props entity.Properties) (Decrypter, error) {
	return &DecrypterRemote{
		hc: resty.New().SetBaseURL(props.BaseURL).SetAuthToken(props.Token),
	}, nil
}

func (d *DecrypterRemote) M3u8URLDecrypt(mu *entity.M3u8URL) error {
	if mu.Encrypted == "" {
		return nil
	}

	resp, err := d.hc.R().SetFormData(map[string]string{"encrypted": mu.Encrypted}).Post(pathM3u8RULDecrypt)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp.Body(), mu)
}

func (d *DecrypterRemote) KeyURLWrapper(uri string) string {
	resp, err := d.hc.R().SetFormData(map[string]string{"uri": uri}).Post(pathKeyURLWrapper)
	if err != nil {
		return ""
	}

	return resp.String()
}

func (d *DecrypterRemote) KeyDecrypt(plainKey []byte) ([]byte, error) {
	resp, err := d.hc.R().SetFormData(map[string]string{"plainKey": base64.URLEncoding.EncodeToString(plainKey)}).Post(pathKeyDecrypt)
	return resp.Body(), err
}
