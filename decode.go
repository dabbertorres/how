package how

import (
	"bufio"
	"io"
	"strings"
)

type Decoder interface {
	Decode(interface{}) error
}

// KeyValueDecoder decodes structs formatted as
// # comment
// <key> = <value>
// comments are ignored
// the key is either the field's (non-empty) "how-long" tag, or the field name
type KeyValueDecoder struct {
	r io.Reader
}

func NewKeyValueDecoder(r io.Reader) *KeyValueDecoder {
	return &KeyValueDecoder{r: r}
}

func (d *KeyValueDecoder) Decode(v interface{}) error {
	rv, err := getReflectStruct(v)
	if err != nil {
		return err
	}

	return d.decode(v, getReflectMap(rv))
}

func (d *KeyValueDecoder) decode(v interface{}, values reflectMap) error {
	scan := bufio.NewScanner(d.r)

	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())

		// skip empty lines and comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			// TODO continue or exit?
			continue
		}

		var (
			key = strings.TrimSpace(kv[0])
			val = strings.TrimSpace(kv[1])
		)

		if rv, ok := values.long[key]; ok {
			err := setValue(rv, val)
			if err != nil {
				// TODO continue or exit?
				continue
			}
		}
	}

	return scan.Err()
}
