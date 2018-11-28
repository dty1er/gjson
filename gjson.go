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

func (d *Decoder) readString() (string, error) {
	d.pos++
	start := d.pos

	for {
		if d.pos >= d.end {
			return "", fmt.Errorf("unexpected EOF")
		}

		c := d.data[d.pos]
		if c == '"' {
			s := string(d.data[start:d.pos])
			d.pos++
			return s, nil
		}
		d.pos++
	}
}

func (d *Decoder) decodeObject() (obj map[string]interface{}, err error) {
	d.pos++

	var c byte
	var key string
	var val string
	obj = make(map[string]interface{})

	// object ends
	if c = d.skipSpaces(); c == '}' {
		d.pos++
		return obj, nil
	}

	for {
		if c = d.skipSpaces(); c != '"' {
			err = fmt.Errorf("key must be string")
			break
		}

		if key, err = d.readString(); err != nil {
			err = fmt.Errorf("key is invalid")
			break
		}

		if c = d.skipSpaces(); c != ':' {
			err = fmt.Errorf("after object key")
			break
		}
		d.pos++

		if c = d.skipSpaces(); c != '"' {
			err = fmt.Errorf("value must be string")
			break
		}

		if val, err = d.readString(); err != nil {
			err = fmt.Errorf("value is invalid")
			break
		}

		obj[key] = val

		if c = d.skipSpaces(); c == '}' {
			d.pos++
			break
		}
	}
	return
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
