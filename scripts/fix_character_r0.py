"""
R0 Fix: character.json data repair
====================================

Phase R0 fixes:
1. Unicode code point generation (U+XXXX format) for ALL characters
2. Bad radical replacement: numeric values (like "187'") -> actual Chinese radicals
3. radical_stroke: fill missing values using verified stroke counts

This script is idempotent: running it multiple times produces the same output.
"""
import json
import re
import sys

INPUT_FILE = "resources/character.json"
OUTPUT_FILE = "resources/character.json"

# Verified Kangxi radical number -> (radical_char, stroke_count)
# This was verified by examining actual characters with each bad radical number.
#
# Verification evidence:
#   120': 绰,绨,缬,纶,绥 -> 糸/纟 (silk radical)
#   147': 规,觎,览,觐,觑 -> 见/見 (see radical)
#   149': 谬,话,诹,诚,谁 -> 言/讠 (speech radical)
#   154': 贠,赔,贵,赌,贼 -> 贝/貝 (shell radical)
#   159': 辚,辏,轮,轹,辊 -> 车/車 (car radical)
#   167': 镲,锣,钭,镣,铤 -> 金/钅 (metal radical)
#   169': 闸,阄,阐,阊,闲 -> 门/門 (gate radical)
#   178': 韦,韪,韧,韩,韨 -> 韦/韋 (leather radical)
#   181': 颇,顸,颗,颧,频 -> 页/頁 (page/head radical)
#   182': 飏,飔,飚,飕,飙 -> 风/風 (wind radical)
#   184': 馄,饹,馇,饨,馋 -> 食/饣 (food radical)
#   187': 骗,驼,骚,驮,骥 -> 马/馬 (horse radical)
#   195': 鱿,鲌,鲋,鳐,鲊 -> 鱼/魚 (fish radical)
#   196': 鸟,鹰,鹦,鹁,鹙 -> 鸟/鳥 (bird radical)
#   199': 麦,麹,麸,麺 -> 麦/麥 (wheat radical)
#   183': 飞 -> 飞/飛 (fly radical)
#   168': 长 -> 长/長 (long radical)
#   205': 黾 -> 黾/黽 (frog radical)
#   208': 鼡 -> 鼠 (mouse radical)
#   210': 齑,齐 -> 齐/齊 (even radical)
#   211': 龉,龂,龀,龆,龇 -> 齿/齒 (tooth radical)
#   212': 龚,龛,龙 -> 龙/龍 (dragon radical)
#   213': 龟 -> 龟/龜 (turtle radical)
#   197': 鹾 -> 卤/鹵 (salt radical)
#   210': 龠 -> 龠 (flute radical)
#
# Simplified form strokes: these are used because character.json uses simplified chars.
RADICAL_FIX_MAP = {
    "120": ("纟", 3),
    "147": ("见", 4),
    "149": ("讠", 2),
    "154": ("贝", 4),
    "159": ("车", 4),
    "167": ("钅", 5),
    "168": ("长", 4),
    "169": ("门", 3),
    "178": ("韦", 4),
    "181": ("页", 6),
    "182": ("风", 4),
    "183": ("飞", 3),
    "184": ("饣", 3),
    "187": ("马", 3),
    "195": ("鱼", 8),
    "196": ("鸟", 5),
    "197": ("卤", 7),
    "199": ("麦", 7),
    "205": ("黾", 8),
    "208": ("鼠", 13),
    "210": ("齐", 6),
    "211": ("齿", 8),
    "212": ("龙", 5),
    "213": ("龟", 7),
}

