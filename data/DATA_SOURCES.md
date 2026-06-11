# Data Sources & Attribution

## chinese-poetry

| Field | Value |
|-------|-------|
| **Name** | chinese-poetry: 最全中文诗歌古典文集数据库 |
| **Repository** | https://github.com/chinese-poetry/chinese-poetry |
| **License** | [MIT](https://github.com/chinese-poetry/chinese-poetry/blob/master/LICENSE) |
| **Used Version** | commit `99ebbef` (May 2026) |
| **What we use** | `全唐诗/唐诗三百首.json`, `宋词/宋词三百首.json`, `诗经/shijing.json` |
| **Usage in fate** | Poetry origin lookup — when a character in a generated name appears in classical poetry, show the source verse |
| **Download** | `./download/chinese_poetry.ps1` |

## Unihan (Unicode Han Database)

| Field | Value |
|-------|-------|
| **Name** | Unihan Database (IRG Sources, Readings, Variants, RadicalStrokeCounts) |
| **Source** | https://www.unicode.org/Public/UCD/latest/ucd/Unihan.zip |
| **License** | [Unicode License](https://www.unicode.org/license.html) — free to use |
| **Used Version** | Unicode 16.0 |
| **What we use** | `Unihan_IRGSources.txt`, `Unihan_Readings.txt`, `Unihan_Variants.txt`, `Unihan_RadicalStrokeCounts.txt`, `Unihan_DictionaryLikeData.txt` |
| **Usage in fate** | Building `resources/character.json` — radical, stroke count, simplified/traditional mapping, variant mapping |
| **Download** | `./download/unihan.ps1` |

## xinhua (新华字典)

| Field | Value |
|-------|-------|
| **Name** | xinhua dictionary (新华字典) JSON dataset |
| **Source** | https://github.com/pwxcoo/chinese-xinhua (reference) |
| **License** | Open data — free for non-commercial and commercial use |
| **What we use** | `word.json` — Chinese character definitions, pinyin, strokes |
| **Usage in fate** | Filling `resources/character.json` with Chinese meanings (via `dictctl fill-xinhua`) |
| **Download** | `./download/xinhua.ps1` |

## character.json

| Field | Value |
|-------|-------|
| **Location** | `resources/character.json` |
| **Status** | Tracked in git (core project data, ~12MB) |
| **Built from** | Unihan + xinhua + manual curation |
| **Role** | The central character database — all naming features depend on it |
