package storage

import (
	"errors"
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"time"
)

type Storage interface {
	Get(string, bool) (*os.File, error)
	Put(string, io.Reader) error
	Del(string) error
}

type Configuration struct {
	BasePath string
	CachePath string
}

func NewLocalStorage(configuration Configuration) Storage {
	sg := &storage{
		BasePath: configuration.BasePath,
		CachePath: path.Join(configuration.BasePath, "_cache"),
	}

	err := os.MkdirAll(sg.BasePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(sg.CachePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return sg
}

type storage struct {
	BasePath string
	CachePath string
}

func (s *storage) Get(key string, checkExpiration bool) (file *os.File, err error) {
	filePath := s.pathResolve(key)

	if !s.checkFileIsExist(filePath) {
		err = errors.New("not found")
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		return
	}

	fileStat, err := f.Stat()
	if err != nil {
		return
	}

	if s.isExpiration(fileStat.ModTime()) {
		_ = s.Del(key)

		err = errors.New("file expiration")
		return
	}

	return f, err
}

func (s *storage) Put(key string, data io.Reader) (err error) {
	tempUUID, err := uuid.NewV4()
	if err != nil {
		return
	}
	tempFileName := path.Join(s.CachePath, tempUUID.String())

	file, err := os.Create(tempFileName)
	if err != nil {
		return
	}
	defer file.Close()

	io.Copy(file, data)

	fileStat, err := file.Stat()
	if err != nil {
		return
	}
	if fileStat.Size() == 0 {
		_ = os.RemoveAll(tempFileName)
		return
	}

	filePath := s.pathResolve(key)

	err = os.MkdirAll(path.Join(filePath, ".."), os.ModePerm)
	if err != nil {
		return
	}

	err = os.Rename(tempFileName, filePath)
	if err != nil {
		_ = os.RemoveAll(tempFileName)
	}

	return
}

func (s *storage) Del(key string) (err error) {
	filePath := s.pathResolve(key)

	err = os.RemoveAll(filePath)

	return
}

func (s *storage) pathResolve(elem ...string) string {
	return path.Join(append([]string{s.BasePath}, elem...)...)
}

func (*storage) checkFileIsExist(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (*storage) isExpiration(t time.Time) bool {
	d, _:= time.ParseDuration(viper.GetString("storage.metadataDuration"))

	c := t.Add(d)

	return c.Before(time.Now())
}
