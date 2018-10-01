package how

import (
	"os"
)

func getEnv(values reflectMap) error {
	for key, rv := range values.env {
		if val, ok := os.LookupEnv(key); ok {
			err := setValue(rv, val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
