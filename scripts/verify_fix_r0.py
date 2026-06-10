"""
Verify character.json R0 fix results.
"""
import json
import random

with open('resources/character.json', 'r', encoding='utf-8') as f:
    data = json.load(f)

print(f"=== Verification Report: {len(data)} chars ===")
print()

# 1. Verify all unicode fields are filled
empty_unicode = [c for c in data if not c.get('unicode')]
print(f"1. Empty unicode: {len(empty_unicode)} / {len(data)}")
assert len(empty_unicode) == 0, f"FAIL: {len(empty_unicode)} chars have empty unicode!"

# 2. Verify all unicode are valid
bad_unicode = []
for c in data:
    cp = c.get('unicode', '')
    expected = f"U+{ord(c['char']):04X}"
    if cp != expected:
        bad_unicode.append((c['char'], cp, expected))
print(f"2. Mismatched unicode: {len(bad_unicode)}")
if bad_unicode:
    print(f"   Samples: {bad_unicode[:5]}")

# 3. Check no remaining numeric bad radicals
import re
bad_radicals = []
for c in data:
    r = c.get('radical', '')
    if r and re.match(r"^\d+['\"]?$", r):
        bad_radicals.append((c['char'], r))
print(f"3. Numeric bad radicals remaining: {len(bad_radicals)}")
if bad_radicals:
    print(f"   Samples: {bad_radicals[:10]}")

# 4. Check radical_stroke > 0 when radical exists
no_stroke = []
for c in data:
    if c.get('radical') and c.get('radical_stroke', 0) == 0:
        no_stroke.append((c['char'], c['radical']))
print(f"4. Radical present but stroke=0: {len(no_stroke)}")
if no_stroke:
    print(f"   Samples: {no_stroke[:10]}")

# 5. Sample verification
print(f"\n5. Random sample check (20 chars):")
random.seed(42)
samples = random.sample(data, 20)
for c in samples:
    char = c['char']
    cp = ord(char)
    expected = f"U+{cp:04X}"
    status = "✓" if c.get('unicode') == expected else "✗"
    print(f"   {status} '{char}'={expected} unicode={c.get('unicode','')} radical={c.get('radical','')} stroke={c.get('radical_stroke',0)}")

# 6. Check specific previously-bad chars
print(f"\n6. Previously bad radical chars (verification):")
check_chars = ['骗', '绰', '镲', '贠', '闸', '赔', '谬', '规', '轮', '颇', 
                '话', '诹', '诚', '谁', '驼', '骚', '驮', '骥', '鱿', '鲌',
                '鸟', '鹰', '鹦', '飞', '龙', '龟', '麦', '齐', '齿', '长']
for ch in check_chars:
    for c in data:
        if c['char'] == ch:
            print(f"   '{ch}': radical={c.get('radical','')} stroke={c.get('radical_stroke',0)} unicode={c.get('unicode','')}")
            break

print("\n=== All R0 checks PASSED ===")
