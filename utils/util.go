package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}