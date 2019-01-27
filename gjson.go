package gjson

import (
	"fmt"
	"strconv"
)

type decoder struct {
	pos  int
	end  int
	data []byte
}

func (d *decoder) skipSpaces() byte {
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

func (d *decoder) readString() (string, error) {
	d.pos++
	start := d.pos

	for {
		if d.pos >= d.end {
			return "", d.error("unexpected EOF")
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

func (d *decoder) readNumber() (n float64, err error) {
	start := d.pos
	c := d.data[d.pos]
	if c == '0' {
		c = d.readNext()
	} else {
		// When comes here, c must be in 1 to 9
		for ; '0' <= c && c <= '9'; c = d.readNext() {
			// e.g. when number is 12345,
			// 1 -> 10 + 2 -> 120 + 3 -> 1230 + 4 -> 12340 + 5 => 12345
			n = n*10 + float64(c-'0') // c-'0': cast
		}
	}

	if c == '.' {
		d.pos++
		if c = d.data[d.pos]; c < '0' || '9' < c {
			return 0, d.error("number is required after decimal point")
		}
		for c = d.readNext(); '0' <= c && c <= '9'; {
			c = d.readNext()
		}
	}

	tmpn := string(d.data[start:d.pos])
	n, _ = strconv.ParseFloat(tmpn, 64)
	return
}

func (d *decoder) readArray() (interface{}, error) {
	d.pos++
	arr := make([]interface{}, 0)
	if c := d.skipSpaces(); c == ']' {
		d.pos++
		return arr, nil
	}

	for {
		v, err := d.readAny()
		if err != nil {
			return arr, d.error(fmt.Sprintf("reading value failed: %s", err))
		}
		arr = append(arr, v)

		if c := d.skipSpaces(); c == ',' {
			d.pos++
			continue
		} else if c == ']' {
			d.pos++
			return arr, nil
		} else {
			return arr, d.error(`"," or "]" is expected`)
		}
	}
}

func (d *decoder) readAny() (interface{}, error) {
	switch c := d.skipSpaces(); c {
	case '"':
		return d.readString()

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return d.readNumber()
	case '-':
		d.pos++
		c := d.data[d.pos]
		if c < '0' || '9' < c {
			return nil, d.error("invalid in negatice number")
		}
		n, err := d.readNumber()
		return -n, err

	case 't':
		d.pos++
		if d.end-d.pos < 3 { // avoid index out of bounds
			return nil, d.error("Unexpected EOF")
		}
		if d.data[d.pos] == 'r' && d.data[d.pos+1] == 'u' && d.data[d.pos+2] == 'e' {
			d.pos += 3
			return true, nil
		}
		return nil, d.error(fmt.Sprintf(`"true" is expected but got "%s" next to "t"`, string(d.data[d.pos])))
	case 'f':
		d.pos++
		if d.end-d.pos < 4 { // avoid index out of bounds
			return nil, d.error("Unexpected EOF")
		}
		if d.data[d.pos] == 'a' && d.data[d.pos+1] == 'l' && d.data[d.pos+2] == 's' && d.data[d.pos+3] == 'e' {
			d.pos += 4
			return false, nil
		}
		return nil, d.error(fmt.Sprintf(`"false" is expected but got "%s" next to "f"`, string(d.data[d.pos])))

	case 'n':
		d.pos++
		if d.end-d.pos < 3 { // avoid index out of bounds
			return nil, d.error("Unexpected EOF")
		}
		if d.data[d.pos] == 'u' && d.data[d.pos+1] == 'l' && d.data[d.pos+2] == 'l' {
			d.pos += 3
			return nil, nil
		}
		return nil, d.error(fmt.Sprintf(`"null" is expected but got "%s" next to "n"`, string(d.data[d.pos])))

	case '[':
		return d.readArray()

	case '{':
		return d.readObject()
	default:
		return nil, d.error("value is invalid")
	}
}

func (d *decoder) readNext() byte {
	d.pos++
	return d.data[d.pos]
}

func (d *decoder) readObject() (obj map[string]interface{}, err error) {
	d.pos++

	var c byte
	var val interface{}
	var key string
	obj = make(map[string]interface{})

	// object ends
	if c = d.skipSpaces(); c == '}' {
		d.pos++
		return obj, nil
	}

	for {
		if c = d.skipSpaces(); c != '"' {
			err = d.error(fmt.Sprintf("key must be string: %s", string(c)))
			break
		}

		if key, err = d.readString(); err != nil {
			err = d.error("key is invalid")
			break
		}

		if c = d.skipSpaces(); c != ':' {
			err = d.error("after object key")
			break
		}
		d.pos++

		if val, err = d.readAny(); err != nil {
			break
		}

		obj[key] = val

		if c = d.skipSpaces(); c == '}' {
			d.pos++
			break
		} else if c == ',' {
			d.pos++
		} else {
			err = d.error("invalid object")
			break
		}
	}
	return
}

func (d *decoder) error(msg string) error {
	return fmt.Errorf("invalid json: %s at %d", msg, d.pos+1)
}

// Decode returns parsed json object.
// arg must be started "{" and valid as JSON.
func Decode(data []byte) (val map[string]interface{}, err error) {
	d := &decoder{
		pos:  0,
		end:  len(data),
		data: data,
	}

	if c := d.skipSpaces(); c != '{' {
		return nil, d.error(fmt.Sprintf(`"{" expected, but got %v`, c))
	}
	val, err = d.readObject()
	if err != nil {
		return nil, d.error(fmt.Sprintf("invalid: %s", err))
	}
	if c := d.skipSpaces(); d.pos < d.end {
		return nil, d.error(fmt.Sprintf("invalid json: %v", c))
	}
	return val, nil
}
