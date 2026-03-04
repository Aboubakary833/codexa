package storage

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aboubakary833/codexa/internal/domain"
)

type registry struct {
	Path string
}

func NewRegistry(path string) *registry {
	return &registry{
		Path: path,
	}
}

// StoreContent store an entry
func (r *registry) StoreContent(ctx context.Context, path, content string) error {

	if !strings.HasPrefix(path, r.Path) {
		path = filepath.Join(r.Path, path)
	}

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	_, err = file.Write([]byte(content))

	return err
}

// LoadContent load a given snippet file content from the local registry
func (r *registry) LoadContent(ctx context.Context, snippet domain.Snippet) (string, error) {
	filePath := path.Join(r.Path, snippet.Filepath)
	rawBytes, err := os.ReadFile(filePath)

	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			return "", domain.ErrSnippetContentNotFound

		case errors.Is(err, os.ErrPermission):
			return "", domain.ErrSnippetContentCantBeRead

		default:
			return "", err
		}
	}

	return string(rawBytes), nil
}

// Stat check if the registry exists and is a directory.
// Return nil on success or error if the check fail
func (r *registry) Stat() error {
	info, err := os.Stat(r.Path)

	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			return domain.ErrRegistryNotFound
		}

		return err
	}

	if !info.IsDir() {
		return domain.ErrRegistryNotFound
	}

	return nil
}
