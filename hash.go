package fate

import (
	"github.com/cespare/xxhash/v2"
)

func hash(name string) uint64 {
	return xxhash.Sum64([]byte(name))
}
