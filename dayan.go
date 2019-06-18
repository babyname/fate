package fate

//DaYan
type DaYan struct {
	Luck    Luck   `json:"luck"`     //吉凶
	SkyNine string `json:"sky_nine"` //天九(天九地十取天九)
	Comment string `json:"comment"`
}
