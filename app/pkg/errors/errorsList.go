package errors

var (
	// Auth middleware
	ErrAuthTokenCookieAbsent   = NewCustomError("No Authorization token provided in Cookies", 401, true)
	ErrCfIPHeaderAbsent        = NewCustomError("Ip header validation err", 401, true)
	ErrUserAgentHeaderAbsent   = NewCustomError("User-Agent header validation err", 401, true)
	ErrFingerprintHeaderAbsent = NewCustomError("Fingerprint header validation err", 401, true)

	// User
	ErrGetUserIDFromCtx = NewCustomError("Cannot get userID from context", 401, true)
)
