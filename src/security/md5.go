package security

import (
	"crypto/md5"
	"io"
)

// Handles MD5 encryption.
func EncryptMD5(data []byte) []byte {
	h := md5.New()
	io.Writer.Write(h, data)

	return h.Sum(nil)
}
