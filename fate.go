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



func (s* stdFate)GenerateOne() *stdFate{
	//先套三才五格，再找字




	return s
}


