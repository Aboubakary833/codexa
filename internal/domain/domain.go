package domain

import (
	"errors"
	"time"
)

var (
	ErrTechNotFound     = errors.New("tech category does not exist")
	ErrSnippetNotFound  = errors.New("snippet is not found")
	ErrRegistryNotFound = errors.New("no registry dir found")

	ErrManifestNotFound      = errors.New("manifest is not found")
	ErrRemoteTechNotFound    = errors.New("remote tech category not found")
	ErrRemoteSnippetNotFound = errors.New("remote snippet not found")

	ErrSnippetContentNotFound   = errors.New("snippet target file don't exists")
	ErrSnippetContentCantBeRead = errors.New("snippet target file can't be read")
)

type Tech struct {
	//unique, canonical and stable identifier(ex: go, javascript, typescript)
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Aliases []TechAlias `json:"aliases"`
}

type TechAlias struct {
	ID     string `json:"id"`
	TechID string `json:"tech_id,omitempty"`
	Name   string `json:"name"`
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

type RemoteTech struct {
	Tech
	Dirname string `json:"dirname"`
}

type RemoteSnippet struct {
	ID       string `json:"id"`
	Topic    string `json:"topic"`
	Filename string `json:"filename"`
}

type Manifest struct {
	Version string       `json:"version"`
	Techs   []RemoteTech `json:"techs"`
}
