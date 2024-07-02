package utf7

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"strings"
	"unicode"
	"unicode/utf16"
)

const (
	rangeASCII = "ascii"
	rangeUTF7  = "utf7"
)

// Range represents a range with a lower and upper bounds. The range has a
// name for easier identification
type Range struct {
	Name string
	Low  int
	High int
}

// Range table for ASCII chars belonging to the "direct character" group
// for UTF-7
var utf7AsciiRT = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x27, 0x29, 1}, // '()
		{0x2c, 0x2f, 1}, // ,-./
		{0x30, 0x39, 1}, // 0-9
		{0x3a, 0x3f, 5}, // :?
		{0x41, 0x5a, 1}, // A-Z
		{0x61, 0x7a, 1}, // a-z
	},
}

// EncodeString will encode the provided UTF-8 string into UTF-7 format
//
// The encoding process will have the following peculiarities
// * Any char outside the "direct characters" will be encoded. This means that
// only "a-z", "A-Z", "0-9" and "'(),-.:?" chars will remain intact while the
// rest will be encoded. "Optional direct chars" (such as the space) will
// be encoded.
// * The "+" char will be encoded as any other character, so the result will
// be "+ACs-", not "+-"
// * Sequences of chars will be encoded as a single group. For example,
// "こんにちは" will be encoded as "+MFMwkzBrMGEwbw-"
// * All encoded sequences will be enclosed between "+" and "-"
func EncodeString(s string) string {
	runes := []rune(s)

	ranges := analyzeRunes(runes)

	var sb strings.Builder
	// doubling the number of bytes of the string is usually enough
	sb.Grow(len(s) * 2)

	for _, v := range ranges {
		if v.Name == rangeASCII {
			for _, v := range runes[v.Low:v.High] {
				sb.WriteRune(v)
			}
		} else {
			utf7Bytes := convertToUtf7(runes[v.Low:v.High])
			sb.Write(utf7Bytes)
		}
	}
	return sb.String()
}

// DecodeString will decode the provided UTF-7 string into UTF-8.
//
// Any valid UTF-7 string can be decoded, not just the ones returned by
// the EncodeString function.
// In particular, UTF-7 strings such as "a+-b" or "a+AD0.b" can be decoded
// even if the EncodeString function won't generate the corresponding
// strings that way.
//
// Note that this function requires the string to contain only ASCII chars
// (as per UTF-7), otherwise an error will be returned.
// Illegal char sequences in the encoded parts of the string will also trigger
// errors.
func DecodeString(s string) (string, error) {
	byteArray := []byte(s)

	ranges, err := analyzeUtf7(byteArray)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.Grow(len(byteArray))

	for _, v := range ranges {
		if v.Name == rangeASCII {
			// if it's an ascii range, just copy it
			sb.Write(byteArray[v.Low:v.High])
		} else {
			// utf7 range
			utf7ByteRange := byteArray[v.Low:v.High]
			if len(utf7ByteRange) == 2 && utf7ByteRange[0] == '+' && utf7ByteRange[1] == '-' {
				// special case for the "+-" sequence -> just write "+" as replacement
				sb.WriteByte('+')
			} else {
				// utf7 range must start with "+" and should (but might not) end with "-"
				// we need to remove those chars before decoding
				toDecode := byteArray[v.Low+1 : v.High-1]
				if byteArray[v.High-1] != '-' {
					toDecode = byteArray[v.Low+1 : v.High]
				}
				runeArray, err := convertFromUtf7(toDecode)
				if err != nil {
					return "", err
				}
				for _, r := range runeArray {
					sb.WriteRune(r)
				}
			}
		}
	}
	return sb.String(), nil
}

// analyzeRunes will analyze the array of runes and provide a list of ranges.
// Each range will be defined by a name and a low and high index. For example,
// an "ascii" range could go from index 0 to 12 and "utf7" range from 12 to 25.
// The range includes the low index but not the high "[0,12)". This means it
// be easily extracted with something like "runes[r.Low:r.High]".
//
// The list of ranges will only include the following names:
// * "ascii" for runes belonging to the "direct characters" group of UTF-7
// (those that can be used directly without encoding them). Note that
// it won't consider every ASCII character.
// * "utf7" for runes that should be encoded for UTF-7.
//
// As said, runes in the ranges marked as "utf7" should be encoded for UTF-7,
// while the others can be used without changes.
//
// This method is intended to be used to detect which ranges need to be
// encoded to UTF-7
func analyzeRunes(runes []rune) []Range {
	ranges := make([]Range, 0)

	var currentRange Range
	for k, v := range runes {
		if unicode.Is(utf7AsciiRT, v) {
			if currentRange.Name == "" {
				// take control of the current range
				currentRange.Name = rangeASCII
				currentRange.Low = k
			} else if currentRange.Name != rangeASCII {
				// close current range and open a new one
				currentRange.High = k
				ranges = append(ranges, currentRange)
				currentRange = Range{
					Name: rangeASCII,
					Low:  k,
				}
			}
		} else {
			if currentRange.Name == "" {
				// take control of the current range
				currentRange.Name = rangeUTF7
				currentRange.Low = k
			} else if currentRange.Name != rangeUTF7 {
				// close current range and open a new one
				currentRange.High = k
				ranges = append(ranges, currentRange)
				currentRange = Range{
					Name: rangeUTF7,
					Low:  k,
				}
			}
		}
	}
	// close the last range
	currentRange.High = len(runes)
	ranges = append(ranges, currentRange)

	return ranges
}

