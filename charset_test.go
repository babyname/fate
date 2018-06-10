package fate_test

import (
	"log"
	"testing"
)

type value struct {
	s1 string
}

func testStruct() value {
	v := value{
		s1: "value",
	}
	log.Printf("%p  %p", &v, &v.s1)
	return v
}

func TestInitAll(t *testing.T) {
	v1 := testStruct()
	log.Printf("%p %p", &v1, &v1.s1)
}
