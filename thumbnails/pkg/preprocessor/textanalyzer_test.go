package preprocessor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzeString(t *testing.T) {
	defaultOpts := AnalysisOpts{
		UseMergeMap: true,
		MergeMap:    DefaultMergeMap,
	}

	tables := []struct {
		input string
		opts  AnalysisOpts
		eOut  TextAnalysis
	}{
		{
			input: "basic latin",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 10, Spaces: []int{5}, TargetScript: "Latin", RuneCount: 11},
				},
				RuneCount: map[string]int{
					"Latin": 11,
				},
				Text: "basic latin",
			},
		},
		{
			input: "trailing tab	",
			opts: defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 12, Spaces: []int{8, 12}, TargetScript: "Latin", RuneCount: 13},
				},
				RuneCount: map[string]int{
					"Latin": 13,
				},
				Text: "trailing tab	",
			},
		},
		{
			input: "Small text. \"$\", \"£\" and \"¥\" are currencies.",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 45, Spaces: []int{5, 11, 16, 21, 25, 30, 34}, TargetScript: "Latin", RuneCount: 44},
				},
				RuneCount: map[string]int{
					"Latin": 44,
				},
				Text: "Small text. \"$\", \"£\" and \"¥\" are currencies.",
			},
		},
		{
			input: "latin with 🖖",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 14, Spaces: []int{5, 10}, TargetScript: "Latin", RuneCount: 12},
				},
				RuneCount: map[string]int{
					"Latin": 12,
				},
				Text: "latin with 🖖",
			},
		},
		{
			input: "기본 한국어",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 15, Spaces: []int{6}, TargetScript: "Hangul", RuneCount: 6},
				},
				RuneCount: map[string]int{
					"Hangul": 6,
				},
				Text: "기본 한국어",
			},
		},
		{
			input: "基本的な日本語",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 20, Spaces: []int{}, TargetScript: "Hiragana", RuneCount: 7},
				},
				RuneCount: map[string]int{
					"Hiragana": 7,
				},
				Text: "基本的な日本語",
			},
		},
		{
			input: "ウーロン茶",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 14, Spaces: []int{}, TargetScript: "Katakana", RuneCount: 5},
				},
				RuneCount: map[string]int{
					"Katakana": 5,
				},
				Text: "ウーロン茶",
			},
		},
		{
			input: "私はエンジニアです",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 26, Spaces: []int{}, TargetScript: "Hiragana", RuneCount: 9},
				},
				RuneCount: map[string]int{
					"Hiragana": 9,
				},
				Text: "私はエンジニアです",
			},
		},
		{
			input: "ティー私はエンジニアです",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 35, Spaces: []int{}, TargetScript: "Hiragana", RuneCount: 12},
				},
				RuneCount: map[string]int{
					"Hiragana": 12,
				},
				Text: "ティー私はエンジニアです",
			},
		},
		{
			input: "私はエンジニアです ティー",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 36, Spaces: []int{27}, TargetScript: "Hiragana", RuneCount: 13},
				},
				RuneCount: map[string]int{
					"Hiragana": 13,
				},
				Text: "私はエンジニアです ティー",
			},
		},
		{
			input: "आधारभूत देवनागरी",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 45, Spaces: []int{21}, TargetScript: "Devanagari", RuneCount: 16},
				},
				RuneCount: map[string]int{
					"Devanagari": 16,
				},
				Text: "आधारभूत देवनागरी",
			},
		},
		{
			input: "mixed 언어 传入 🚀!",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 5, Spaces: []int{5}, TargetScript: "Latin", RuneCount: 6},
					ScriptRange{Low: 6, High: 12, Spaces: []int{12}, TargetScript: "Hangul", RuneCount: 3},
					ScriptRange{Low: 13, High: 24, Spaces: []int{19}, TargetScript: "Han", RuneCount: 5}, // 🚀 and ! are "Common" script and will be merged with "Han"
				},
				RuneCount: map[string]int{
					"Latin":  6,
					"Hangul": 3,
					"Han":    5,
				},
				Text: "mixed 언어 传入 🚀!",
			},
		},
		{
			input: "/k͜p/",
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 5, Spaces: []int{}, TargetScript: "Latin", RuneCount: 5},
				},
				RuneCount: map[string]int{
					"Latin": 5,
				},
				Text: "/k͜p/",
			},
		},
		{
			input: "ä ä", // ä and a + ¨
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 5, Spaces: []int{2}, TargetScript: "Latin", RuneCount: 4},
				},
				RuneCount: map[string]int{
					"Latin": 4,
				},
				Text: "ä ä",
			},
		},
		{
			input: "базовый русский", // cyrillic script isn't part of our default
			opts:  defaultOpts,
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 28, Spaces: []int{14}, TargetScript: "_unknown", RuneCount: 15},
				},
				RuneCount: map[string]int{
					"_unknown": 15,
				},
				Text: "базовый русский",
			},
		},
	}

	for _, table := range tables {
		testname := fmt.Sprintf("Analyzing \"%s\" string", table.input)
		t.Run(testname, func(t *testing.T) {
			ta := NewTextAnalyzer(DefaultScripts)
			result := ta.AnalyzeString(table.input, table.opts)
			if table.opts.UseMergeMap {
				result.MergeCommon(table.opts.MergeMap)
			} else {
				result.MergeCommon(MergeMap{})
			}
			assert.Equal(t, table.eOut, result)
		})
	}
}

