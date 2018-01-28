package fate

type Fate interface{

}

type stdFate struct{
	name Name
}

//Start the fate main entrance
func Start(last string) stdFate{
	return stdFate{
		name: NewName(last),
	}
}



func (s* stdFate)OneFirstName() *stdFate{
	return s
}