// analyzeUtf7 will analyze the provided byte sequence and return a list of
// ranges.
// The byte sequence is considered as UTF-7, so if there is a non-ASCII char
// in the sequence, an error will be returned (it isn't a valid UTF-7 string).
//
// Each returned range will have either "ascii" or "utf7" as name for the range.
// "ascii" ranges won't require any change and can be used directly. "utf7"
// ranges are encoded in UTF-7 and will require decoding.
//
// This method is intended to be used to detect which ranges need to be
// decoded from UTF-7
func analyzeUtf7(byteArray []byte) ([]Range, error) {
	base64chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	base64ByteArray := []byte(base64chars)

	ranges := make([]Range, 0)

	currentRange := Range{
		Name: rangeASCII,
		Low:  0,
	}

	for k, v := range byteArray {
		if v > unicode.MaxASCII {
			return nil, errors.New("Byte sequence contains a non-ASCII char")
		}

		if v == '+' && currentRange.Name != rangeUTF7 {
			// start utf7-encoded range
			currentRange.High = k
			ranges = append(ranges, currentRange)
			currentRange = Range{
				Name: rangeUTF7,
				Low:  k,
			}
		} else if v == '-' {
			// close utf7-encoded range
			currentRange.High = k + 1 // the '-' char is part of the range
			ranges = append(ranges, currentRange)
			currentRange = Range{
				Name: rangeASCII,
				Low:  k + 1,
			}
		} else if bytes.IndexByte(base64ByteArray, v) == -1 && currentRange.Name == rangeUTF7 {
			// found invalid base64 char, so need to close the utf7 range
			currentRange.High = k
			ranges = append(ranges, currentRange)
			currentRange = Range{
				Name: rangeASCII,
				Low:  k,
			}
		}
	}

	// close the last range
	currentRange.High = len(byteArray)
	ranges = append(ranges, currentRange)

	// there might be empty ranges we need to clear
	// empty ranges have Low = High
	realRanges := make([]Range, 0, len(ranges))
	for _, v := range ranges {
		if v.Low != v.High {
			realRanges = append(realRanges, v)
		}
	}

	return realRanges, nil
}

// convertToUtf7 will convert the provided runes to a UTF-7 sequence of bytes.
// The function assumes that all the provided runes must be converted to UTF-7
func convertToUtf7(runes []rune) []byte {
	byteArray := make([]byte, 0, len(runes)*2)

	u16 := utf16.Encode(runes)
	for _, v := range u16 {
		byteArray = binary.BigEndian.AppendUint16(byteArray, v)
	}

	dst := make([]byte, base64.RawStdEncoding.EncodedLen(len(byteArray))+2)
	dst[0] = '+'
	base64.RawStdEncoding.Encode(dst[1:len(dst)-1], byteArray)
	dst[len(dst)-1] = '-'
	return dst
}

// convertFromUtf7 will convert the sequence of bytes to runes. The sequence
// of bytes is assumed to be an UTF-7 encoded sequence (without the "+" and
// "-" limiters)
// The returned runes should be UTF-8 encoded and can be converted to a
// regular string easily.
// Note that errors can be returned if the decoding process fails
func convertFromUtf7(byteArray []byte) ([]rune, error) {
	dst := make([]byte, base64.RawStdEncoding.DecodedLen(len(byteArray)))

	_, err := base64.RawStdEncoding.Decode(dst, byteArray)
	if err != nil {
		return []rune{}, err
	}

	if len(dst)%2 != 0 {
		// some data can't be represented as utf16, and can't be decoded
		return []rune{}, errors.New("some utf7 data can't be represented as utf16")
	}

	u16array := make([]uint16, 0, len(dst)/2)
	for i := 0; i < len(dst); i++ {
		u16array = append(u16array, binary.BigEndian.Uint16(dst[i:i+2]))
		i = i + 1
	}
	return utf16.Decode(u16array), nil
}
