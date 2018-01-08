package fate

type NameFate struct {
	FirstName string
	LastName  string
	FiveGrid
	ThreeTalent
}

//输入姓
func InsertLastName(name string) NameFate {
	return NameFate{}
}

//输入名
func InsertFirstName(name string) NameFate {
	return NameFate{}
}
