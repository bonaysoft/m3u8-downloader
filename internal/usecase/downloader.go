package usecase

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bonaysoft/m3u8-downloader/internal/entity"
	"github.com/bonaysoft/m3u8-downloader/pkg/aesutil"
	"github.com/bonaysoft/m3u8-downloader/pkg/urifixer"
	"github.com/grafov/m3u8"
)

type Downloader struct {
	de Decrypter
}

func NewDownloader(de Decrypter) *Downloader {
	return &Downloader{
		de: de,
	}
}

func (dl *Downloader) Fetch(m3u8Addr string) (m3u8.Playlist, m3u8.ListType, error) {
	resp, err := http.Get(m3u8Addr)
	if err != nil {
		return nil, 0, err
	}

	return m3u8.DecodeFrom(resp.Body, true)
}

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

	playList, listType, err := dl.Fetch(u.String())
	if err != nil {
		return fmt.Errorf("failed to fetch playlist: %v", err)
	}

	switch listType {
	case m3u8.MEDIA:
		mediapl := playList.(*m3u8.MediaPlaylist)
		for idx, seg := range mediapl.GetAllSegments() {
			seg.URI = urifixer.MakeUp(seg.URI, u, urifixer.FixerOpt(mu.TsURLPart))
			fmt.Printf("downloading %d: %s\n", idx, seg.URI)
			if err := dl.download(seg); err != nil {
				return fmt.Errorf("failed to download: %s", err)
			}
		}

	case m3u8.MASTER:
		masterpl := playList.(*m3u8.MasterPlaylist)
		for _, variant := range masterpl.Variants {
			if strings.HasSuffix(variant.URI, ".m3u8") {
				fmt.Printf("Downloading the MasterPlaylist from %s\n", variant.URI)
				if err := dl.Download(entity.NewM3u8URL(variant.URI)); err != nil {
					return err
				}
			}

			if variant.Chunklist == nil {
				break
			}
			for idx, seg := range variant.Chunklist.Segments {
				seg.URI = urifixer.MakeUp(seg.URI, u, urifixer.FixerOpt(mu.TsURLPart))
				fmt.Printf("downloading %d: %s\n", idx, seg.URI)

				if err := dl.download(seg); err != nil {
					return err
				}
			}
		}
	}

	return dl.merge()
}

func (dl *Downloader) download(seg *m3u8.MediaSegment) error {
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

	return dl.save(chunk)
}

func (dl *Downloader) decrypt(chunk []byte, key *m3u8.Key) ([]byte, error) {
	keyURL := dl.de.KeyURLWrapper(key.URI)
	resp, err := http.Get(keyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawKey, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	plainKey, err := dl.de.KeyDecrypt(rawKey)
	if err != nil {
		return nil, err
	}

	return aesutil.AES128Decrypt(chunk, plainKey, []byte(key.IV))
}

func (dl *Downloader) save(chunk []byte) error {
	return nil
}

func (dl *Downloader) merge() error {
	return nil
}