func TestAnalyzeStringRaw(t *testing.T) {
	tables := []struct {
		input string
		eOut  TextAnalysis
	}{
		{
			input: "basic latin",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 4, Spaces: []int{}, TargetScript: "Latin", RuneCount: 5},
					ScriptRange{Low: 5, High: 5, Spaces: []int{5}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 6, High: 10, Spaces: []int{}, TargetScript: "Latin", RuneCount: 5},
				},
				RuneCount: map[string]int{
					"Latin":  10,
					"Common": 1,
				},
				Text: "basic latin",
			},
		},
		{
			input: "trailing tab	",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 7, Spaces: []int{}, TargetScript: "Latin", RuneCount: 8},
					ScriptRange{Low: 8, High: 8, Spaces: []int{8}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 9, High: 11, Spaces: []int{}, TargetScript: "Latin", RuneCount: 3},
					ScriptRange{Low: 12, High: 12, Spaces: []int{12}, TargetScript: "Common", RuneCount: 1},
				},
				RuneCount: map[string]int{
					"Latin":  11,
					"Common": 2,
				},
				Text: "trailing tab	",
			},
		},
		{
			input: "Small text. \"$\", \"£\" and \"¥\" are currencies.",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 4, Spaces: []int{}, TargetScript: "Latin", RuneCount: 5},
					ScriptRange{Low: 5, High: 5, Spaces: []int{5}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 6, High: 9, Spaces: []int{}, TargetScript: "Latin", RuneCount: 4},
					ScriptRange{Low: 10, High: 21, Spaces: []int{11, 16, 21}, TargetScript: "Common", RuneCount: 11}, // £ takes 2 bytes
					ScriptRange{Low: 22, High: 24, Spaces: []int{}, TargetScript: "Latin", RuneCount: 3},
					ScriptRange{Low: 25, High: 30, Spaces: []int{25, 30}, TargetScript: "Common", RuneCount: 5}, // ¥ takes 2 bytes
					ScriptRange{Low: 31, High: 33, Spaces: []int{}, TargetScript: "Latin", RuneCount: 3},
					ScriptRange{Low: 34, High: 34, Spaces: []int{34}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 35, High: 44, Spaces: []int{}, TargetScript: "Latin", RuneCount: 10},
					ScriptRange{Low: 45, High: 45, Spaces: []int{}, TargetScript: "Common", RuneCount: 1},
				},
				RuneCount: map[string]int{
					"Latin":  25,
					"Common": 19,
				},
				Text: "Small text. \"$\", \"£\" and \"¥\" are currencies.",
			},
		},
		{
			input: "latin with 🖖",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 4, Spaces: []int{}, TargetScript: "Latin", RuneCount: 5},
					ScriptRange{Low: 5, High: 5, Spaces: []int{5}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 6, High: 9, Spaces: []int{}, TargetScript: "Latin", RuneCount: 4},
					ScriptRange{Low: 10, High: 14, Spaces: []int{10}, TargetScript: "Common", RuneCount: 2},
				},
				RuneCount: map[string]int{
					"Latin":  9,
					"Common": 3,
				},
				Text: "latin with 🖖",
			},
		},
		{
			input: "기본 한국어",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 5, Spaces: []int{}, TargetScript: "Hangul", RuneCount: 2},
					ScriptRange{Low: 6, High: 6, Spaces: []int{6}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 7, High: 15, Spaces: []int{}, TargetScript: "Hangul", RuneCount: 3},
				},
				RuneCount: map[string]int{
					"Hangul": 5,
					"Common": 1,
				},
				Text: "기본 한국어",
			},
		},
		{
			input: "基本的な日本語",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 8, Spaces: []int{}, TargetScript: "Han", RuneCount: 3},
					ScriptRange{Low: 9, High: 11, Spaces: []int{}, TargetScript: "Hiragana", RuneCount: 1},
					ScriptRange{Low: 12, High: 20, Spaces: []int{}, TargetScript: "Han", RuneCount: 3},
				},
				RuneCount: map[string]int{
					"Hiragana": 1,
					"Han":      6,
				},
				Text: "基本的な日本語",
			},
		},
		{
			input: "ウーロン茶",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 2, Spaces: []int{}, TargetScript: "Katakana", RuneCount: 1},
					ScriptRange{Low: 3, High: 5, Spaces: []int{}, TargetScript: "Common", RuneCount: 1}, // ー U+30FC (KATAKANA-HIRAGANA PROLONGED SOUND MARK) seems to be counted as Common
					ScriptRange{Low: 6, High: 11, Spaces: []int{}, TargetScript: "Katakana", RuneCount: 2},
					ScriptRange{Low: 12, High: 14, Spaces: []int{}, TargetScript: "Han", RuneCount: 1},
				},
				RuneCount: map[string]int{
					"Katakana": 3,
					"Common":   1,
					"Han":      1,
				},
				Text: "ウーロン茶",
			},
		},
		{
			input: "私はエンジニアです",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 2, Spaces: []int{}, TargetScript: "Han", RuneCount: 1},
					ScriptRange{Low: 3, High: 5, Spaces: []int{}, TargetScript: "Hiragana", RuneCount: 1},
					ScriptRange{Low: 6, High: 20, Spaces: []int{}, TargetScript: "Katakana", RuneCount: 5},
					ScriptRange{Low: 21, High: 26, Spaces: []int{}, TargetScript: "Hiragana", RuneCount: 2},
				},
				RuneCount: map[string]int{
					"Han":      1,
					"Hiragana": 3,
					"Katakana": 5,
				},
				Text: "私はエンジニアです",
			},
		},
		{
			input: "ティー私はエンジニアです",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 5, Spaces: []int{}, TargetScript: "Katakana", RuneCount: 2},
					ScriptRange{Low: 6, High: 8, Spaces: []int{}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 9, High: 11, Spaces: []int{}, TargetScript: "Han", RuneCount: 1},
					ScriptRange{Low: 12, High: 14, Spaces: []int{}, TargetScript: "Hiragana", RuneCount: 1},
					ScriptRange{Low: 15, High: 29, Spaces: []int{}, TargetScript: "Katakana", RuneCount: 5},
					ScriptRange{Low: 30, High: 35, Spaces: []int{}, TargetScript: "Hiragana", RuneCount: 2},
				},
				RuneCount: map[string]int{
					"Han":      1,
					"Hiragana": 3,
					"Katakana": 7,
					"Common":   1,
				},
				Text: "ティー私はエンジニアです",
			},
		},
		{
			input: "私はエンジニアです ティー",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 2, Spaces: []int{}, TargetScript: "Han", RuneCount: 1},
					ScriptRange{Low: 3, High: 5, Spaces: []int{}, TargetScript: "Hiragana", RuneCount: 1},
					ScriptRange{Low: 6, High: 20, Spaces: []int{}, TargetScript: "Katakana", RuneCount: 5},
					ScriptRange{Low: 21, High: 26, Spaces: []int{}, TargetScript: "Hiragana", RuneCount: 2},
					ScriptRange{Low: 27, High: 27, Spaces: []int{27}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 28, High: 33, Spaces: []int{}, TargetScript: "Katakana", RuneCount: 2},
					ScriptRange{Low: 34, High: 36, Spaces: []int{}, TargetScript: "Common", RuneCount: 1},
				},
				RuneCount: map[string]int{
					"Han":      1,
					"Hiragana": 3,
					"Katakana": 7,
					"Common":   2,
				},
				Text: "私はエンジニアです ティー",
			},
		},
		{
			input: "आधारभूत देवनागरी",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 20, Spaces: []int{}, TargetScript: "Devanagari", RuneCount: 7},
					ScriptRange{Low: 21, High: 21, Spaces: []int{21}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 22, High: 45, Spaces: []int{}, TargetScript: "Devanagari", RuneCount: 8},
				},
				RuneCount: map[string]int{
					"Devanagari": 15,
					"Common":     1,
				},
				Text: "आधारभूत देवनागरी",
			},
		},
		{
			input: "mixed 언어 传入 🚀!",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 4, Spaces: []int{}, TargetScript: "Latin", RuneCount: 5},
					ScriptRange{Low: 5, High: 5, Spaces: []int{5}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 6, High: 11, Spaces: []int{}, TargetScript: "Hangul", RuneCount: 2},
					ScriptRange{Low: 12, High: 12, Spaces: []int{12}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 13, High: 18, Spaces: []int{}, TargetScript: "Han", RuneCount: 2},
					ScriptRange{Low: 19, High: 24, Spaces: []int{19}, TargetScript: "Common", RuneCount: 3},
				},
				RuneCount: map[string]int{
					"Latin":  5,
					"Hangul": 2,
					"Han":    2,
					"Common": 5,
				},
				Text: "mixed 언어 传入 🚀!",
			},
		},
		{
			input: "/k͜p/",
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 0, Spaces: []int{}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 1, High: 1, Spaces: []int{}, TargetScript: "Latin", RuneCount: 1},
					ScriptRange{Low: 2, High: 3, Spaces: []int{}, TargetScript: "Inherited", RuneCount: 1},
					ScriptRange{Low: 4, High: 4, Spaces: []int{}, TargetScript: "Latin", RuneCount: 1},
					ScriptRange{Low: 5, High: 5, Spaces: []int{}, TargetScript: "Common", RuneCount: 1},
				},
				RuneCount: map[string]int{
					"Latin":     2,
					"Common":    2,
					"Inherited": 1,
				},
				Text: "/k͜p/",
			},
		},
		{
			input: "ä ä", // ä and a + ¨
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 1, Spaces: []int{}, TargetScript: "Latin", RuneCount: 1},
					ScriptRange{Low: 2, High: 2, Spaces: []int{2}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 3, High: 3, Spaces: []int{}, TargetScript: "Latin", RuneCount: 1},
					ScriptRange{Low: 4, High: 5, Spaces: []int{}, TargetScript: "Inherited", RuneCount: 1},
				},
				RuneCount: map[string]int{
					"Latin":     2,
					"Common":    1,
					"Inherited": 1,
				},
				Text: "ä ä",
			},
		},
		{
			input: "базовый русский", // cyrillic script isn't part of our default
			eOut: TextAnalysis{
				ScriptRanges: []ScriptRange{
					ScriptRange{Low: 0, High: 13, Spaces: []int{}, TargetScript: "_unknown", RuneCount: 7},
					ScriptRange{Low: 14, High: 14, Spaces: []int{14}, TargetScript: "Common", RuneCount: 1},
					ScriptRange{Low: 15, High: 28, Spaces: []int{}, TargetScript: "_unknown", RuneCount: 7},
				},
				RuneCount: map[string]int{
					"_unknown": 14,
					"Common":   1,
				},
				Text: "базовый русский",
			},
		},
	}

	for _, table := range tables {
		testname := fmt.Sprintf("Raw-Analyzing \"%s\" string", table.input)
		t.Run(testname, func(t *testing.T) {
			ta := NewTextAnalyzer(DefaultScripts)
			result := ta.AnalyzeString(table.input, AnalysisOpts{})

			assert.Equal(t, table.eOut, result)
		})
	}
}
