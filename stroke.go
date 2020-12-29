package fate

//Strokes 推荐笔画数
type Strokes struct {
	SimplifiedChinese  int //简体中文
	TraditionalChinese int //繁体中文
	KangxiDictionary   int //康熙字典
	GodBook            int //神册
}

//FindCharacterStrokes 通过文字查询推荐笔画数
func FindCharacterStrokes(char string) int {
	//TODO:find
	s := &Strokes{}
	if s.GodBook > 0 {
		return s.GodBook
	}
	if s.KangxiDictionary > 0 {
		return s.KangxiDictionary
	}
	if s.TraditionalChinese > 0 {
		return s.TraditionalChinese
	}
	if s.SimplifiedChinese > 0 {
		return s.SimplifiedChinese
	}
	return 0
}
