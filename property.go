package fate

type Property interface {

}

type stdProperty struct {
	useThreeTalent bool //三才
	useFiveGrid bool //五格
	useFivePhase bool //字符五行
	useZodiac bool //生肖
}

func loadProperty() Property {
	property := stdProperty{}
	return property
}