# Known good radicals -> stroke count (for filling missing radical_stroke)
RADICAL_STROKE_MAP = {
    "一": 1, "丨": 1, "丶": 1, "丿": 1, "乙": 1, "亅": 1,
    "二": 2, "亠": 2, "人": 2, "儿": 2, "入": 2, "八": 2, "冂": 2, "冖": 2, "冫": 2,
    "几": 2, "凵": 2, "刀": 2, "力": 2, "勹": 2, "匕": 2, "匚": 2, "十": 2, "卜": 2,
    "卩": 2, "厂": 2, "厶": 2, "又": 2,
    "口": 3, "囗": 3, "土": 3, "士": 3, "夂": 3, "夊": 3, "夕": 3, "大": 3, "女": 3,
    "子": 3, "宀": 3, "寸": 3, "小": 3, "尢": 3, "尸": 3, "屮": 3, "山": 3, "川": 3,
    "工": 3, "己": 3, "巾": 3, "干": 3, "幺": 3, "广": 3, "廴": 3, "弋": 3, "弓": 3,
    "彐": 3, "彡": 3, "彳": 3,
    "心": 4, "戈": 4, "户": 4, "手": 4, "支": 4, "攵": 4, "文": 4, "斗": 4, "斤": 4,
    "方": 4, "无": 4, "日": 4, "曰": 4, "月": 4, "木": 4, "欠": 4, "止": 4, "歹": 4,
    "殳": 4, "毋": 4, "比": 4, "毛": 4, "氏": 4, "气": 4, "水": 4, "火": 4, "爪": 4,
    "父": 4, "爻": 4, "片": 4, "牙": 4, "牛": 4, "犬": 4,
    "玄": 5, "玉": 5, "瓜": 5, "瓦": 5, "甘": 5, "生": 5, "用": 5, "田": 5, "疋": 5,
    "疒": 5, "癶": 5, "白": 5, "皮": 5, "皿": 5, "目": 5, "矛": 5, "矢": 5, "石": 5,
    "示": 5, "禸": 5, "禾": 5, "穴": 5, "立": 5,
    "竹": 6, "米": 6, "糸": 6, "缶": 6, "网": 6, "羊": 6, "羽": 6, "老": 6, "而": 6,
    "耒": 6, "耳": 6, "聿": 6, "肉": 6, "臣": 6, "自": 6, "至": 6, "臼": 6, "舌": 6,
    "舛": 6, "舟": 6, "艮": 6, "色": 6, "艸": 6, "虍": 6, "虫": 6, "血": 6, "行": 6,
    "衣": 6, "西": 6,
    "见": 7, "角": 7, "言": 7, "谷": 7, "豆": 7, "豕": 7, "豸": 7, "贝": 7, "赤": 7,
    "走": 7, "足": 7, "身": 7, "车": 7, "辛": 7, "辰": 7, "辵": 7, "邑": 7, "酉": 7,
    "釆": 7, "里": 7, "见": 4, "讠": 2, "贝": 4, "车": 4,
    "金": 8, "长": 4, "门": 3, "阜": 8, "隶": 8, "隹": 8, "雨": 8, "青": 8, "非": 8,
    "钅": 5, "长": 4, "门": 3,
    "革": 9, "韦": 4, "韭": 9, "音": 9, "页": 6, "风": 4, "飞": 3, "食": 9, "首": 9,
    "香": 9, "页": 6, "风": 4, "飞": 3, "饣": 3,
    "马": 3, "骨": 10, "高": 10, "髟": 10, "斗": 10, "鬯": 10, "鬲": 10, "鬼": 10,
    "马": 3,
    "鱼": 8, "鸟": 5, "卤": 7, "鹿": 11, "麦": 7, "麻": 11, "鱼": 11, "鸟": 11,
    "黄": 12, "黍": 12, "黑": 12, "黹": 12, "黄": 12,
    "黾": 13, "鼎": 13, "鼓": 13, "鼠": 13,
    "鼻": 14, "齐": 6,
    "齿": 8, "龙": 5,
    "龟": 7, "龠": 17,
}

# Common simplified radicals with their stroke counts (overrides)
SIMPLIFIED_RADICAL_STROKES = {
    "纟": 3, "讠": 2, "钅": 5, "贝": 4, "见": 4, "门": 3,
    "车": 4, "马": 3, "鱼": 8, "鸟": 5, "风": 4, "飞": 3,
    "饣": 3, "页": 6, "韦": 4, "长": 4, "龙": 5, "齐": 6,
    "齿": 8, "龟": 7, "卤": 7, "麦": 7, "黾": 8, "鼠": 13,
    "龠": 17,
}

# Build consolidated stroke map
RADICAL_STROKE_MAP.update(SIMPLIFIED_RADICAL_STROKES)


