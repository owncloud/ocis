package checkers

// NoopChecker doesn't check anything
type NoopChecker struct {
}

// NewNoopChecker creates a new NoopChecker
func NewNoopChecker() *NoopChecker {
	return &NoopChecker{}
}

// CheckClaims won't do anything and won't return an error
func (nc *NoopChecker) CheckClaims(_ map[string]interface{}) error {
	return nil
}

func (nc *NoopChecker) RequireMap() map[string]string {
	return nil
}
