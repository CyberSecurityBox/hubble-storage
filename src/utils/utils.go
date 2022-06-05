package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

func sha256Sum(s []byte) string {
	h := sha256.New()
	h.Write(s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func FlowHash(s []string) string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(s)
	return sha256Sum(b.Bytes())
}
