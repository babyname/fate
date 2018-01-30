package fate

type Property interface {
	UseThreeTalent() bool
	UseFiveGrid() bool
	UseFivePhase() bool
	UseZodiac() bool
}

type stdProperty struct {
	useThreeTalent bool //三才
	useFiveGrid bool //五格
	useFivePhase bool //字符五行
	useZodiac bool //生肖
}

func (s *stdProperty) UseThreeTalent() bool {
	panic("implement me")
}

func (s *stdProperty) UseFiveGrid() bool {
	panic("implement me")
}

func (s *stdProperty) UseFivePhase() bool {
	panic("implement me")
}

func (s *stdProperty) UseZodiac() bool {
	panic("implement me")
}

func loadProperty() Property {
	property := &stdProperty{}
	return property
}