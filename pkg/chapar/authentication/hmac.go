package authentication

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func VerifyHMAC(header, token string) error {
	if !strings.HasPrefix(header, "hmac ") {
		return errors.New("invalid header prefix")
	}

	auth := strings.TrimPrefix(header, "hmac ")
	parts := strings.SplitN(auth, ":", 2)
	if len(parts) != 2 {
		return errors.New("invalid format")
	}

	tsStr, sig := parts[0], parts[1]
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return errors.New("invalid timestamp")
	}

	if time.Since(time.Unix(ts, 0)).Abs() > time.Minute {
		return errors.New("expired")
	}

	mac := hmac.New(sha256.New, []byte(token))
	fmt.Fprintf(mac, "%d", ts)
	expected := hex.EncodeToString(mac.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(sig), []byte(expected)) == 1 {
		return nil
	}

	return errors.New("unauthorized")
}
