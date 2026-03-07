package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrTechNotFound     = errors.New("tech category does not exist")
	ErrSnippetNotFound  = errors.New("snippet is not found")
	ErrRegistryNotFound = errors.New("no registry dir found")

	ErrManifestNotFound      = errors.New("manifest is not found")
	ErrRemoteTechNotFound    = errors.New("remote tech category not found")
	ErrRemoteSnippetNotFound = errors.New("remote snippet not found")

	ErrCachedManifestNotFound   = errors.New("cached manifest file is not found")
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

// Match check if a given tech category name match the current tech category
func (rt RemoteTech) Match(name string) bool {
	if strings.EqualFold(rt.Name, name) {
		return true
	}

	for _, alias := range rt.Aliases {
		if strings.EqualFold(alias.ID, name) || strings.EqualFold(alias.Name, name) {
			return true
		}
	}

	return false
}

type RemoteSnippet struct {
	ID       string `json:"id"`
	Topic    string `json:"topic"`
	Filename string `json:"filename"`
}

func (rs RemoteSnippet) Match(topic string) bool {
	if strings.EqualFold(rs.ID, topic) || strings.EqualFold(rs.Topic, topic) {
		return true
	}

	return strings.Contains(strings.ToLower(rs.Filename), topic)
}

type Manifest struct {
	Version string       `json:"version"`
	Techs   []RemoteTech `json:"techs"`
}

type CachedManifest struct {
	Manifest
	UpdatedAt time.Time
}

// IsTrustWorthy check if the cached manifest has'nt expired and not empty
func (cm CachedManifest) IsTrustWorthy() bool {
	expirationTime := cm.UpdatedAt.Add(time.Minute*30)
	if cm.UpdatedAt.IsZero() || expirationTime.Before(time.Now()) {
		return false
	}

	if len(cm.Techs) == 0 {
		return false
	}

	return true
}
