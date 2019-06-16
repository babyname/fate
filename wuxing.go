package fate

type Luck string

//WuXing 五行：five elements of metal,wood,water,fire and earth
type WuXing struct {
	WuXing  string `json:"wu_xing"`
	Luck    Luck   `json:"luck"`
	Comment string `json:"comment"`
}

//FindWuXing find a wuxing
func FindWuXing(wx string) *WuXing {
	return nil
}
