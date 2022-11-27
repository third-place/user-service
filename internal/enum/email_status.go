package enum

type EmailStatusType string

const (
	EmailStatusVerified EmailStatusType = "unverified"
	EmailStatusUnverified EmailStatusType = "verified"
)
