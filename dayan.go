package fate

var daYanList = []*DaYan{
	{
		//Luck:    Poin0,
		SkyNine: "太极之数",
		Comment: "太极之数，万物开泰，生发无穷，利禄亨通。",
	},
}

//DaYan
type DaYan struct {
	Luck    Luck   `json:"luck"`     //吉凶
	SkyNine string `json:"sky_nine"` //天九(天九地十取天九)
	Comment string `json:"comment"`
}
