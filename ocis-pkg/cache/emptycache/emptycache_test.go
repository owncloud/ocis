package emptycache

import (
	"fmt"
	"testing"
)

func TestInitializeMultiple(t *testing.T) {
	var tests = []map[string]interface{}{
		{},
		{
			"capacity": 1024,
		},
		{
			"randomKey":     "RandomValue",
			"notConsidered": true,
			"anotherType":   map[string]string{},
		},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			eCache := NewEmptyCache()
			err := eCache.Initialize(test)
			if err != nil {
				t.Errorf("Initialization failed for params %+v, got error %v", test, err)
			}
		})
	}
}

func TestStoreMultiple(t *testing.T) {
	var tests = []struct {
		key   string
		value string
		ttl   int64
	}{
		{"key1", "value1", 222},
		{"üñòKey", "välòç", 555},
		{"5566", "välòç", 1555},
		{"üñòKey", "98685", 555},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			eCache := NewEmptyCache()
			eCache.Initialize(map[string]interface{}{})
			err := eCache.Store(test.key, test.value, test.ttl)
			if err != nil {
				t.Errorf("Storing failed for params %+v, got error %v", test, err)
			}
		})
	}
}

func TestRetrieveMultiple(t *testing.T) {
	type expectedResult struct {
		value  string
		exists bool
		err    error
	}

	var tests = []struct {
		key      string
		expected expectedResult
	}{
		{"key1", expectedResult{"", false, nil}},
		{"üñòKey", expectedResult{"", false, nil}},
		{"5566", expectedResult{"", false, nil}},
		{"üñòKey", expectedResult{"", false, nil}},
		{"missing", expectedResult{"", false, nil}},
	}

	eCache := NewEmptyCache()
	eCache.Initialize(map[string]interface{}{})
	eCache.Store("key1", "value1", 222)
	eCache.Store("üñoKey", "välòç", 555)
	eCache.Store("5566", "välòç", 1555)
	eCache.Store("üñòKey", "98685", 55)
	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			value, exists, err := eCache.Retrieve(test.key)
			if value != test.expected.value {
				t.Errorf("Wrong value for params %+v, expected %s, got %s", test, test.expected.value, value)
			}
			if exists != test.expected.exists {
				t.Errorf("Wrong exist for params %+v, expected %t, got %t", test, test.expected.exists, exists)
			}
			if err != test.expected.err {
				t.Errorf("Wrong error for params %+v, expected %v, got %v", test, test.expected.err, err)
			}
		})
	}
}

func TestRemoveMultiple(t *testing.T) {
	var tests = []string{
		"key1",
		"üñòKey",
		"5566",
		"üñòKey",
		"missing",
	}

	eCache := NewEmptyCache()
	eCache.Initialize(map[string]interface{}{})
	eCache.Store("key1", "value1", 222)
	eCache.Store("üñoKey", "välòç", 555)
	eCache.Store("5566", "välòç", 1555)
	eCache.Store("üñòKey", "98685", 55)
	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			err := eCache.Remove(test)
			if err != nil {
				t.Errorf("Storing failed for params %+v, got error %v", test, err)
			}
		})
	}
}