def get_unicode_code_point(char):
    """Generate U+XXXX unicode code point for a character."""
    if not char or len(char) == 0:
        return ""
    cp = ord(char[0])
    return f"U+{cp:04X}"


def is_bad_radical(radical):
    """Check if radical is a numeric bad value like '187'."""
    if not radical:
        return False
    # Match patterns like "187'", "120", etc.
    return bool(re.match(r"^\d+['\"]?$", radical.strip()))


def fix_radical(radical_str):
    """Fix a bad numeric radical string -> (correct radical char, stroke count)."""
    match = re.match(r"(\d+)", radical_str.strip())
    if not match:
        return None, 0
    
    num_str = match.group(1)
    if num_str in RADICAL_FIX_MAP:
        return RADICAL_FIX_MAP[num_str]
    return None, 0


def get_radical_stroke(radical):
    """Get stroke count for a radical character."""
    if not radical:
        return 0
    # Check simplified first (more common in this dataset)
    if radical in SIMPLIFIED_RADICAL_STROKES:
        return SIMPLIFIED_RADICAL_STROKES[radical]
    if radical in RADICAL_STROKE_MAP:
        return RADICAL_STROKE_MAP[radical]
    return 0


def main():
    print(f"Reading {INPUT_FILE}...")
    with open(INPUT_FILE, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    print(f"Total characters: {len(data)}")
    
    unicode_fixed = 0
    radical_fixed = 0
    stroke_fixed = 0
    unicode_already_good = 0
    radical_already_good = 0
    stroke_already_good = 0
    
    for i, c in enumerate(data):
        char = c.get('char', '')
        if not char:
            continue
        
        # 1. Fix unicode
        current_unicode = c.get('unicode', '')
        expected_unicode = get_unicode_code_point(char)
        if not current_unicode:
            c['unicode'] = expected_unicode
            unicode_fixed += 1
        elif current_unicode == expected_unicode:
            unicode_already_good += 1
        
        # 2. Fix bad radical
        current_radical = c.get('radical', '')
        if is_bad_radical(current_radical):
            new_radical, new_stroke = fix_radical(current_radical)
            if new_radical:
                c['radical'] = new_radical
                radical_fixed += 1
                # Also fix radical_stroke if it was 0
                if c.get('radical_stroke', 0) == 0 and new_stroke > 0:
                    c['radical_stroke'] = new_stroke
                    stroke_fixed += 1
        elif current_radical:
            radical_already_good += 1
        
        # 3. Fill missing radical_stroke if we have the radical but stroke is 0
        if c.get('radical_stroke', 0) == 0 and c.get('radical'):
            stroke = get_radical_stroke(c['radical'])
            if stroke > 0:
                c['radical_stroke'] = stroke
                stroke_fixed += 1
            else:
                stroke_already_good += 1
        elif c.get('radical_stroke', 0) > 0:
            stroke_already_good += 1
    
    print("\n=== Fix Summary ===")
    print(f"Unicode fixed:        {unicode_fixed}")
    print(f"Radicals fixed:       {radical_fixed}")
    print(f"Radical strokes fixed:{stroke_fixed}")
    print(f"Unicode already good: {unicode_already_good}")
    print(f"Radicals already good:{radical_already_good}")
    print(f"Strokes already good: {stroke_already_good}")
    
    # Print a few samples to verify
    print("\n=== Sample verification (first 10 chars) ===")
    for c in data[:10]:
        print(f"  '{c['char']}' unicode={c.get('unicode','')} radical={c.get('radical','')} stroke={c.get('radical_stroke',0)}")
    
    # Print some previously-bad radical chars
    print("\n=== Previously bad radical chars (checking) ===")
    found = 0
    for c in data:
        char = c['char']
        if char in ['骗', '绰', '镲', '贠', '闸', '赔', '谬', '规', '轮', '颇']:
            print(f"  '{char}' radical={c.get('radical','')} stroke={c.get('radical_stroke',0)}")
            found += 1
            if found >= 10:
                break
    
    print(f"\nWriting {OUTPUT_FILE}...")
    with open(OUTPUT_FILE, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
    
    print("Done!")
    return 0


if __name__ == '__main__':
    sys.exit(main())
