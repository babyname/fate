package fate

import "strings"

const (
	ZodiacPig = iota
)

var ZodiacList = []Zodiac{}

//Zodiac 生肖
type Zodiac struct {
	Xi        string //喜
	XiRadical string
	Ji        string //忌
	JiRadical string
}

func ZodiacPoint(z *Zodiac, character *Character) int {
	return z.Point(character)
}

func (z *Zodiac) zodiacJi(character *Character) int {
	if strings.IndexRune(z.Ji, []rune(character.Ch)[0]) != -1 {
		return -3
	}
	return 0
}

//Point 喜忌对冲，理论上喜忌都有的话，最好不要选给1，忌给0，喜给5，都没有给3
func (z *Zodiac) Point(character *Character) int {
	dp := 3
	dp += z.zodiacJi(character)
	dp += z.zodiacXi(character)
	return dp
}

func (z *Zodiac) zodiacXi(character *Character) int {
	if strings.IndexRune(z.Xi, []rune(character.Ch)[0]) != -1 {
		return 2
	}
	return 0
}
