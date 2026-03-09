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

	ErrManifestNotFound   = errors.New("manifest is not found")

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
	ID        string    `json:"id"`
	TechID    string    `json:"tech_id,omitempty"`
	Topic     string    `json:"topic"`
	Filepath  string    `json:"filepath,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Match check if the given topic match the snippet
func (s Snippet) Match(topic string) bool {
	name := strings.Split(s.ID, ":")[1]
	
	if strings.EqualFold(name, topic) || strings.EqualFold(s.Topic, topic) {
		return true
	}
	
	return false
}

type RemoteTech struct {
	Tech
	Dirname string `json:"dirname"`
}

// Match check if the given name match the tech category
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
	expirationTime := cm.UpdatedAt.Add(time.Minute * 30)
	if cm.UpdatedAt.IsZero() || expirationTime.Before(time.Now()) {
		return false
	}

	if len(cm.Techs) == 0 {
		return false
	}

	return true
}
