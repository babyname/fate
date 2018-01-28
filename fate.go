package fate

type Fate interface{

}


type defaultFate struct{
	data NameData
}

//Start the fate main entrance
func Start(last string) Fate{
	return defaultFate{
		data: NameData{LastName:last},
	}
}