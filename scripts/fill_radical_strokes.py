"""
Fill remaining radical_stroke entries.

Remaining radicals not in stroke table:
- 攴 (radical 65, 4 strokes - tap/knock)
- 黃 (radical 201, 12 strokes - yellow, traditional form)
- 戶 (radical 63, 4 strokes - door, traditional form)
- 長 (radical 168, 8 strokes - long, traditional form)

Plus other traditional/rare radicals we may have missed.
"""
import json
import re

EXTRA_RADICAL_STROKES = {
    # Traditional radicals with their stroke counts
    "攴": 4,     # Radical 65 - tap/knock
    "黃": 12,    # Traditional yellow (simplified: 黄=11)
    "戶": 4,     # Traditional door/household (simplified: 户=4)
    "長": 8,     # Traditional long (simplified: 长=4)
    "齒": 15,    # Traditional tooth
    "龍": 16,    # Traditional dragon
    "龜": 16,    # Traditional turtle
    "馬": 10,    # Traditional horse
    "魚": 11,    # Traditional fish
    "鳥": 11,    # Traditional bird
    "風": 9,     # Traditional wind
    "飛": 9,     # Traditional fly
    "頁": 9,     # Traditional page/head
    "貝": 7,     # Traditional shell
    "見": 7,     # Traditional see
    "門": 8,     # Traditional gate
    "言": 7,     # Speech
    "糸": 6,     # Silk
    "金": 8,     # Metal
    "食": 9,     # Food/feed
    "車": 7,     # Car/vehicle
    "辵": 7,     # Walk
    "邑": 7,     # City
    "阜": 8,     # Mound
    "韋": 9,     # Leather traditional
    "麥": 11,    # Wheat traditional
    "齊": 14,    # Even traditional
    "黽": 13,    # Frog traditional
    "鹵": 11,    # Salt traditional
    "鼠": 13,    # Mouse (standalone char as radical)
    "鼻": 14,    # Nose
    "龠": 17,    # Flute
    # More rare radicals
    "禸": 5,     # Radical 111
    "癶": 5,     # Radical 102
    "舛": 6,     # Radical 132
    "臼": 6,     # Radical 131
    "艮": 6,     # Radical 135
    "釆": 7,     # Radical 162
    "鬯": 10,    # Radical 188
    "鬲": 10,    # Radical 189
    "黹": 12,    # Radical 200
    "鼎": 13,    # Radical 202
    "鼓": 13,    # Radical 203
    # Common components
    "亠": 2,
    "冂": 2,
    "冖": 2,
    "冫": 2,
    "勹": 2,
    "匚": 2,
    "尢": 3,
    "屮": 3,
    "彐": 3,
    "彡": 3,
    "彳": 3,
    "夂": 3,
    "夊": 3,
    "廴": 3,
    "弋": 3,
    "卩": 2,
    "厶": 2,
    "囗": 3,
    "忄": 3,
    "扌": 3,
    "艹": 3,
    "辶": 3,
    "钅": 5,
    "阝": 2,
    "饣": 3,
    "纟": 3,
    "讠": 2,
    "衤": 5,
    "礻": 4,
    "犭": 3,
    "罒": 5,
    "灬": 4,
    "氵": 3,
    "⺮": 6,
}

def main():
    with open('resources/character.json', 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    fixed = 0
    for c in data:
        if c.get('radical_stroke', 0) == 0 and c.get('radical'):
            r = c['radical']
            if r in EXTRA_RADICAL_STROKES:
                c['radical_stroke'] = EXTRA_RADICAL_STROKES[r]
                fixed += 1
    
    print(f"Fixed {fixed} additional radical_stroke entries")
    
    # Count remaining zeros
    still_zero = sum(1 for c in data if c.get('radical_stroke', 0) == 0 and c.get('radical'))
    print(f"Remaining with radical_stroke=0: {still_zero}")
    if still_zero > 0:
        remaining = set()
        for c in data:
            if c.get('radical_stroke', 0) == 0 and c.get('radical'):
                remaining.add(c['radical'])
        print(f"Unmapped radicals: {sorted(remaining)}")
    
    with open('resources/character.json', 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
    
    print("Done!")


if __name__ == '__main__':
    main()
