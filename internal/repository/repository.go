package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/ent/poemchar"
)

type Repository struct {
	*ent.Client
	cache *ModelCache
}

func (m *Repository) Initialize(ctx context.Context) error {
	return m.Schema.Create(ctx)
}

func makeFallbackChar(ch string) *ent.Character {
	stroke, wuxing := surnameStrokeAndWuxing(ch)
	return &ent.Character{
		Char:              ch,
		IsSimplified:      true,
		Pinyin:            []string{},
		SimplifiedStroke:  stroke,
		TraditionalStroke: stroke,
		KangxiStroke:      stroke,
		ScienceStroke:     stroke,
		WuXing:            wuxing,
		Regular:           true,
		Nameable:          false,
		Meaning:           "姓氏用字",
		Source:            "fallback",
	}
}

var surnameData = map[string][2]string{
	"王": {"4", "土"}, "李": {"7", "木"}, "张": {"11", "火"}, "刘": {"15", "火"},
	"陈": {"16", "火"}, "杨": {"13", "木"}, "赵": {"14", "火"}, "黄": {"12", "土"},
	"周": {"8", "金"}, "吴": {"7", "木"}, "徐": {"10", "金"}, "孙": {"10", "金"},
	"胡": {"11", "土"}, "朱": {"6", "木"}, "高": {"10", "木"}, "林": {"8", "木"},
	"何": {"7", "木"}, "郭": {"15", "木"}, "马": {"10", "水"}, "罗": {"20", "火"},
	"梁": {"11", "火"}, "宋": {"7", "金"}, "郑": {"19", "火"}, "谢": {"17", "金"},
	"韩": {"17", "水"}, "唐": {"10", "火"}, "冯": {"12", "水"}, "于": {"3", "土"},
	"董": {"15", "木"}, "萧": {"18", "木"}, "程": {"12", "火"}, "曹": {"11", "金"},
	"袁": {"10", "土"}, "邓": {"19", "火"}, "许": {"11", "木"}, "傅": {"12", "水"},
	"沈": {"8", "水"}, "曾": {"12", "金"}, "彭": {"12", "水"}, "吕": {"7", "火"},
	"苏": {"22", "木"}, "卢": {"16", "火"}, "蒋": {"17", "木"}, "蔡": {"17", "木"},
	"贾": {"13", "水"}, "丁": {"2", "火"}, "魏": {"18", "木"}, "薛": {"19", "木"},
	"叶": {"15", "土"}, "阎": {"16", "木"}, "余": {"7", "土"}, "潘": {"16", "水"},
	"杜": {"7", "木"}, "戴": {"18", "火"}, "夏": {"10", "火"}, "钟": {"17", "金"},
	"汪": {"8", "水"}, "田": {"5", "火"}, "任": {"6", "金"}, "姜": {"9", "木"},
	"范": {"15", "水"}, "方": {"4", "水"}, "石": {"5", "金"}, "姚": {"9", "土"},
	"谭": {"19", "火"}, "廖": {"14", "火"}, "邹": {"17", "金"}, "熊": {"14", "水"},
	"金": {"8", "金"}, "陆": {"16", "火"}, "郝": {"14", "金"}, "孔": {"4", "木"},
	"白": {"5", "水"}, "崔": {"11", "木"}, "康": {"11", "木"}, "毛": {"4", "水"},
	"邱": {"12", "木"}, "秦": {"10", "火"}, "江": {"7", "水"}, "史": {"5", "金"},
	"顾": {"21", "木"}, "侯": {"9", "水"}, "邵": {"12", "金"}, "孟": {"8", "水"},
	"龙": {"16", "火"}, "万": {"15", "水"}, "段": {"9", "火"}, "雷": {"13", "水"},
	"钱": {"16", "金"}, "汤": {"13", "水"}, "尹": {"4", "土"}, "黎": {"15", "火"},
	"易": {"8", "火"}, "常": {"11", "金"}, "武": {"8", "水"}, "乔": {"12", "木"},
	"贺": {"12", "水"}, "赖": {"16", "火"}, "龚": {"22", "木"}, "文": {"4", "水"},
	"西": {"6", "金"}, "门": {"8", "水"}, "东": {"8", "木"}, "南": {"9", "火"},
	"北": {"5", "水"}, "司": {"5", "金"}, "上": {"3", "金"},
	"长": {"8", "火"}, "欧阳": {"15", "土"}, "司马": {"15", "金"},
	"诸葛": {"16", "木"}, "上官": {"3", "金"}, "慕容": {"15", "水"},
}

func surnameStrokeAndWuxing(ch string) (int, string) {
	if data, ok := surnameData[ch]; ok {
		stroke, _ := strconv.Atoi(data[0])
		if stroke > 0 {
			return stroke, data[1]
		}
	}
	return 1, "金"
}

func (m *Repository) QueryLastName(ctx context.Context, last [2]string) (lastName [2]*ent.Character, err error) {
	lastName[0], err = m.Character.Query().Where(character.CharEQ(last[0])).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			lastName[0] = makeFallbackChar(last[0])
		} else {
			return lastName, fmt.Errorf("query last name 0:%v", err)
		}
	}
	if last[1] != "" {
		lastName[1], err = m.Character.Query().Where(character.CharEQ(last[1])).First(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				lastName[1] = makeFallbackChar(last[1])
			} else {
				return lastName, fmt.Errorf("query last name 1:%v", err)
			}
		}
	}
	return lastName, nil
}

func ID(name string) string {
	sum := md5.Sum([]byte(name))
	return hex.EncodeToString(sum[:])
}

func New(client *ent.Client) *Repository {
	Logger("Repository")
	return &Repository{Client: client, cache: NewModelCache()}
}

func (m *Repository) QueryPoetryChars(ctx context.Context) ([]string, error) {
	chars, err := m.PoemChar.Query().
		Select(poemchar.FieldChar).
		Strings(ctx)
	if err != nil {
		return nil, fmt.Errorf("query poetry chars: %w", err)
	}
	seen := make(map[string]bool)
	var unique []string
	for _, ch := range chars {
		if !seen[ch] {
			seen[ch] = true
			unique = append(unique, ch)
		}
	}
	return unique, nil
}
