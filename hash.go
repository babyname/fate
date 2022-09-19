package fate

import (
	"github.com/cespare/xxhash"
)

func hash(name string) uint64 {
	return xxhash.Sum64([]byte(name))
}
