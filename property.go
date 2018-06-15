package fate

type Property struct {
	UseThreeTalent bool //三才
	UseFiveGrid    bool //五格
	UseFivePhase   bool //字符五行
	UseZodiac      bool //生肖喜忌
	UseLikeUseGod  bool //喜用神
}

var property *Property

func init() {
	property = DefaultProperty()
}

func DefaultProperty() *Property {
	property := &Property{
		UseThreeTalent: false,
		UseFiveGrid:    false,
		UseFivePhase:   false,
		UseZodiac:      false,
	}
	return property
}

func SetProperty(p *Property) {
	property = p
}

func GetProperty() *Property {
	return property
}
