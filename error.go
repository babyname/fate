package fate

import (
	"fmt"
	"strings"
)

// Wrap ...
func Wrap(err error, msg ...string) error {
	if err != nil {
		m := strings.Join(msg, " ")
		return fmt.Errorf("%s:%w", m, err)
	}
	return nil
}
