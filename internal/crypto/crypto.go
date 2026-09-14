package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func CreateMD5Hash(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
