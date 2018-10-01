package how

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

// PrintHelp calls PrintHelpTo with os.Stderr
func PrintHelp(header string, config interface{}) error {
	return PrintHelpTo(os.Stderr, header, config)
}

// PrintHelp prints header, followed by the help messages for each field in config
// If a field does not have a 'how-help' tag, one is generated for the field
// A default value for each field is printed using the values of each field in config
func PrintHelpTo(w io.Writer, header string, config interface{}) error {
	if header != "" {
		fmt.Fprintln(w, header)
	}

	rv, err := getReflectStruct(config)
	if err != nil {
		return err
	}
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		var (
			fv = rv.Field(i)
			ft = rt.Field(i)
		)

		if tag := ft.Tag.Get("how-long"); tag != "" {
			fmt.Fprintf(w, "--%s", tag)
		} else {
			fmt.Fprintf(w, "--%s", ft.Name)
		}

		if tag := ft.Tag.Get("how-short"); tag != "" {
			fmt.Fprintf(w, ", -%s", tag)
		}

		if tag := ft.Tag.Get("how-env"); tag != "" {
			if runtime.GOOS == "windows" {
				fmt.Fprintf(w, ", %%%s%%", tag)
			} else {
				fmt.Fprintf(w, ", $%s", tag)
			}
		}

		if tag := ft.Tag.Get("how-help"); tag != "" {
			fmt.Fprintf(w, "\t- %s\n", tag)
		} else {
			fmt.Fprintf(w, "\t- %s\n", ft.Type.String())
		}

		fmt.Fprintf(w, "\t\tdefault: %v\n\n", fv)
	}

	return nil
}
