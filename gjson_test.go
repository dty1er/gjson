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
		{in: `{}  1`, invalid: true},
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
		{in: `{"key": 1a}`, invalid: true},
		{in: `{"key": --1}`, invalid: true},
		{in: `{"key": +0}`, invalid: true},
		{in: `{5: 0}`, invalid: true},
		// bool
		{in: `{"key": true}`},
		{in: `{"key": false}`},
		{in: `{"key": trux}`, invalid: true},
		{in: `{"key": falsx}`, invalid: true},
		{in: `{"key": tru}`, invalid: true},
		{in: `{"key": fals}`, invalid: true},
	} {
		expected := make(map[string]interface{})
		err1 := json.Unmarshal([]byte(testcase.in), &expected)
		out, err2 := Decode([]byte(testcase.in))
		if testcase.invalid {
			if err1 == nil || err2 == nil {
				t.Errorf("error is expected: err1: %v, err2: %v", err1, err2)
			}
			continue
		}

		if err1 != nil {
			t.Errorf("err1 is not nil: %v", err1)
		}

		if err2 != nil {
			t.Errorf("err2 is not nil: %v", err2)
		}
		if !reflect.DeepEqual(out, expected) {
			t.Errorf("#%d: got %v, want %v", i, out, expected)
		}
	}
}
