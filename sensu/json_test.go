package sensu

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_JsonExtend(t *testing.T) {
	for i, tuple := range []struct {
		src           string
		dst           string
		expected      string
		errorExpected bool
	}{
		{
			src:           `{}`,
			dst:           `{}`,
			expected:      `{}`,
			errorExpected: false,
		},
		{
			src:           `{"b":2}`,
			dst:           `{"a":1}`,
			expected:      `{"a":1,"b":2}`,
			errorExpected: false,
		},
		{
			src:           `{"a":0}`,
			dst:           `{"a":1}`,
			expected:      `{"a":0}`,
			errorExpected: false,
		},
		{
			src:           `{"a":{       "y":2}}`,
			dst:           `{"a":{"x":1       }}`,
			expected:      `{"a":{"x":1, "y":2}}`,
			errorExpected: false,
		},
		{
			src:           `{"a":{"x":2}}`,
			dst:           `{"a":{"x":1}}`,
			expected:      `{"a":{"x":2}}`,
			errorExpected: false,
		},
		{
			src:           `{"a":{       "y":7, "z":8}}`,
			dst:           `{"a":{"x":1, "y":2       }}`,
			expected:      `{"a":{"x":1, "y":7, "z":8}}`,
			errorExpected: false,
		},
		{
			src:           `{"1":{"n":[1,2]}}`,
			dst:           `{"1":{"n":"xxx"}}`,
			expected:      `{}`,
			errorExpected: true,
		},
		{
			src:           `{"1": { "b":1, "2": { "3": {         "b":3, "n":[1,2    ]} }        }}`,
			dst:           `{"1": {        "2": { "3": {"a":"A",        "n":[    3,4]} }, "a":3 }}`,
			expected:      `{"1": { "b":1, "2": { "3": {"a":"A", "b":3, "n":[1,2,3,4]} }, "a":3 }}`,
			errorExpected: false,
		},
		{
			src:           `{"1": { "b":1, "2": { "3": {"a":"A", "b":3, "n":[1,2    ]} }        }}`,
			dst:           `{"1": {        "2": { "3": {"a":"B",        "n":[    3,4]} }, "a":3 }}`,
			expected:      `{"1": { "b":1, "2": { "3": {"a":"A", "b":3, "n":[1,2,3,4]} }, "a":3 }}`,
			errorExpected: false,
		},
	} {
		var src map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.src), &src); err != nil {
			t.Error(err)
			continue
		}

		var dst map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.dst), &dst); err != nil {
			t.Error(err)
			continue
		}

		var expected map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.expected), &expected); err != nil {
			t.Error(err)
			continue
		}

		s := &Json{src}
		d := &Json{dst}
		e := &Json{expected}

		err := s.Extend(d)
		if tuple.errorExpected {
			if err == nil {
				t.Error("Error expected, none returned")
				continue
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}

			if !reflect.DeepEqual(s, e) {
				t.Errorf("%d. expected %v, got %v", i, e, s)
			}
		}
	}
}
