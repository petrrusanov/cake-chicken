package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc64"
	"io"
	"log"
	"regexp"
)

var reValidSha = regexp.MustCompile("^[a-fA-F0-9]{40}$")

// StrongHashValue makes SHA-512 hmac with secret
func StrongHashValue(val string, secret string) string {
	if val == "" || reValidSha.MatchString(val) {
		return val // already hashed or empty
	}
	key := []byte(secret)

	h := hmac.New(sha512.New, key)
	return hashWithFallback(h, val)
}

// hashWithFallback tries to has val with hash.Hash and failback to crc if needed
func hashWithFallback(h hash.Hash, val string) string {
	if _, err := io.WriteString(h, val); err != nil {
		// fail back to crc64
		log.Printf("[WARN] can't hash id %s, %s", val, err)
		return fmt.Sprintf("%x", crc64.Checksum([]byte(val), crc64.MakeTable(crc64.ECMA)))
	}
	return hex.EncodeToString(h.Sum(nil))
}
