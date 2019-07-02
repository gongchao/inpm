package services

import (
	"github.com/spf13/viper"
	"io"
	"npm/utils"
)

type PackageService interface {
	GetDataFromRemote(url string) (data io.ReadCloser, err error)
}

func NewPackageService() PackageService {
	return &packageService{}
}

type packageService struct {
}

func (*packageService) GetDataFromRemote(url string) (data io.ReadCloser, err error) {
	data, err = utils.Request(viper.GetString("uplinks.npm.url") + url)
	return
}