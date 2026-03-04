package main

import (
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/aboubakary833/codexa/internal/adapters/cli"
	"github.com/aboubakary833/codexa/internal/adapters/sqlite"
	"github.com/aboubakary833/codexa/internal/adapters/storage"
	"github.com/aboubakary833/codexa/internal/application"
	_ "github.com/mattn/go-sqlite3"
)

var rootDirName = ".codexa"

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

	snippetRepository := sqlite.NewSnippetRepository(db)
	techRepository := sqlite.NewTechRepository(db)
	registry := storage.NewRegistry(config.LocalRegistryPath())

	app := application.New(snippetRepository, techRepository, registry)
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	cli.NewCommandWrapper(app, logger).Execute()
}
