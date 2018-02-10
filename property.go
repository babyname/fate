package fate

type Property interface {
	UseThreeTalent() bool
	UseFiveGrid() bool
	UseFivePhase() bool
	UseZodiac() bool
}

type stdProperty struct {
	useThreeTalent bool //三才
	useFiveGrid    bool //五格
	useFivePhase   bool //字符五行
	useZodiac      bool //生肖
}

var property Property

func init() {
	property = DefaultProperty()
}

func (s *stdProperty) UseThreeTalent() bool {
	return s.useThreeTalent
}

func (s *stdProperty) UseFiveGrid() bool {
	return s.useFiveGrid
}

func (s *stdProperty) UseFivePhase() bool {
	return s.useFivePhase
}

func (s *stdProperty) UseZodiac() bool {
	return s.useZodiac
}

func DefaultProperty() Property {
	property := &stdProperty{
		useThreeTalent: false,
		useFiveGrid:    false,
		useFivePhase:   false,
		useZodiac:      false,
	}
	return property
}

func SetProperty(p Property) {
	property = p
}

func GetProperty() Property {
	return property
}
