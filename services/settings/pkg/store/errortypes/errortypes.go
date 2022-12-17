package errortypes

// BundleNotFoundError is the error to use when a bundle is not found.
type BundleNotFoundError string

func (e BundleNotFoundError) Error() string { return "error: bundle not found: " + string(e) }

// IsBundleNotFound implements the IsBundleNotFound interface.
func (e BundleNotFoundError) IsBundleNotFound() {}
