package main

import "path/filepath"

var (
	dbname            = "database.sqlite"
	remoteRegistryUrl = "https://raw.githubusercontent.com/aboubakary833/cx-registry/main"
)

// Config is Codexa environment config struct.
// It resolved the paths to the root dir, the database and the registries.
type Config struct {
	rootPath string
}

func NewConfig(rootPath string) Config {
	return Config{
		rootPath: rootPath,
	}
}

// RootPath return the .codexa dir path
func (cfg Config) RootPath() string {
	return cfg.rootPath
}

// DatabasePath return the sqlite database path
func (cfg Config) DatabasePath() string {
	return filepath.Join(cfg.rootPath, dbname)
}

// LocalRegistryPath return the local registry dir path
func (cfg Config) LocalRegistryPath() string {
	return filepath.Join(cfg.rootPath, "registry")
}

// RemoteRegistryPath return the remote registry url
func (cfg Config) RemoteRegistryPath() string {
	return remoteRegistryUrl
}
