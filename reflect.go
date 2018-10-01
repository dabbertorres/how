package how

import (
	"reflect"
	"strconv"
)

func setValue(rv reflect.Value, val string) (err error) {
	switch rv.Kind() {
	case reflect.Bool:
		bv := true
		if val != "" {
			bv, err = strconv.ParseBool(val)
			if err != nil {
				return
			}
		}
		rv.SetBool(bv)

	case reflect.Int:
		err = setValueInt(rv, val, 0)
	case reflect.Int8:
		err = setValueInt(rv, val, 8)
	case reflect.Int16:
		err = setValueInt(rv, val, 16)
	case reflect.Int32:
		err = setValueInt(rv, val, 32)
	case reflect.Int64:
		err = setValueInt(rv, val, 64)

	case reflect.Uint:
		err = setValueUint(rv, val, 0)
	case reflect.Uint8:
		err = setValueUint(rv, val, 8)
	case reflect.Uint16:
		err = setValueUint(rv, val, 16)
	case reflect.Uint32:
		err = setValueUint(rv, val, 32)
	case reflect.Uint64:
		err = setValueUint(rv, val, 64)

	case reflect.Float32:
		err = setValueFloat(rv, val, 32)
	case reflect.Float64:
		err = setValueFloat(rv, val, 64)

	case reflect.String:
		rv.SetString(val)

	case reflect.Slice:
		// TODO

	default:
		err = errUnsupportedType(rv.Type().Name())
	}

	return
}

func setValueInt(rv reflect.Value, val string, bitSize int) error {
	iv, err := strconv.ParseInt(val, 0, bitSize)
	if err != nil {
		return err
	}
	rv.SetInt(iv)
	return nil
}

func setValueUint(rv reflect.Value, val string, bitSize int) error {
	uv, err := strconv.ParseUint(val, 0, bitSize)
	if err != nil {
		return err
	}
	rv.SetUint(uv)
	return nil
}

func setValueFloat(rv reflect.Value, val string, bitSize int) error {
	fv, err := strconv.ParseFloat(val, bitSize)
	if err != nil {
		return err
	}
	rv.SetFloat(fv)
	return nil
}

func getReflectStruct(v interface{}) (reflect.Value, error) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || (rv.Kind() != reflect.Ptr && rv.Kind() != reflect.Interface) {
		return rv, ErrInvalidValue
	}
	rv = rv.Elem()

	if rv.Kind() != reflect.Struct {
		return rv, ErrNotStruct
	}

	return rv, nil
}

type reflectMap struct {
	long  map[string]reflect.Value
	short map[string]reflect.Value
	env   map[string]reflect.Value
}

// build a map of field identifiers to the corresponding (reflective) values
func getReflectMap(v reflect.Value) (rm reflectMap) {
	rm.long = make(map[string]reflect.Value)
	rm.short = make(map[string]reflect.Value)
	rm.env = make(map[string]reflect.Value)

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		var (
			fv = v.Field(i)
			ft = t.Field(i)
		)

		if tag := ft.Tag.Get("how-long"); tag != "" {
			rm.long[tag] = fv
		} else {
			rm.long[ft.Name] = fv
		}

		if short := ft.Tag.Get("how-short"); short != "" {
			rm.short[short[:1]] = fv
		}

		if env := ft.Tag.Get("how-env"); env != "" {
			rm.env[env] = fv
		}
	}

	return
}
