package audit

func NewUserSignupEvent() *AuditEntry {
	return &AuditEntry{
		Action: ActionUserSignup,
	}
}
