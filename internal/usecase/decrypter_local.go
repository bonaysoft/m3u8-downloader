package usecase

import (
	"encoding/base64"
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/bonaysoft/m3u8-downloader/internal/entity"
)

var _ Decrypter = (*DecrypterLocal)(nil)

type DecrypterLocal struct {
	Cmd []string
}

func NewDecrypterLocal(properties entity.Properties) (Decrypter, error) {
	return &DecrypterLocal{
		Cmd: properties.Cmd,
	}, nil
}

func (d *DecrypterLocal) M3u8URLDecrypt(mu *entity.M3u8URL) error {
	cmd := exec.Command(d.Cmd[0], append(d.Cmd[1:], "m3u8Decrypt", mu.Encrypted)...)
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	return json.Unmarshal(out, &mu)
}

func (d *DecrypterLocal) KeyURLWrapper(uri string) string {
	cmd := exec.Command(d.Cmd[0], append(d.Cmd[1:], "keyURLWrapper", uri)...)
	out, err := cmd.Output()
	if err != nil {
		return uri
	}

	return strings.TrimSpace(string(out))
}

func (d *DecrypterLocal) KeyDecrypt(plainKey []byte) ([]byte, error) {
	cmd := exec.Command(d.Cmd[0], append(d.Cmd[1:], "keyDecrypt", base64.StdEncoding.EncodeToString(plainKey))...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(out))
}
