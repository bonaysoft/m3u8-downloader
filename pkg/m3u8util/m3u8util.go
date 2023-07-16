package m3u8util

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/grafov/m3u8"
)

func FetchM3u8Do(m3u8Addr string, mediaPlDoFunc func(playlist *m3u8.MediaPlaylist) error, masterPlDoFunc func(playlist *m3u8.MasterPlaylist) error) error {
	resp, err := http.Get(m3u8Addr)
	if err != nil {
		return err
	}

	playList, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return err
	}

	switch listType {
	case m3u8.MEDIA:
		return mediaPlDoFunc(playList.(*m3u8.MediaPlaylist))
	case m3u8.MASTER:
		return masterPlDoFunc(playList.(*m3u8.MasterPlaylist))
	}

	return nil
}

func FetchKey(key *m3u8.Key, keyDecrypt func(rawKey []byte) ([]byte, error)) ([]byte, error) {
	resp, err := http.Get(key.URI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawKey, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if keyDecrypt == nil {
		return rawKey, nil
	}

	plainKey, err := keyDecrypt(rawKey)
	if err != nil {
		return nil, err
	}

	return plainKey, err
}

func Walk2Merge(chunkTsFileDir string, filename string) error {
	finalTsFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer finalTsFile.Close()

	return filepath.WalkDir(chunkTsFileDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		tsBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// https://en.wikipedia.org/wiki/MPEG_transport_stream
		if syncByteIdx := bytes.IndexRune(tsBytes, 0x47); syncByteIdx > 0 {
			tsBytes = tsBytes[syncByteIdx:]
		}

		_, err = finalTsFile.Write(tsBytes)
		return err
	})
}
