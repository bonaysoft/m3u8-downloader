package m3u8util

import (
	"bytes"
	"fmt"
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
		return fmt.Errorf("failed to fetch m3u8 URL %s: %w", m3u8Addr, err)
	}

	playList, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return fmt.Errorf("failed to decode m3u8 URL response: %w", err)
	}

	switch listType {
	case m3u8.MEDIA:
		return mediaPlDoFunc(playList.(*m3u8.MediaPlaylist))
	case m3u8.MASTER:
		return masterPlDoFunc(playList.(*m3u8.MasterPlaylist))
	}

	return fmt.Errorf("invalid listType for the m3u8")
}

func FetchKey(key *m3u8.Key, keyDecrypt func(rawKey []byte) ([]byte, error)) ([]byte, error) {
	resp, err := http.Get(key.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch key using uri %s: %w", key.URI, err)
	}
	defer resp.Body.Close()

	rawKey, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("faild to read key body: %w", err)
	}

	if keyDecrypt == nil {
		return rawKey, nil
	}

	plainKey, err := keyDecrypt(rawKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}

	return plainKey, nil
}

func Walk2Merge(chunkTsFileDir string, filename string) error {
	finalTsFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create ts file %s: %w", filename, err)
	}
	defer finalTsFile.Close()

	return filepath.WalkDir(chunkTsFileDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		tsBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read chunk ts file %s: %w", path, err)
		}

		// https://en.wikipedia.org/wiki/MPEG_transport_stream
		if syncByteIdx := bytes.IndexRune(tsBytes, 0x47); syncByteIdx > 0 {
			tsBytes = tsBytes[syncByteIdx:]
		}

		if _, err = finalTsFile.Write(tsBytes); err != nil {
			return fmt.Errorf("failed to write ts file %s: %w", filename, err)
		}

		return nil
	})
}
