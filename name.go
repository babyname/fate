package fate

import (
	"strconv"
	"strings"

	"github.com/godcong/chronos"
	"github.com/godcong/yi"
)

//Name 姓名
type Name struct {
	FirstName         []*Character //名
	LastName          []*Character //姓
	born              *chronos.Calendar
	baZi              *BaZi
	nameScienceStroke *NameStroke //康熙（古典）
	nameStroke        *NameStroke //简体
	zodiac            *Zodiac
	zodiacPoint       int
}

// String ...
func (n Name) String() string {
	var s string
	for _, l := range n.LastName {
		s += l.Ch
	}
	for _, f := range n.FirstName {
		s += f.Ch
	}
	return s
}

// Strokes ...
func (n Name) Strokes() string {
	var s []string
	for _, l := range n.LastName {
		s = append(s, strconv.Itoa(l.ScienceStroke))
	}

	for _, f := range n.FirstName {
		s = append(s, strconv.Itoa(f.ScienceStroke))
	}
	return strings.Join(s, ",")
}

// PinYin ...
func (n Name) PinYin() string {
	var s string
	for _, l := range n.LastName {
		s += "[" + strings.Join(l.PinYin, ",") + "]"
	}

	for _, f := range n.FirstName {
		s += "[" + strings.Join(f.PinYin, ",") + "]"
	}
	return s
}

// WuXing ...
func (n Name) WuXing() string {
	var s string
	for _, l := range n.LastName {
		s += l.WuXing
	}
	for _, f := range n.FirstName {
		s += f.WuXing
	}
	return s
}

// XiYongShen ...
func (n Name) XiYongShen() string {
	return n.baZi.XiYongShen()
}

//GetName
func (impl *fateImpl) createName(first []*Character, nk NameStroke) *Name {
	if len(first) > 2 || len(first) < 1 || len(impl.last) > 2 || len(impl.last) < 1 {
		panic("input error")
	}

	var first1, first2, last1, last2 *Character
	first1 = first[0]

	if len(first) == 2 {

		first2 = first[1]

	} else {
		first2 = nil
	}

	last1 = impl.lastChar[0]

	if len(impl.lastChar) == 2 {
		last2 = impl.lastChar[1]
	} else {
		last2 = nil
	}

	name := Name{
		FirstName: first,
		LastName:  impl.lastChar,
	}

	if first2 == nil {
		if last2 == nil {
			name.nameScienceStroke = GetNameStroke(last1.getStrokeScience(), 0, first1.getStrokeScience(), 0)
			name.nameStroke = NewNameStroke(last1.getStroke(), 0, first1.getStroke(), 0)
		} else {
			name.nameScienceStroke = GetNameStroke(last1.getStrokeScience(), last2.getStrokeScience(), first1.getStrokeScience(), 0)
			name.nameStroke = NewNameStroke(last1.getStroke(), last2.getStroke(), first1.getStroke(), 0)
		}
	} else {
		if last2 == nil {
			name.nameScienceStroke = GetNameStroke(last1.getStrokeScience(), 0, first1.getStrokeScience(), first2.getStrokeScience())
			name.nameStroke = NewNameStroke(last1.getStroke(), 0, first1.getStroke(), first2.getStroke())
		} else {
			name.nameScienceStroke = GetNameStroke(last1.getStrokeScience(), last2.getStrokeScience(), first1.getStrokeScience(), first2.getStrokeScience())
			name.nameStroke = NewNameStroke(last1.getStroke(), last2.getStroke(), first1.getStroke(), first2.getStroke())
		}
	}

	return &name
}

func (n *Name) IsLucky(sex yi.Sex, filter_gua bool, filter_hard bool) bool {
	if !n.nameScienceStroke.IsWuGeLucky(sex, filter_hard) {
		return false
	}

	if filter_gua {
		if filter_hard {
			if !n.nameScienceStroke.IsGuaLucky(sex, true) {
				return false
			} else if !n.nameStroke.IsGuaLucky(sex, true) {
				return false
			}
		} else {
			if !n.nameScienceStroke.IsGuaLucky(sex, false) {
				return false
			} else if !n.nameStroke.IsGuaLucky(sex, false) {
				return false
			}
		}
	}

	return true
}

// BaZi ...
func (n Name) BaZi() string {
	return n.baZi.String()
}
