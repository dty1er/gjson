package gjson

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	for i, testcase := range []struct {
		invalid bool
		in      string
	}{
		{in: `{}`},
		{in: `  {}`},
		{in: `{  }`},
		{in: `{}  `},
		{in: `  {  }  `},
		{in: `a`, invalid: true},
		{in: `{"a:`, invalid: true},
		{in: `{}  1`, invalid: true},
		{in: ``, invalid: true},
		// string
		{in: `{"key": "value"}`},
		{in: `  {"key": "value"}`},
		{in: `  {  "key"  :    "value"   }    `},
		{in: `{key": "value"}`, invalid: true},
		{in: `{"key: "value"}`, invalid: true},
		{in: `{"key": value"}`, invalid: true},
		{in: `{"key": "value}`, invalid: true},
		{in: `{"key"; "value"}`, invalid: true},
		// number
		{in: `{"key": 1}`},
		{in: `{"key": 1.234}`},
		{in: `{"key": -23.45}`},
		{in: `{"key": 0.5}`},
		{in: `{"key": 1a}`, invalid: true},
		{in: `{"key": --1}`, invalid: true},
		{in: `{"key": -a}`, invalid: true},
		{in: `{"key": +0}`, invalid: true},
		{in: `{"key": 0.a)}`, invalid: true},
		{in: `{5: 0}`, invalid: true},
		// bool
		{in: `{"key": true}`},
		{in: `{"key": false}`},
		{in: `{"key": trux}`, invalid: true},
		{in: `{"key": falsx}`, invalid: true},
		{in: `{"key": tru}`, invalid: true},
		{in: `{"key": t}`, invalid: true},
		{in: `{"key": fals}`, invalid: true},
		{in: `{"key": f}`, invalid: true},
		// array
		{in: `{"key": []}`},
		{in: `{"key": ["a"]}`},
		{in: `{"key": ["a", "b"]}`},
		{in: `{"key": ["a", 1]}`},
		{in: `{"key": ["a", -1]}`},
		{in: `{"key": ["a", -1]}`},
		{in: `{"key": [true, false, [true]]}`},
		{in: `{"key": [["a", 4.2], [-1, ["b", "c"]]]}`},
		{in: `{"key": [["a", 4.2], [-1, ["b", "c"]]}`, invalid: true},
		{in: `{"key": [["a", 4.2], "bcd", 123, [-1, ["b", "c"]]}`, invalid: true},
		{in: `{"key": [}`, invalid: true},
	} {
		expected := make(map[string]interface{})
		err1 := json.Unmarshal([]byte(testcase.in), &expected)
		out, err2 := Decode([]byte(testcase.in))
		if testcase.invalid {
			if err1 == nil || err2 == nil {
				t.Fatalf("case #%d:%s: error is expected: err1: %v, err2: %v", i, testcase.in, err1, err2)
			}
			continue
		}

		if err1 != nil {
			t.Fatalf("case #%d:%s: err1 is not nil: %v", i, testcase.in, err1)
		}

		if err2 != nil {
			t.Fatalf("case #%d:%s: err2 is not nil: %v", i, testcase.in, err1)
		}
		if !reflect.DeepEqual(out, expected) {
			t.Errorf("#%d: got %v, want %v", i, out, expected)
		}
	}
}
