package how

import (
	"reflect"
	"testing"
)

type argsTestConfig struct {
	Foo   string `how-long:"foo"`
	Bar   string `how-long:"bar"`
	Baz   string `how-long:"baz"`
	Qux   string `how-long:"qux"`
	Xyzzy bool   `how-long:"xyzzy"`

	FooB   bool `how-short:"f"`
	BarB   bool `how-short:"b"`
	QuxB   bool `how-short:"q"`
	XyzzyB bool `how-short:"x"`

	Other bool `how-long:"other" how-short:"o"`
}

func Test_parseArgs(t *testing.T) {
	cfg := &argsTestConfig{}
	rv, err := getReflectStruct(cfg)
	if err != nil {
		t.Skip("getReflectStruct() failed:", err)
	}

	values := getReflectMap(rv)

	err = parseArgs(values, []string{"--foo=value0",
		"-f",
		"--bar=", "value1",
		"--xyzzy",
		"--baz", "=value2",
		"-b", "true",
		"-qx",
		"--qux", "=", "value3",
		"other",
	})
	if err != nil {
		t.Error(err)
	}

	checkStr := func(variable, str, expect string) {
		t.Helper()
		if str != expect {
			t.Errorf("%s: expected '%s', got '%s'", variable, expect, str)
		}
	}

	checkBool := func(variable string, b bool) {
		t.Helper()
		if !b {
			t.Errorf("%s: expected to be true", variable)
		}
	}

	checkStr("cfg.Foo", cfg.Foo, "value0")
	checkBool("cfg.FooB", cfg.FooB)
	checkStr("cfg.Bar", cfg.Bar, "value1")
	checkBool("cfg.Xyzzy", cfg.Xyzzy)
	checkStr("cfg.Baz", cfg.Baz, "value2")
	checkBool("cfg.BarB", cfg.BarB)
	checkBool("cfg.QuxB", cfg.QuxB)
	checkBool("cfg.XyzzyB", cfg.XyzzyB)
	checkStr("cfg.Qux", cfg.Qux, "value3")
	checkBool("cfg.Other", cfg.Other)
}

func Test_parseLongFlag(t *testing.T) {
	args := []string{"--foo=value",
		"--bar=", "value",
		"--baz", "=value",
		"--qux", "=", "value",
		"--xyzzy",
	}

	// dummy value
	values := reflectMap{}

	var (
		key  string
		val  string
		used int
		err  error
	)

	checks := func(printErr func(), expectedUsed int, expectedKey, expectedVal string) {
		t.Helper()
		if err != nil {
			printErr()
			t.Error(err)
			return
		}

		if used != expectedUsed {
			printErr()
			t.Errorf("used: expected %d, got %d", expectedUsed, used)
			return
		}

		if key != expectedKey {
			printErr()
			t.Errorf("key: expected '%s', got '%s'", expectedKey, key)
			return
		}

		if val != expectedVal {
			printErr()
			t.Errorf("val: expected '%s', got '%s'", expectedVal, val)
			return
		}
	}

	key, val, _, used, err = parseLongFlag(values, args)
	checks(func() { t.Errorf("'%s' failed:", args[0]) }, 0, "foo", "value")

	key, val, _, used, err = parseLongFlag(values, args[1:])
	checks(func() { t.Errorf("'%s %s' failed:", args[1], args[2]) }, 1, "bar", "value")

	key, val, _, used, err = parseLongFlag(values, args[3:])
	checks(func() { t.Errorf("'%s %s' failed:", args[3], args[4]) }, 1, "baz", "value")

	key, val, _, used, err = parseLongFlag(values, args[5:])
	checks(func() { t.Errorf("'%s = %s' failed:", args[5], args[7]) }, 2, "qux", "value")

	key, val, _, used, err = parseLongFlag(values, args[8:])
	checks(func() { t.Errorf("'%s' failed:", args[8]) }, 0, "xyzzy", "")

	key, val, _, used, err = parseLongFlag(values, []string{"--foo="})
	if !IsMissingValueError(err) {
		t.Errorf("'--foo=' failed: expected missing value error, got: %v", err)
	}

	key, val, _, used, err = parseLongFlag(values, []string{"--foo", "="})
	if !IsMissingValueError(err) {
		t.Errorf("'--foo =' failed: expected missing value error, got: %v", err)
	}
}

func Test_parseShortFlag(t *testing.T) {
	args := []string{"-f",
		"-b", "value",
		"-fbqx",
		"-fz",
	}

	cfg := &argsTestConfig{}
	rv, err := getReflectStruct(cfg)
	if err != nil {
		t.Skip("getReflectStruct() failed:", err)
	}

	values := getReflectMap(rv)

	var (
		key  string
		val  string
		vm   map[string]reflect.Value
		used int
	)

	checks := func(printErr func(), expectedUsed int, expectedKey, expectedVal string) {
		t.Helper()
		if err != nil {
			printErr()
			t.Error(err)
			return
		}

		if used != expectedUsed {
			printErr()
			t.Errorf("used: expected %d, got %d", expectedUsed, used)
			return
		}

		if key != expectedKey {
			printErr()
			t.Errorf("key: expected '%s', got '%s'", expectedKey, key)
			return
		}

		if val != expectedVal {
			printErr()
			t.Errorf("val: expected '%s', got '%s'", expectedVal, val)
			return
		}
	}

	key, val, _, used, err = parseShortFlag(values, args)
	checks(func() { t.Errorf("'%s' failed:", args[0]) }, 0, "f", "")

	key, val, _, used, err = parseShortFlag(values, args[1:])
	checks(func() { t.Errorf("'%s %s' failed:", args[1], args[2]) }, 1, "b", "value")

	key, val, vm, used, err = parseShortFlag(values, args[3:])
	if err != nil {
		t.Errorf("'%s' failed: %v", args[3], err)
	} else if vm != nil {
		t.Errorf("'%s' failed: expected nil valuesMap (multiple short flags)", args[3])
	} else if used != 0 {
		t.Errorf("'%s' failed: used: expected 0, got %d", args[3], used)
	} else if !cfg.FooB {
		t.Errorf("'%s' failed: expected cfg.FooB to be true", args[3])
	} else if !cfg.BarB {
		t.Errorf("'%s' failed: expected cfg.BarB to be true", args[3])
	} else if !cfg.QuxB {
		t.Errorf("'%s' failed: expected cfg.QuxB to be true", args[3])
	} else if !cfg.XyzzyB {
		t.Errorf("'%s' failed: expected cfg.XyzzyB to be true", args[3])
	}

	key, val, vm, used, err = parseShortFlag(values, args[4:])
	if !IsNotFlagError(err) {
		t.Errorf("'%s' failed: expected ErrNotFlag error, got: %v", args[4], err)
	}
}
