package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

func Hash(v ...string) string {
	sort.Strings(v)

	h := sha256.New()
	h.Write([]byte(strings.Join(v, " ")))

	return hex.EncodeToString(h.Sum(nil))
}
