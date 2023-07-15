package usecase

import (
	"fmt"

	"github.com/bonaysoft/m3u8-downloader/internal/entity"
)

type DecrypterConstructor func(properties entity.Properties) (Decrypter, error)

var supportDecrypter = map[string]DecrypterConstructor{
	"none":   NewDecrypterNone,
	"local":  NewDecrypterLocal,
	"remote": NewDecrypterRemote,
}

func NewDecrypter(name string, properties entity.Properties) (Decrypter, error) {
	d, ok := supportDecrypter[name]
	if !ok {
		return nil, fmt.Errorf("unknown decrypter: %v", name)
	}

	return d(properties)
}
