package gjson

import (
	"fmt"
)

// Decoder ...
type Decoder struct {
	pos  int
	end  int
	data []byte
}

// NewDecoder ...
func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		pos:  0,
		end:  len(data),
		data: data,
	}
}

func (d *Decoder) skipSpaces() byte {
	for {
		if d.pos == d.end {
			return 0
		}
		switch c := d.data[d.pos]; c {
		case ' ':
			d.pos++
			continue
		default:
			return c
		}
	}
}

func (d *Decoder) decodeObject() (map[string]interface{}, error) {
	d.pos++

	var c byte
	obj := make(map[string]interface{})

	// look ahead for } - if the object has no keys.
	if c = d.skipSpaces(); c == '}' {
		d.pos++
		return obj, nil
	}
	return nil, fmt.Errorf("\"}\" expected, but got %v", c)
}

// Decode ...
func Decode(data []byte) (val map[string]interface{}, err error) {
	d := NewDecoder(data)
	debug("d.end: %v", d.end)

	if c := d.skipSpaces(); c != '{' {
		return nil, fmt.Errorf("\"{\" expected, but got %v", c)
	}
	val, err = d.decodeObject()
	if err != nil {
		return nil, fmt.Errorf("invalid")
	}
	if c := d.skipSpaces(); d.pos < d.end {
		return nil, fmt.Errorf("invalid json: %v", c)
	}
	return val, nil
}

func debug(msg string, args ...interface{}) {
	s := fmt.Sprintf(msg, args...)
	fmt.Println(fmt.Sprintf("[DEBUG] %v", s))
}
