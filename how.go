// how is simple package for parsing arguments
// config values are sourced from a config file, the environment, and os arguments (ordered in increasing priority)
//
// the struct field tags to use for customizing behavior are:
// how-long
// how-short
// how-env
// how-help
package how

import (
	"errors"
	"os"
)

var (
	// ShowHelp is returned from any of the Parse* functions if "-h" or "--help" is one of the arguments
	// If returned, no other arguments/settings are parsed!
	// PrintHelp can then be used for convenience
	ShowHelp = errors.New("show help")
)

// Parse parses the environment and os.Args into config
func Parse(config interface{}) error {
	return parse(config, os.Args, nil)
}

// ParseArgs parses the environment and args into config
func ParseArgs(config interface{}, args []string) error {
	return parse(config, args, nil)
}

// ParseWithFile decodes the (key = value encoded) file at path into config, and then parses the environment and os.Args into config
func ParseWithFile(config interface{}, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return parse(config, os.Args, NewKeyValueDecoder(file))
}

// ParseWithDecoder decodes (using decoder) into config and then parses the environment and os.Args into config
func ParseWithDecoder(config interface{}, decoder Decoder) error {
	return parse(config, os.Args, decoder)
}

func parse(config interface{}, args []string, decoder Decoder) error {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return ShowHelp
		}
	}

	rv, err := getReflectStruct(config)
	if err != nil {
		return err
	}

	refMap := getReflectMap(rv)

	// file settings are overridden by the environment, which are overridden by flags
	if kvd, ok := decoder.(*KeyValueDecoder); ok {
		err = kvd.decode(config, refMap)
	} else {
		err = decoder.Decode(config)
	}

	if err != nil {
		return err
	}

	err = getEnv(refMap)
	if err != nil {
		return err
	}

	err = parseArgs(refMap, args)
	return err
}
