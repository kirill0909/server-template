package auth

import (
	"time"
)

// for redis
type Session struct {
	ClientID     string `json:"clientID"`
	IP           string `json:"ip"`
	UserAgent    string `json:"userAgent"`
	Fingerprint  string `json:"fingerprint"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpireAt     int    `json:"expireAt"`
	Duration     time.Duration
}

func (s *Session) IsAuthDataValid(authData AuthHeaders) bool {
	if s.IP != authData.IP {
		return false
	}
	if s.UserAgent != authData.UserAgent {
		return false
	}
	if s.Fingerprint != authData.Fingerprint {
		return false
	}
	if s.AccessToken != authData.Authorization {
		return false
	}
	return true
}
