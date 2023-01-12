package errortypes

// BundleNotFound is the error to use when a bundle is not found.
//
// Deprecated: use the genreric services/settings/pkg/settings.NotFound error
type BundleNotFound string

func (e BundleNotFound) Error() string { return "error: bundle not found: " + string(e) }

// IsBundleNotFound implements the IsBundleNotFound interface.
func (e BundleNotFound) IsBundleNotFound() {}
