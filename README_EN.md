# 命运起名 Fate

<p align="center">
  <strong>Modern Scientific Naming Tool — Intelligent Naming System Based on Bazi Wuxing · Sancai Wuge · Zhouyi Hexagrams</strong>
</p>

<p align="center">
  <img src="docs/images/architecture.svg" alt="Architecture" width="800"/>
</p>

<p align="center">
  English | <a href="README.md">中文</a>
</p>

---

## ✨ Features

- 🎂 **Bazi Calculation** — Four Pillars, Wuxing strength, Tiaohou Shen
- ☯️ **Dual Xi-Yong Algorithm** — Balance Method + GeJu Method (10 pattern types)
- 📐 **Sancai Wuge** — 81 Dayan numbers, Yin-Yang Wuxing, O(1) precomputed lookup (3.89ns/op)
- 🔮 **Zhouyi Hexagrams** — 64 hexagram interpretations (Daxiang/Career/Business/Fame/Marriage/Decision)
- 📊 **4-Dimension Scoring** — Cultural Impression / Wuxing Bazi / Zodiac / Wuge Shuli
- 📝 **Multi-format Output** — Text / Markdown / JSON
- 🏛️ **Simplified-Traditional Mapping** — 400+ character mapping table

## 📸 Report Preview

<p align="center">
  <img src="docs/images/report_preview.svg" alt="Report Preview" width="800"/>
</p>

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- GCC (SQLite3 CGO dependency)
- SQLite3

### Install

```bash
git clone https://github.com/babyname/fate.git
cd fate
go mod download
```

### Generate Name Analysis Report

```bash
go run github.com/babyname/fate/cmd/gen_report

# Output files in output/ directory
ls output/
# 张_姓名分析报告_平衡用神法.txt
# 张_姓名分析报告_平衡用神法.md
# 张_姓名分析报告_平衡用神法.json
# 张_姓名分析报告_格局用神法.txt
# 张_姓名分析报告_格局用神法.md
# 张_姓名分析报告_格局用神法.json
```

### Usage in Code

```go
package main

import (
    "time"

    "github.com/babyname/fate/internal/analysis"
    v2 "github.com/godcong/chronos/v2"
)

func main() {
    born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")

    // Balance Method
    fateData, _ := v2.GetFateData(&v2.FateInput{
        BirthDate: born,
        Gender:    1,
        Surname:   "张",
        Method:    v2.XiYongMethodBalance,
    })

    // GeJu Method
    fateData2, _ := v2.GetFateData(&v2.FateInput{
        BirthDate: born,
        Gender:    1,
        Surname:   "张",
        Method:    v2.XiYongMethodGeJu,
    })

    // Build name result
    c1 := &ent.Character{Char: "驰", WuXing: "火", ScienceStroke: 13, ...}
    c2 := &ent.Character{Char: "筎", WuXing: "木", ScienceStroke: 12, ...}
    result := analysis.BuildNameResult(1, "张", c1, c2, 11, 0, fateData)

    // Generate report
    report := analysis.NewReport("张", "2024年06月15日", "男", fateData, 1000)
    report.TopNames = append(report.TopNames, result)

    f := &analysis.MarkdownFormatter{}
    f.Format(os.Stdout, report)
}
```

## 🏗️ Project Structure

```
fate/
├── cmd/gen_report/           # Report generation tool
├── ent/                      # Ent ORM entity definitions
├── internal/
│   ├── analysis/             # Core analysis module
│   │   ├── analysis.go       # Data structures + formatters
│   │   ├── sancai_data.go    # Sancai/Jichuyun/Chenggongyun/Renji data
│   │   ├── zhouyi.go         # Zhouyi hexagram calculation
│   │   ├── zhouyi_data.go    # 64 hexagram interpretation data
│   │   └── simplified_traditional.go  # Simplified-Traditional mapping
│   ├── wuge/                 # Sancai Wuge (3-5 elements)
│   ├── wuxing/               # Wuxing analysis
│   ├── rating/               # Scoring system
│   ├── filter/               # Name filtering
│   ├── naming/               # Naming logic
│   └── session/              # Session management
├── chronos/                  # Bazi calculation submodule
│   ├── fate.go               # FateData main entry
│   ├── xiyong_balance.go     # Balance Xi-Yong method
│   ├── xiyong_geju.go        # GeJu Xi-Yong method
│   └── fate_helpers.go       # Wuxing calculation helpers
└── yi/                       # Zhouyi hexagram submodule
```

## ☯️ Two Xi-Yong Shen Algorithms

### Balance Method

Based on day-master strength, using ally/enemy force comparison:

| Day Master | Yong Shen | Xi Shen | Ji Shen | Chou Shen |
|-----------|-----------|---------|---------|-----------|
| Strong | Officer (Ke-Wo) | Wealth + Seal | Peer (Bi-Jie) | Ji Shen's parent |
| Weak | Seal (Sheng-Wo) | Peer (Bi-Jie) | Officer + Food | Ji Shen's parent |

### GeJu (Pattern) Method

First determine the pattern (10 types), then select Yong Shen based on pattern + strength:

| Pattern | Strong Yong | Weak Yong |
|---------|------------|-----------|
| Zheng Guan | Officer | Seal |
| Qi Sha | Food controls Sha | Seal transforms Sha |
| Shi Shen | Food generates Wealth | Peer helps Self |
| Shang Guan | Seal controls Shang | Seal |
| Zheng/Pian Cai | Wealth | Seal |
| Jian Lu/Yang Ren | Officer | Seal |

## 📊 Scoring System

| Dimension | Description |
|-----------|-------------|
| Cultural Impression | Common characters, regular script, meaning richness |
| Wuxing Bazi | Name Wuxing matching with Xi-Yong Shen |
| Zodiac | Zodiac Wuxing generation/restriction with name |
| Wuge Shuli | Tian/Ren/Di/Wai/Zong Ge luck/inauspicious |

## 🔧 Performance

| Metric | Value |
|--------|-------|
| WuGeLucky Lookup | 3.89 ns/op, 0 B alloc |
| Zhang Surname Full | ~50ms / 61,730 names |
| Li Surname Full | ~127ms / 165,852 names |

## 🛠️ Tech Stack

- **Go 1.22+** — Main language
- **Ent ORM** — Database ORM
- **SQLite3** — Data storage (`github.com/sqlite3ent/sqlite3` driver)
- **chronos/v2** — Bazi calculation submodule (local `replace`)
- **yi** — Zhouyi hexagram calculation (`github.com/godcong/yi`)

## 📜 License

MIT License
