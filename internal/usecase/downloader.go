package usecase

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/bonaysoft/m3u8-downloader/internal/entity"
	"github.com/bonaysoft/m3u8-downloader/pkg/aesutil"
	"github.com/bonaysoft/m3u8-downloader/pkg/m3u8util"
	"github.com/bonaysoft/m3u8-downloader/pkg/urifixer"
	"github.com/grafov/m3u8"
	"github.com/saltbo/gopkg/strutil"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	de Decrypter
}

func NewDownloader(de Decrypter) *Downloader {
	return &Downloader{
		de: de,
	}
}

// Download fetch the m3u8 and download
func (dl *Downloader) Download(mu *entity.M3u8URL) error {
	if mu.Encrypted != "" {
		if err := dl.de.M3u8URLDecrypt(mu); err != nil {
			return err
		}
	}

	u, err := url.Parse(mu.PlainURL)
	if err != nil {
		return fmt.Errorf("failed to parse playlist: %v", err)
	}

	mediaPLFunc := func(playlist *m3u8.MediaPlaylist) error {
		tsFileDir, err := os.MkdirTemp("", "bonaysoft-m3u8-downloader-"+strutil.Md5HexShort(u.String()))
		if err != nil {
			return fmt.Errorf("failed to create temporary download dir: %v", err)
		}
		defer os.RemoveAll(tsFileDir)

		bar := progressbar.Default(int64(len(playlist.GetAllSegments())), "Downloading...")
		for _, seg := range playlist.GetAllSegments() {
			_ = bar.Add(1)

			seg.URI = urifixer.MakeUp(seg.URI, u, urifixer.FixerOpt(mu.TsURLPart))
			if err := dl.download(seg, tsFileDir); err != nil {
				return fmt.Errorf("failed to download: %s", err)
			}
		}
		_ = bar.Finish()
		fmt.Println("Download done!")

		return m3u8util.Walk2Merge(tsFileDir, strings.ReplaceAll(path.Base(u.Path), ".m3u8", ".ts"))
	}

	masterPLFunc := func(playlist *m3u8.MasterPlaylist) error {
		// TODO	select a variant
		variant := playlist.Variants[0]
		return m3u8util.FetchM3u8Do(variant.URI, mediaPLFunc, nil)
	}

	return m3u8util.FetchM3u8Do(u.String(), mediaPLFunc, masterPLFunc)
}

// download DL the target seg into the specify directory
func (dl *Downloader) download(seg *m3u8.MediaSegment, chunkTsFileDir string) error {
	resp, err := http.Get(seg.URI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	chunk, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("%v", string(chunk))
	}

	if seg.Key != nil {
		decrypted, err := dl.decrypt(chunk, seg.Key)
		if err != nil {
			return err
		}

		chunk = decrypted
	}

	f, err := os.Create(path.Join(chunkTsFileDir, path.Base(seg.URI)))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(chunk)
	return err
}

func (dl *Downloader) decrypt(chunk []byte, key *m3u8.Key) ([]byte, error) {
	key.URI = dl.de.KeyURLWrapper(key.URI)

	keyBytes, err := m3u8util.FetchKey(key, dl.de.KeyDecrypt)
	if err != nil {
		return nil, err
	}

	return aesutil.AES128Decrypt(chunk, keyBytes, []byte(key.IV))
}
