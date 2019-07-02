package models

import (
	"path"
	"strings"
)

type PackagePath struct {
	ScopeName   string
	PackageName string
	TarballName string
}

func (p *PackagePath) ToKey() string {
	return strings.Join([]string{p.ScopeName, p.PackageName}, "/")
}

func (p *PackagePath) GetMKey() string {
	return path.Join(p.ToKey(), "package.json")
}

func (p *PackagePath) GetMRUrl() string {
	return p.ToKey()
}

func (p *PackagePath) GetTKey() string {
	return path.Join(p.ToKey(), p.TarballName)
}

func (p *PackagePath) GetTRUrl() string {
	return path.Join(p.ToKey(), "/-/", p.TarballName)
}
