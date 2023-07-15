package usecase

import (
	"fmt"

	"github.com/bonaysoft/m3u8-downloader/internal/entity"
)

var _ Decrypter = (*DecrypterNone)(nil)

type DecrypterNone struct {
}

func NewDecrypterNone(properties entity.Properties) (Decrypter, error) {
	return &DecrypterNone{}, nil
}

func (d *DecrypterNone) M3u8URLDecrypt(mu *entity.M3u8URL) error {
	if mu.Encrypted != "" {
		return fmt.Errorf("encrypted m3u8 url, you must specify a decrypter to decrypt it")
	}

	return nil
}

func (d *DecrypterNone) KeyURLWrapper(uri string) string {
	return uri
}

func (d *DecrypterNone) KeyDecrypt(plainKey []byte) ([]byte, error) {
	return plainKey, nil
}
