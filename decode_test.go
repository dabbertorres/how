package how

import (
	"bytes"
	"testing"
)

func TestKeyValueDecoder_Decode(t *testing.T) {
	type testCfg struct {
		Foo   string  `how-long:"foo" how-help:"Foo does something"`
		Bar   bool    `how-help:"Bar signals something"`
		Baz   int     `how-long:"baz" how-help:"Baz is a signed number"`
		Qux   uint    `how-long:"qux" how-help:"qux is an unsigned number"`
		Xyzzy float64 `how-long:"xyzzy" how-help:"xyzzy is a 64-bit floating point number"`
	}

	const input = `# Foo does something
foo = foo value
# Bar signals something
Bar = true
# Baz is a signed number
baz = -5
# qux is an unsigned number
qux = 33
# xyzzy is a 64-bit floating point number
xyzzy = -1.5
`

	buf := bytes.NewBufferString(input)

	cfg := &testCfg{}
	err := NewKeyValueDecoder(buf).Decode(cfg)
	if err != nil {
		t.Error(err)
	} else if cfg.Foo != "foo value" {
		t.Errorf("cfg.Foo: expected 'foo value', got '%s'", cfg.Foo)
	} else if !cfg.Bar {
		t.Errorf("cfg.Bar: expected true, got false")
	} else if cfg.Baz != -5 {
		t.Errorf("cfg.Baz: expected -5, got %d", cfg.Baz)
	} else if cfg.Qux != 33 {
		t.Errorf("cfg.Qux: expected 33, got %d", cfg.Qux)
	} else if cfg.Xyzzy != -1.5 {
		t.Errorf("cfg.Xyzzy: expected -1.5, got %f", cfg.Xyzzy)
	}
}
