Enhancement: Correct shutdown of services under runtime

Supervised goroutines now shut themselves down on context cancellation propagation.

https://github.com/owncloud/ocis/pull/2843
