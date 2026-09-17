package identificador

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func UUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func Segredo() string { b := make([]byte, 32); _, _ = rand.Read(b); return hex.EncodeToString(b) }
func Hash(valor string) string {
	soma := sha256.Sum256([]byte(valor))
	return hex.EncodeToString(soma[:])
}
