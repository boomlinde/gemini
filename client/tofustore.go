package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type tofustore struct {
	dir string
}

func (s *tofustore) Pin(host string, signature string) error {
	err := os.MkdirAll(s.dir, 0700)
	if err != nil {
		return fmt.Errorf("failed to create tofu dir: %w", err)
	}

	tfile, err := ioutil.TempFile(s.dir, "sig_*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tname := tfile.Name()

	if _, err := tfile.Write([]byte(signature + "\n")); err != nil {
		tfile.Close()
		os.Remove(tname)
		return fmt.Errorf("failed to write tofu: %w", err)
	}

	if err := tfile.Sync(); err != nil {
		tfile.Close()
		os.Remove(tname)
		return fmt.Errorf("failed to sync: %w", err)
	}

	tfile.Close()

	wpath := filepath.Join(s.dir, host)
	if err := os.Rename(tname, wpath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func (s *tofustore) Get(host string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(s.dir, host))
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}
