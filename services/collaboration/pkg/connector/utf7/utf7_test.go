package utf7_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/connector/utf7"
)

var _ = Describe("Utf7", func() {
	DescribeTable(
		"Encode",
		func(input, output string) {
			Expect(utf7.EncodeString(input)).To(Equal(output))
		},
		Entry("regular filename", "private.txt", "private.txt"),
		Entry("direct chars", "is-better?yes:no(3).pdf", "is-better?yes:no(3).pdf"),
		Entry("with spaces", "a big file with spaces.docx", "a+ACA-big+ACA-file+ACA-with+ACA-spaces.docx"),
		Entry("with symbols", "a+b=c", "a+ACs-b+AD0-c"),
		Entry("unicode filenames", "超極秘文書.doc", "+jYVpdXnYZYdm+A-.doc"),
		Entry("emoji and symbols", "💰🔜™️.pdf", "+2D3csNg93RwhIv4P-.pdf"),
	)

	DescribeTable(
		"Decode",
		func(input, output string) {
			Expect(utf7.DecodeString(input)).To(Equal(output))
		},
		Entry("regular filename", "private.txt", "private.txt"),
		Entry("direct chars", "is-better?yes:no(3).pdf", "is-better?yes:no(3).pdf"),
		Entry("with spaces", "a+ACA-big+ACA-file+ACA-with+ACA-spaces.docx", "a big file with spaces.docx"),
		Entry("with symbols", "a+ACs-b+AD0-c", "a+b=c"),
		Entry("unicode filenames", "+jYVpdXnYZYdm+A-.doc", "超極秘文書.doc"),
		Entry("emoji and symbols", "+2D3csNg93RwhIv4P-.pdf", "💰🔜™️.pdf"),
		Entry("special case +- chars", "a+-b+AD0-c", "a+b=c"),
		Entry("optional direct chars", "1 +- 1 +AD0- 2", "1 + 1 = 2"),
		Entry("missing - char", "+jYVpdXnYZYdm+A.doc", "超極秘文書.doc"),
		Entry("missing - char2", "a+AD0.b", "a=.b"),
		Entry("missing - char end", "88+xUixVdVYwTjGlA", "88안녕하세요"),
	)

	DescribeTable(
		"EncodeDecode",
		func(input string) {
			output, err := utf7.DecodeString(utf7.EncodeString(input))
			Expect(err).To(Succeed())
			Expect(output).To(Equal(input))
		},
		Entry("regular filename", "private.txt"),
		Entry("direct chars", "is-better?yes:no(3).pdf"),
		Entry("with spaces", "a big file with spaces.docx"),
		Entry("with symbols", "a+b=c"),
		Entry("unicode filenames", "超極秘文書.doc"),
		Entry("emoji and symbols", "💰🔜™️.pdf"),
	)

	DescribeTable(
		"DecodeFailures",
		func(input string) {
			out, err := utf7.DecodeString(input)
			Expect(err).To(HaveOccurred())
			Expect(out).To(Equal(""))
		},
		Entry("Non-utf16 sequence in base64", "a+-b+AD-c"),
		Entry("Illegal base64 char", "a+-b+A.D-c"),
		Entry("Non-ascii string", "⇉a+-b"),
	)
})
