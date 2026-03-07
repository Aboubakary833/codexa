package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/aboubakary833/codexa/internal/domain"
)

type registry struct {
	// the registry dir path
	Path string
}

func NewRegistry(path string) *registry {
	return &registry{
		Path: path,
	}
}

// CreateOrUpdateSnippet store a snippet content or update it if the snippet file already exist.
func (r *registry) CreateOrUpdateSnippet(ctx context.Context, path, content string) error {
	file, err := r.openFile(path)
	
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = file.Write([]byte(content))

	return err
}

// LoadSnippet load a given snippet file content from the local registry.
func (r *registry) LoadSnippet(ctx context.Context, snippet domain.Snippet) (string, error) {
	rawBytes, err := r.readFile(snippet.Filepath)

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

// GetManifest return a cached version of the Manifest.
// If no manifest cache file exists, a new one will automatically be created.
func (r *registry) GetManifest(ctx context.Context) (domain.CachedManifest, error) {
	rawBytes, err := r.readFile("index.json")

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.CachedManifest{}, domain.ErrCachedManifestNotFound
		}

		return domain.CachedManifest{}, err
	}
	var manifest domain.CachedManifest

	if err = json.Unmarshal(rawBytes, &manifest); err != nil {
		return domain.CachedManifest{}, domain.ErrCachedManifestNotFound
	}

	return manifest, nil
}


func (r *registry) CreateOrUpdateManifest(ctx context.Context, manifest domain.CachedManifest) error {
	file, err := r.openFile("index.json")

	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	return encoder.Encode(manifest)
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

// readFile is a helper function for reading and returning a file content
func (r *registry) readFile(path string) ([]byte, error) {
	filePath := r.getFileFullPath(path)
	return os.ReadFile(filePath)
}

// openFile is the registry helper function for opening file with "READ" and "WRITE" permissions.
// If the given file does'nt exist, it will be automatically be created.
func (r *registry) openFile(path string) (*os.File, error) {
	filePath := r.getFileFullPath(path)
	return os.OpenFile(filePath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, os.ModePerm)
}

func (r *registry) getFileFullPath(path string) string {
	if !strings.HasPrefix(path, r.Path) {
		path = filepath.Join(r.Path, path)
	}

	return path
}
