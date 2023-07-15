package usecase

import (
	"github.com/bonaysoft/m3u8-downloader/internal/entity"
)

type Decrypter interface {
	M3u8URLDecrypt(mu *entity.M3u8URL) error
	KeyURLWrapper(uri string) string
	KeyDecrypt(plainKey []byte) ([]byte, error)
}
