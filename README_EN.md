# fate

<p align="center">
  <strong>fate — Bazi Wuxing Naming Algorithm Engine</strong><br>
  <em>Intelligent naming algorithm engine and CLI based on Bazi Wuxing, Sancai Wuge, and Zhouyi hexagrams</em>
</p>

<p align="center">
  English | <a href="README.md">中文</a>
</p>

---

## Features

- Bazi Calculation — Four Pillars, Wuxing strength, Tiaohou Shen
- Dual Xi-Yong Algorithm — Balance Method + GeJu Method (10 pattern types)
- Sancai Wuge — 81 Dayan numbers, Yin-Yang Wuxing, O(1) lookup
- Zhouyi Hexagrams — 64 hexagram interpretations
- 5-Dimension Scoring — Cultural/Wuxing/Zodiac/Wuge/Yinyun

---

## About fate and qiming

**fate** is an open-source naming algorithm engine that provides core Bazi calculation, Xi-Yong analysis, Wuge filtering, and name generation capabilities, delivered as both a CLI tool and a Go library.

**qiming** is a commercial naming service built on top of fate, providing a web interface, poetry-based naming, and other end-user features.

| | fate | qiming |
|---|---|---|
| Purpose | Algorithm engine + CLI | Commercial service |
| Bazi/Xi-Yong/5D scoring | Included | Included |
| Poetry naming | Not included | Included |
| Web interface | Not included | Included |
| Open source | Yes | No |

---

## Quick Start

### Requirements

- Go 1.22+

### Install

```bash
git clone https://github.com/babyname/fate.git
cd fate
go mod download
```

---

## Two Usage Modes

### Mode 1: Command Line

Fast name generation with clean output:

```bash
# Generate names
go run ./cmd/console name -s 张 -b "2024/06/15 10:30" -g boy

# View detailed analysis for a specific name
go run ./cmd/console name detail 峰 瑞 -s 张 -b "2024/06/15 10:30" -g boy

# View all options
go run ./cmd/console name -h
```

**Options:**

| Flag | Description | Example |
|------|-------------|---------|
| `-s, --surname` | Surname | `-s 张` |
| `-b, --born` | Birth date | `-b "2024/06/15 10:30"` |
| `-g, --gender` | Gender | `-g boy` or `-g girl` |
| `--xiyong` | Xi-Yong algorithm | `--xiyong balance` or `--xiyong geju` |
| `--strictness` | Wuge filter strictness | `--strictness moderate` |
| `-f, --filter` | Filter out specific characters | `-f 病死` |
| `-o, --output` | Output to file | `-o result.txt` |

---

### Mode 2: Library Import

Integrate into your Go project:

```bash
go get github.com/babyname/fate
```

```go
package main

import (
    "fmt"
    "time"

    "github.com/babyname/fate"
    "github.com/babyname/fate/config"
)

func main() {
    cfg := config.DefaultConfig()

    f, err := fate.New(cfg)
    if err != nil {
        panic(err)
    }

    filter := fate.NewFilter(fate.FilterOption{
        CharacterFilter:     true,
        CharacterFilterType: fate.CharacterFilterTypeDefault,
        MinStroke:           3,
        MaxStroke:           18,
        RegularFilter:       true,
        DaYanFilter:         true,
        WuXingFilter:        true,
    })

    s := f.NewSessionWithFilter(filter)

    born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")
    input := &fate.Input{
        Last: [2]string{"张", ""},
        Born: born,
        Sex:  fate.SexBoy,
    }

    err = s.Start(input)
    if err != nil {
        panic(err)
    }
    s.Wait()

    output := input.Output()
    fmt.Printf("Generated %d names\n", output.Total())

    for _, nr := range output.TopNames() {
        fmt.Printf("  %s - Score: %.1f\n", nr.FullName, nr.Score)
    }
}
```

**Core API:**

| Type | Description |
|------|-------------|
| `fate.New(cfg)` | Create Fate instance |
| `fate.NewSessionWithFilter(filter)` | Create session with filter options |
| `fate.NewFilter(option)` | Create filter from options |
| `fate.FilterOption{}` | Filter options (stroke range / Wuxing / Dayan / gender etc.) |
| `fate.Input{}` | Input parameters (surname / birthday / gender) |
| `fate.Output` | Output results (TopNames / AllNames / Total) |
| `fate.Session` | Session interface (Start / Stop / Wait) |

---

## Xi-Yong Algorithms

### Balance Method

Determines Xi-Yong based on Day Master's strength (same vs opposite camp):

| Day Master | Yong Shen | Xi Shen | Ji Shen |
|------------|-----------|---------|---------|
| Strong | Officer/Kill (克制) | Wealth+Resource (我克+生我) | Peer/Rob (同我) |
| Weak | Resource (生我) | Peer/Rob (同我) | Officer+Output (克我+我生) |

### GeJu Method

First determines pattern type (10 types: Zheng Guan, Qi Sha, Shi Shen, Shang Guan, Zheng Cai, Pian Cai, Jian Lu, Yang Ren, etc.), then selects Xi-Yong accordingly.

---

## 5-Dimension Scoring

| Dimension | Weight | Description |
|-----------|--------|-------------|
| Cultural Impression | 20% | Common characters, regular forms, meaning richness |
| Wuxing Bazi | 25% | Name Wuxing match with Xi-Yong |
| Zodiac | 10% | Zodiac and name Wuxing relationship |
| Wuge Shuli | 25% | Tian/Ren/Di/Wai/Zong Ge fortune |
| Yinyun | 20% | Name phonetic harmony |

---

## Tech Stack

- **Go 1.22+** — Main language
- **Ent ORM** — Database ORM
- **SQLite3** — Data storage
- **chronos/v2** — Bazi calculation
- **yi** — Zhouyi hexagrams

---

## License

MIT License
