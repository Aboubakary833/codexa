package main

import (
	"log"
	"os"
	"path"

	"github.com/aboubakary833/codexa/internal/adapters/cli"
	"github.com/aboubakary833/codexa/internal/adapters/remote"
	"github.com/aboubakary833/codexa/internal/adapters/sqlite"
	"github.com/aboubakary833/codexa/internal/adapters/storage"
	"github.com/aboubakary833/codexa/internal/application"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Version		= "dev"
	rootDirName = ".codexa"
)

func main() {
	homePath, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln(err)
	}
	rootDirPath := path.Join(homePath, rootDirName)
	config := NewConfig(rootDirPath)

	bootstrapper := Bootstrapper{cfg: config}

	if err = bootstrapper.Setup(); err != nil {
		log.Fatalln(err)
	}

	db, err := bootstrapper.InitDatabase()

	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	logger := bootstrapper.InitLogger()

	registry := storage.NewRegistry(config.LocalRegistryPath())
	fetcher := remote.NewFetcher(config.RemoteRegistryPath())
	snippetRepository := sqlite.NewSnippetRepository(db)
	techRepository := sqlite.NewTechRepository(db)

	app := application.New(
		snippetRepository, techRepository,
		registry, fetcher,
	)

	cli.NewCommandWrapper(app, Version, logger).Execute()
}
