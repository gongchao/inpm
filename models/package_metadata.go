package models

import "time"

type PackageMetadata struct {
	Name     string `json:"name"`
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	Time struct {
		Modified time.Time `json:"modified"`
	} `json:"time"`
	Versions map[string]ModuleMetadataVersion `json:"versions"`
}

type ModuleMetadataVersion struct {
	Name            string                 `json:"name"`
	Version         string                 `json:"version"`
	Description     string                 `json:"description"`
	Dependencies    map[string]interface{} `json:"dependencies"`
	DevDependencies map[string]interface{} `json:"devDependencies"`
	Dist            struct {
		Shasum  string `json:"shasum"`
		Tarball string `json:"tarball"`
	} `json:"dist"`
}
