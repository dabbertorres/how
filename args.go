package how

import (
	"reflect"
	"strings"
)

func parseArgs(values reflectMap, args []string) (err error) {
	for i := 0; i < len(args); i++ {
		var (
			key       string
			val       string
			valuesMap map[string]reflect.Value
			used      int
		)

		if isLongFlag(args[i]) {
			key, val, valuesMap, used, err = parseLongFlag(values, args[i:])
		} else if isShortFlag(args[i]) {
			key, val, valuesMap, used, err = parseShortFlag(values, args[i:])
		} else {
			// pass without processing as a long key without a value
			valuesMap = values.long
			key = args[i]
		}

		if err != nil {
			return
		}

		i += used

		// flags already added
		if valuesMap == nil {
			continue
		}

		if rv, ok := valuesMap[key]; ok {
			if val == "" && rv.Kind() != reflect.Bool {
				return errMissingValue(key)
			}

			err := setValue(rv, val)
			if err != nil {
				return err
			}
		} else {
			return errNotFlag(key)
		}
	}

	return nil
}

func parseLongFlag(values reflectMap, args []string) (key, val string, valuesMap map[string]reflect.Value, used int, err error) {
	valuesMap = values.long

	key = strings.TrimPrefix(args[0], "--")
	// long form (key = value) could be 1 - 3 'args'

	if strings.Contains(key, "=") {
		// could be "key=value"
		kv := strings.SplitN(key, "=", 2)
		key = kv[0]

		if len(kv) != 2 || kv[1] == "" {
			// ...or "key= value"

			if len(args) > 1 {
				val = args[1]
				used++
			} else {
				err = errMissingValue(key)
				return
			}
		} else {
			val = kv[1]
		}
	} else if len(args) > 1 && strings.HasPrefix(args[1], "=") {
		// looks like "key =value" or "key = value"

		val = strings.TrimPrefix(args[1], "=")
		if val == "" {
			// it's "key = value"

			if len(args) > 2 {
				val = args[2]
				used += 2
			} else {
				err = errMissingValue(key)
				return
			}
		} else {
			// it's "key =value"
			used++
		}
	} else {
		// just an on/off flag, don't need to do anything
		// (or the next arg is a flag, and this flag's value is missing, which we'll find out later)
	}

	return
}

func parseShortFlag(values reflectMap, args []string) (key, val string, valuesMap map[string]reflect.Value, used int, err error) {
	key = strings.TrimPrefix(args[0], "-")

	// is it a group of on/off flags?
	if len(key) > 1 {
		// going to handle setting the flags here, since we have multiple
		for start, end := 0, 1; start < len(key); start, end = start+1, end+1 {
			flag := key[start:end]
			if rv, ok := values.short[flag]; ok {
				setValue(rv, "")
			} else {
				err = errNotFlag(flag)
				return
			}
		}
	} else {
		valuesMap = values.short

		// if the next arg isn't a flag, then it's the value for this flag
		if len(args) > 1 && !(isLongFlag(args[1]) || isShortFlag(args[1])) {
			val = args[1]
			used++
		} else {
			// on/off flag, don't need to do anything
			// (or the next arg is a flag, and this flag's value is missing, which we'll find out later)
		}
	}

	return
}

func isLongFlag(arg string) bool {
	return strings.HasPrefix(arg, "--")
}

func isShortFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}
