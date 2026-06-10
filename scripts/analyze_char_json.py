import json

with open('resources/character.json', 'r', encoding='utf-8') as f:
    data = json.load(f)

print(f'Total characters: {len(data)}')
print()

no_unicode = [c for c in data if c.get('unicode', '') == '']
print(f'Empty unicode: {len(no_unicode)} / {len(data)} ({100*len(no_unicode)/len(data):.1f}%)')

zero_stroke = [c for c in data if c.get('radical_stroke', 0) == 0]
print(f'Zero radical_stroke: {len(zero_stroke)} / {len(data)}')

has_trad = [c for c in data if c.get('traditional_chars')]
has_simp = [c for c in data if c.get('simplified_chars')]
has_variant = [c for c in data if c.get('variant_chars')]
print(f'With traditional_chars: {len(has_trad)}')
print(f'With simplified_chars: {len(has_simp)}')
print(f'With variant_chars: {len(has_variant)}')

has_simp_of = [c for c in data if c.get('simplified_of_char')]
has_var_of = [c for c in data if c.get('variant_of_char')]
print(f'With simplified_of_char: {len(has_simp_of)}')
print(f'With variant_of_char: {len(has_var_of)}')

print()
print('=== Radical analysis ===')
# Find non-Chinese radicals (like numbers, garbage)
bad_radicals = []
good_radicals = []
for c in data:
    r = c.get('radical', '')
    if not r:
        continue
    # Check if radical contains non-Chinese characters
    has_bad = any(ord(ch) < 0x4e00 or ord(ch) > 0x9fff for ch in r)
    if has_bad:
        bad_radicals.append((c['char'], r, c.get('radical_stroke', 0)))
    else:
        good_radicals.append((c['char'], r, c.get('radical_stroke', 0)))

print(f'Bad radicals (non-Chinese): {len(bad_radicals)}')
print(f'Good radicals (Chinese): {len(good_radicals)}')
print()

if bad_radicals:
    print('Bad radical samples:')
    for c, r, s in bad_radicals[:10]:
        print(f'  {c}: radical={repr(r)} stroke={s}')

# Check what the bad radical patterns are
bad_patterns = {}
for c, r, s in bad_radicals:
    bad_patterns[r] = bad_patterns.get(r, 0) + 1

print()
print('Unique bad radical values (top 20):')
for val, cnt in sorted(bad_patterns.items(), key=lambda x: -x[1])[:20]:
    print(f'  {repr(val)}: {cnt}')

print()
print('Good radical samples (first 10):')
for c, r, s in good_radicals[:10]:
    print(f'  {c}: radical={repr(r)} stroke={s}')
