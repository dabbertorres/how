package how

import (
	"fmt"
	"io"
)

// NewKeyValueEncoder encodes structs formatted as
// # comment
// <key> = <value>
// comments are placed above key-value pairs if a field has a non-empty "how-help" tag
// the key is either the field's (non-empty) "how-long" tag, or the field name
func NewKeyValueEncoder(w io.Writer) *KeyValueEncoder {
	return &KeyValueEncoder{w: w}
}

type KeyValueEncoder struct {
	w io.Writer
}

func (e *KeyValueEncoder) Encode(v interface{}) error {
	rv, err := getReflectStruct(v)
	if err != nil {
		return err
	}
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		var (
			fv = rv.Field(i)
			ft = rt.Field(i)
		)

		if help := ft.Tag.Get("how-help"); help != "" {
			fmt.Fprintln(e.w, "#", help)
		}

		var key string
		if tag := ft.Tag.Get("how-long"); tag != "" {
			key = tag
		} else {
			key = ft.Name
		}

		fmt.Fprintf(e.w, "%s = %v\n", key, fv)
	}

	return nil
}
