package fate

import "strings"

type Name struct {
	FirstName []string
	LastName  []string
}

func NewName(last string) Name {
	name := Name{}
	if len(last)>1{
		name.LastName = strings.Split(last,"")
		return name
	}
	name.LastName[0] = last
	return name
}