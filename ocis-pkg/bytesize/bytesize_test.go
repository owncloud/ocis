package bytesize_test

import (
	"fmt"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/bytesize"
	"github.com/test-go/testify/require"
)

func TestParseSpecial(t *testing.T) {
	testCases := []struct {
		Alias          string
		Input          string
		ExpectedOutput uint64
		ExpectError    bool
	}{
		{
			Alias:          "it assumes bytes",
			Input:          "100",
			ExpectedOutput: 100,
		},
		{
			Alias:          "it accepts a space between value and unit",
			Input:          "1 MB",
			ExpectedOutput: 1000000,
		},
		{
			Alias:          "it accepts also more spaces between value and unit",
			Input:          "1                                            MB",
			ExpectedOutput: 1000000,
		},
		{
			Alias:          "it ignores trailing and leading spaces",
			Input:          " 1MB ",
			ExpectedOutput: 1000000,
		},
		{
			Alias:          "it errors on unknown units",
			Input:          "1SB",
			ExpectedOutput: 0,
			ExpectError:    true,
		},
		{
			Alias:          "it multiplies correctly",
			Input:          "16MB",
			ExpectedOutput: 16000000,
		},
		{
			Alias:          "it errors when no value is given",
			Input:          "GB",
			ExpectedOutput: 0,
			ExpectError:    true,
		},
		{
			Alias:          "it errors when bad input is given",
			Input:          ",as!@@delta",
			ExpectedOutput: 0,
			ExpectError:    true,
		},
		{
			Alias:          "it errors when using floats",
			Input:          "1.024GB",
			ExpectedOutput: 0,
			ExpectError:    true,
		},
	}

	for _, tc := range testCases {
		actual, err := bytesize.Parse(tc.Input)
		if tc.ExpectError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		require.Equal(t, tc.ExpectError, err != nil, tc.Alias)
		require.Equal(t, int(tc.ExpectedOutput), int(actual), tc.Alias)
	}
}

func TestParseHappy(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected uint64
	}{
		{Input: "1", Expected: 1},
		{Input: "1KB", Expected: 1000},
		{Input: "1MB", Expected: 1000000},
		{Input: "1GB", Expected: 1000000000},
		{Input: "1TB", Expected: 1000000000000},
		{Input: "1PB", Expected: 1000000000000000},
		{Input: "1EB", Expected: 1000000000000000000},
		{Input: "1MiB", Expected: 1048576},
		{Input: "1GiB", Expected: 1073741824},
		{Input: "1TiB", Expected: 1099511627776},
		{Input: "1PiB", Expected: 1125899906842624},
		{Input: "1EiB", Expected: 1152921504606846976},
	}

	for _, tc := range testCases {
		actual, err := bytesize.Parse(tc.Input)
		require.NoError(t, err)
		require.Equal(t, int(tc.Expected), int(actual), fmt.Sprintf("case %s", tc.Input))
	}
}
