package middleware

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("authentication helpers", func() {
	DescribeTable("isPublicPath should recognize public paths",
		func(input string, expected bool) {
			isPublic := isPublicPath(input)
			Expect(isPublic).To(Equal(expected))
		},
		Entry("public files path", "/remote.php/dav/public-files/", true),
		Entry("public files path without remote.php", "/remote.php/dav/public-files/", true),
		Entry("token info path", "/remote.php/ocs/apps/files_sharing/api/v1/tokeninfo/unprotected", true),
		Entry("capabilities", "/ocs/v1.php/cloud/capabilities", true),
	)
})
