package domain

import (
	"errors"
	"time"
)

var (
	ErrTechNotFound             = errors.New("tech category does not exist")
	ErrSnippetNotFound          = errors.New("snippet is not found")
	ErrRegistryNotFound         = errors.New("no registry dir found")
	ErrSnippetContentNotFound   = errors.New("snippet target file don't exists")
	ErrSnippetContentCantBeRead = errors.New("snippet target file can't be read")
)

type Tech struct {
	//unique, canonical and stable identifier(ex: go, javascript, typescript)
	ID      string
	Name    string
	Aliases []TechAlias
}

type TechAlias struct {
	ID     string
	TechID string
	Name   string
}

type Snippet struct {
	// Identifier composed of tech category ID + topic in lowercase(ex: go:channels)
	ID     string
	TechID string
	// The actual entry topic(ex: channels, slice, map)
	Topic     string
	Filepath  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Manifest struct {
	Version string `json:"version"`
	Techs   []struct {
		Name    string   `json:"name"`
		Aliases []string `json:"aliases"`
		DirName string   `json:"dirname"`
	} `json:"tech"`
}
