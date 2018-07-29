package mongo

//XiYongShen 喜用神
type XiYongShen struct {
	XiShen   string `json:"xi_shen"`   //喜神
	YongShen string `json:"yong_shen"` //用神
	Jin      string `json:"jin"`       //五行:金
	Mu       string `json:"mu"`        //五行:木
	Shui     string `json:"shui"`      //五行:水
	Huo      string `json:"huo"`       //五行:火
	Tu       string `json:"tu"`        //五行:土
}
