package config

import (
	"fmt"
	"reflect"
)

func validVolumeModes(v interface{}, param string) error {
	validNames := map[string]bool{"ReadWriteOnce": true, "ReadOnlyMany": true, "ReadWriteMany": true}
	st := reflect.ValueOf(v)

	if st.Kind() != reflect.String {
		return fmt.Errorf(
			"invalid volume mode: %v. Valid modes: %s",
			st,
			"ReadWriteOnce,ReadOnlyMany,ReadWriteMany",
		)
	}

	if !validNames[st.String()] {
		return fmt.Errorf("invalid volume mode: %v", st)
	}
	return nil
}
