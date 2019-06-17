package fate

type Luck int

var luckPoint = []string{"大凶", "凶", "凶多于吉", "吉凶参半", "吉多于凶", "吉", "大吉"}

func (l *Luck) Point() int {
	return int(*l) + 1
}

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
