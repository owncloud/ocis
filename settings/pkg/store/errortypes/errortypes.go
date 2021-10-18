package errortypes

// BundleNotFound is the error to use when a bundle is not found.
type BundleNotFound string

func (e BundleNotFound) Error() string { return "error: bundle not found: " + string(e) }

// IsBundleNotFound implements the IsBundleNotFound interface.
func (e BundleNotFound) IsBundleNotFound() {}
