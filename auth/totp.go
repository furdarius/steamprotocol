package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	chars    = "23456789BCDFGHJKMNPQRTVWXY"
	charsLen = uint32(len(chars))
)

type TimeTip struct {
	Time                              int64  `json:"server_time,string"`
	SkewToleranceSeconds              uint32 `json:"skew_tolerance_seconds,string"`
	LargeTimeJink                     uint32 `json:"large_time_jink,string"`
	ProbeFrequencySeconds             uint32 `json:"probe_frequency_seconds"`
	AdjustedTimeProbeFrequencySeconds uint32 `json:"adjusted_time_probe_frequency_seconds"`
	HintProbeFrequencySeconds         uint32 `json:"hint_probe_frequency_seconds"`
	SyncTimeout                       uint32 `json:"sync_timeout"`
	TryAgainSeconds                   uint32 `json:"try_again_seconds"`
	MaxAttempts                       uint32 `json:"max_attempts"`
}

// TOTPGenerator is used to generate two factor auth code synced by time with Steam.
// TOTP = Time-based One-time Password Algorithm
type TOTPGenerator struct {
	cl *http.Client
}

// NewTOTPGenerator initialize new instance of TOTPGenerator
func NewTOTPGenerator(cl *http.Client) *TOTPGenerator {
	return &TOTPGenerator{cl}
}

// TwoFactorSynced fetch time from Steam and return generated two-factor code synced by time with Steam
func (gen *TOTPGenerator) TwoFactorSynced(sharedSecret string) (string, error) {
	timeTip, err := gen.FetchTimeTip()
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch time tip")
	}

	return gen.GenerateTwoFactorCode(sharedSecret, timeTip.Time)
}

// GenerateTwoFactorCode generate Steam two-factor code using current timestamp as parameter.
func (gen *TOTPGenerator) GenerateTwoFactorCode(sharedSecret string, currentTimestamp int64) (string, error) {
	data, err := base64.StdEncoding.DecodeString(sharedSecret)
	if err != nil {
		return "", err
	}

	ful := make([]byte, 8)
	binary.BigEndian.PutUint32(ful[4:], uint32(currentTimestamp/30))

	hmacBuf := hmac.New(sha1.New, data)
	hmacBuf.Write(ful)

	sum := hmacBuf.Sum(nil)
	start := sum[19] & 0x0F
	slice := binary.BigEndian.Uint32(sum[start:start+4]) & 0x7FFFFFFF

	buf := make([]byte, 5)
	for i := 0; i < 5; i++ {
		buf[i] = chars[slice%charsLen]
		slice /= charsLen
	}

	return string(buf), nil
}

// FetchTimeTip fetch time from Steam.
func (gen *TOTPGenerator) FetchTimeTip() (*TimeTip, error) {
	resp, err := gen.cl.Post(
		"https://api.steampowered.com/ITwoFactorService/QueryTime/v1/",
		"application/x-www-form-urlencoded",
		nil,
	)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Inner *TimeTip `json:"response"`
	}

	var response Response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	resp.Body.Close()

	return response.Inner, nil
}
