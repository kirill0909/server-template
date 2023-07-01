package auth

import (
	customErrors "auth-svc/pkg/errors"
	"log"
)

type AuthHeaders struct {
	IP            string `json:"Cf-Connecting-Ip"`
	UserAgent     string `json:"User-Agent"`
	Fingerprint   string `json:"Fingerprint"`
	Authorization string `json:"Authorization"`
}

func (a *AuthHeaders) Validate() error {
	log.Println("-------AuthHeaders.Validate start")
	if a.IP == "" {
		return customErrors.ErrCfIPHeaderAbsent
	}
	if a.UserAgent == "" {
		return customErrors.ErrUserAgentHeaderAbsent
	}
	if a.Fingerprint == "" {
		return customErrors.ErrFingerprintHeaderAbsent
	}
	if a.Authorization == "" {
		return customErrors.ErrAuthTokenCookieAbsent
	}
	log.Println("-------AuthHeaders.Validate end")
	return nil
}
