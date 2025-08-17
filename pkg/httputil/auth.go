package httputil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func buildHMACHeader(token string) string {
	ts := time.Now().Unix()
	mac := hmac.New(sha256.New, []byte(token))
	fmt.Fprintf(mac, "%d", ts)
	sig := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("hmac %d:%s", ts, sig)
}
