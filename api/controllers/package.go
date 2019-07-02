package controllers

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"npm/lib/storage"
	"npm/models"
	"npm/services"
	"strconv"
)

type PackageController struct {
	Service services.PackageService
	Sg storage.Storage
}

func (p *PackageController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/{package}", "GetMetadata")
	b.Handle("GET", "/{scope}/{package}", "RedirectGetMetadata")

	b.Handle("GET", "/{package}/-/{tarball}", "GetTarball")
	b.Handle("GET", "/{scope}/{package}/-/{tarball}", "RedirectGetTarball")
}

func (p *PackageController) GetMetadata(ctx iris.Context) {
	packagePath := parsePackagePath(ctx)

	key := packagePath.GetMKey()

	metadataFile, err := p.Sg.Get(key, true)
	if err == nil {
		defer metadataFile.Close()

		_, _ = io.Copy(ctx.ResponseWriter(), metadataFile)
		return
	}

	data, err := p.Service.GetDataFromRemote(packagePath.GetMRUrl())
	if err != nil {
		ctx.Application().Logger().Errorf("没找到", key)
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}
	defer data.Close()

	var metadata models.PackageMetadata
	err = jsoniter.NewDecoder(data).Decode(&metadata)
	if err != nil {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	// Rewrite dist tarball url
	rewriteTarballUrl(&metadata)

	ctx.JSON(metadata)

	go func() {
		source, err := jsoniter.Marshal(metadata)
		if err != nil {
			return
		}

		_ = p.Sg.Put(key, bytes.NewBuffer(source))
	}()
}

func (p *PackageController) RedirectGetMetadata(ctx iris.Context) {
	p.GetMetadata(ctx)
}

func (p *PackageController) GetTarball(ctx iris.Context) {
	packagePath := parsePackagePath(ctx)

	key := packagePath.GetTKey()

	data, err := p.Sg.Get(key, false)
	if err == nil {
		defer data.Close()

		dataStat, _ := data.Stat()
		if dataStat.Size() == 0 {
			ctx.Application().Logger().Errorf("异常: GetTarball-", key)
		}

		ctx.Header("content-length", strconv.FormatInt(dataStat.Size(), 10))

		_ = ctx.ServeContent(data, packagePath.TarballName, dataStat.ModTime(), false)

		return
	}

	reData, err := p.Service.GetDataFromRemote(packagePath.GetTRUrl())
	if err != nil {
		ctx.Application().Logger().Errorf("异常: GetTarball-GetDataFromRemote", key)
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}
	defer reData.Close()

	_ = p.Sg.Put(key, io.TeeReader(reData, ctx.ResponseWriter()))
}

func (p *PackageController) RedirectGetTarball(ctx iris.Context) {
	p.GetTarball(ctx)
}

func parsePackagePath(ctx iris.Context) *models.PackagePath {
	return &models.PackagePath{
		ScopeName: ctx.Params().GetString("scope"),
		PackageName: ctx.Params().GetString("package"),
		TarballName: ctx.Params().GetString("tarball"),
	}
}

func rewriteTarballUrl(metadata *models.PackageMetadata) {
	for key := range metadata.Versions {
		version := metadata.Versions[key]

		var _url *url.URL

		_url, err := url.Parse(version.Dist.Tarball)
		if err != nil {
			return
		}

		_url.Scheme = viper.GetString("app.protocol")
		_url.Host = viper.GetString("app.host")
		version.Dist.Tarball = _url.String()
		metadata.Versions[key] = version
	}
}