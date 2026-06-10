"""
Final radical_stroke fix for remaining unmapped radicals.

Note: Some radicals like "208''" appear to be double-quoted versions that need cleaning.
"""
import json
import re

# Remaining unmapped radicals and their stroke counts
FINAL_RADICAL_STROKES = {
    # These need their own right:
    "匸": 2,      # variant of Radical 22
    "巛": 3,      # River/water variant (radical 47 variant, 3 strokes
    "爿": 4,      # Wood/bed wooden table
    "襾": 6,      # Clothing variant (radical 145, 6 strokes
    "靑": 8,      # Green variant
    "面": 9,      # Face (also radical 176
}

def clean_radical(r):
    """Clean stray quote marks and normalize radicals"""
    if not r:
        return r, 0
    
    # Handle "208''" -> extract number
    match = re.match(r"(\d+)", r)
    if match:
        # If we get a number, it's a kangxi radical number
        num = int(match.group(1))
        # Number -> radical mapping we don't have. Let's see which one:
        if num == 208:
            return ("鼠", 13)  # 鼠 is radical 208 mouse
        if num == 145:
            return ("见", 4)
        return None, 0
    
    return r.strip().strip("'\"").strip("'\"")
    return r, 0

def main():
    with open('resources/character.json', 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    # First: clean any remaining number-quoted radicals
    fixed_radical = 0
    fixed_stroke = 0
    remaining_zero = []
    
    for c in data:
        r = c.get('radical', '')
        if not r:
            continue
        
        # Check if it's a bad number radical
        match = re.match(r"(\d+)['\"]*", r)
        if match:
            num = int(match.group(1))
            if num == 208:
                c['radical'] = '鼠'
                c['radical_stroke'] = 13
                fixed_radical += 1
            # Check the char '鼡' had "208''"
            continue
        
        # Fill radical_stroke still 0
        if c.get('radical_stroke', 0) == 0:
            if r in FINAL_RADICAL_STROKES:
                c['radical_stroke'] = FINAL_RADICAL_STROKES[r]
                fixed_stroke += 1
            else:
                remaining_zero.append((c['char'], r))
    
    print(f"Fixed radicals: {fixed_radical}")
    print(f"Fixed strokes: {fixed_stroke}")
    
    # Final check
    still_zero = sum(1 for c in data if c.get('radical_stroke', 0) == 0 and c.get('radical'))
    print(f"Remaining radical_stroke=0 with radical: {still_zero}")
    
    if still_zero > 0:
        remaining = set()
        for c in data:
            if c.get('radical_stroke', 0) == 0 and c.get('radical'):
                remaining.add(c['radical'])
        print(f"Remaining unmapped: {sorted(remaining)}")
    
    with open('resources/character.json', 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
    
    print("Done!")


if __name__ == '__main__':
    main()
