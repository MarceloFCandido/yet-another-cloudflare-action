package utils

// Deprecated: Use HandleError from error-handler.go instead
// This function is kept for backward compatibility but should not be used
var PanicOnError = panicOnError

// Deprecated: Use HandleError from error-handler.go instead
func panicOnError(err error) {
	if err != nil {
		// For backward compatibility with tests that expect panic
		panic(err)
	}
}
