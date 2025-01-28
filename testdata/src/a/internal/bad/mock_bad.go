package bad

// want "mock types are not allowed in internal/ directories"
type MockBad struct{}
