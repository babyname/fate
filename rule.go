package fate

type RuleIndex int

const (
	RuleIndexSupply  RuleIndex = iota //过滤补八字
	RuleIndexZodiac                   //过滤生肖
	RuleIndexBagua                    //过滤卦象
	RuleIndexDayan                    //过滤大衍
	RuleIndexRegular                  //常用字过滤
	RuleIndexMax
)

type Rule struct {
	sex   Sex
	rules [RuleIndexMax]bool
}

type RuleOption func(rule *Rule)

func (r *Rule) Set(idx RuleIndex, b bool) {
	if idx >= RuleIndexMax {
		return
	}
	r.rules[idx] = b
}

func (r Rule) Get(idx RuleIndex) bool {
	return r.rules[idx]
}

// RuleFilterBagua ...
// @Description: 过滤卦象
// @param bool
// @return func(rule *Rule)
func RuleFilterBagua(b bool) func(rule *Rule) {
	return func(rule *Rule) {
		rule.Set(RuleIndexBagua, b)
	}
}

// RuleFilterZodiac ...
// @Description: 生肖
// @param bool
// @return func(rule *Rule)
func RuleFilterZodiac(b bool) func(rule *Rule) {
	return func(rule *Rule) {
		rule.Set(RuleIndexZodiac, b)
	}
}

// RuleFilterDayan ...
// @Description: 过滤大衍
// @param bool
// @return func(rule *Rule)
func RuleFilterDayan(b bool) func(rule *Rule) {
	return func(rule *Rule) {
		rule.Set(RuleIndexDayan, b)
	}
}

// RuleFilterRegular ...
// @Description: 常用字过滤
// @param bool
// @return func(rule *Rule)
func RuleFilterRegular(b bool) func(rule *Rule) {
	return func(rule *Rule) {
		rule.Set(RuleIndexRegular, b)
	}
}

// RuleFilterSupply ...
// @Description: 过滤补八字
// @param bool
// @return func(rule *Rule)
func RuleFilterSupply(b bool) func(rule *Rule) {
	return func(rule *Rule) {
		rule.Set(RuleIndexSupply, b)
	}
}

// RuleSetSex ...
// @Description: set sex to rule
// @param Sex
// @return func(rule *Rule)
func RuleSetSex(s Sex) func(rule *Rule) {
	return func(rule *Rule) {
		rule.sex = s
	}
}

// RuleSetCopy ...
// @Description: set a rule copy to rule
// @param *Rule
// @return func(rule *Rule)
func RuleSetCopy(r *Rule) func(rule *Rule) {
	return func(rule *Rule) {
		*rule = *r
	}
}

// DefaultRule ...
// @Description:
// @param Sex
// @return *Rule
func DefaultRule() *Rule {
	return &Rule{
		sex: false,
		rules: [5]bool{
			true, true, true, true, true,
		},
	}
}
