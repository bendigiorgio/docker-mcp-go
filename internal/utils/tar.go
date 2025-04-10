package utils

import (
	"archive/tar"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
)

func Base64ToTarReader(base64Data string) (io.ReadCloser, error) {
	// Decode the base64 string
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	reader := bytes.NewReader(decoded)

	tarReader := tar.NewReader(reader)
	_, err = tarReader.Next()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("invalid tar archive: %w", err)
	}

	reader.Seek(0, io.SeekStart)

	return &nopCloser{reader}, nil
}

// nopCloser wraps an io.Reader to implement io.ReadCloser with a no-op Close method
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
