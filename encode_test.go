package how

import (
	"bytes"
	"testing"
)

func TestKeyValueEncoder_Encode(t *testing.T) {
	type testCfg struct {
		Foo   string  `how-long:"foo" how-help:"Foo does something"`
		Bar   bool    `how-help:"Bar signals something"`
		Baz   int     `how-long:"baz" how-help:"Baz is a signed number"`
		Qux   uint    `how-long:"qux" how-help:"qux is an unsigned number"`
		Xyzzy float64 `how-long:"xyzzy" how-help:"xyzzy is a 64-bit floating point number"`
	}

	cfg := &testCfg{
		Foo:   "foo value",
		Bar:   true,
		Baz:   -5,
		Qux:   33,
		Xyzzy: -1.5,
	}

	const expected = `# Foo does something
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

	buf := bytes.NewBuffer(nil)

	err := NewKeyValueEncoder(buf).Encode(cfg)
	if err != nil {
		t.Error(err)
	}

	if buf.String() != expected {
		t.Errorf("got:\n%s\n\nexpected:\n%s", buf.String(), expected)
	}
}
