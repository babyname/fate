<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/version-4.2.0-blue" alt="Version" />
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License" />
  <img src="https://img.shields.io/badge/module-v4-7B42BC" alt="Module v4" />
</p>

<p align="center">
  <h1 align="center">🌌 fate · Celestial Naming Engine</h1>
  <p align="center">
    <strong>Intelligent Chinese name generator powered by Bazi (八字), Wuxing (五行), Zhouyi (周易), and classical poetry</strong><br>
    <em>FATE — Celestial Naming Engine</em>
  </p>
</p>

<p align="center">
  English | <a href="README.md">中文</a>
</p>

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🎯 Bazi Analysis
- Four Pillars (year/month/day/hour) precision calculation
- Wuxing strength analysis + Tiaohou Yongshen
- **Dual Xi-Yong algorithm**: Balance Method & GeJu Method (10 pattern types)
- Chinese zodiac integration

</td>
<td width="50%">

### 📊 Multi-Dimension Scoring
- **5-dimension scoring**: Cultural · Wuxing · Zodiac · Wuge · Phonetics
- Sancai configuration analysis
- 81 Dayan numbers O(1) lookup
- 64 Zhouyi hexagram interpretations

</td>
</tr>
<tr>
<td width="50%">

### 🌐 Modern Web UI
- React 19 + Tailwind CSS responsive design
- Real-time streaming generation with polling
- Candidate table & card dual view
- Explore mode: random sampling with filters

</td>
<td width="50%">

### ⚡ High Performance
- Ent ORM + SQLite for efficient data access
- ExcellentTable streaming Top-N data structure
- Async concurrent generation (session mode)
- Embedded web frontend (single binary distribution)

</td>
</tr>
</table>

---

## 📸 Screenshots

> 💡 **Tip**: Start the server and visit `http://localhost:18080`.

| Home | Results |
|------|---------|
| ![Home](docs/screenshots/home.png) | ![Results](docs/screenshots/results.png) |

| Candidate Table | Card View |
|-----------------|-----------|
| ![Table](docs/screenshots/table.png) | ![Cards](docs/screenshots/cards.png) |

---

## 🚀 Quick Start

### Requirements

- **Go** 1.22+
- **Bun** (for frontend dev only; not needed for production)

### Install & Run

```bash
# Clone
git clone https://github.com/babyname/fate.git
cd fate

# Download dependencies
go mod download

# Initialize database (required)
go run ./cmd/dbinit

# Start web server
go run ./cmd/server
# Open http://localhost:18080
```

### CLI

```bash
# Build
go build -o fate ./cmd/console

# Generate names
./fate name -s Zhang -b "2024/06/15 10:30" -g boy

# Detailed analysis
./fate name detail Shi Lun -s Zhang -b "2024/06/15 10:30" -g boy
```

### CLI Options

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--surname` | `-s` | Surname | `-s Zhang` |
| `--born` | `-b` | Birth date & time | `-b "2024/06/15 10:30"` |
| `--gender` | `-g` | Gender | `-g boy` / `-g girl` |
| `--xiyong` | | Xi-Yong algorithm | `balance` / `geju` |
| `--strictness` | | Wuge strictness | `loose` / `moderate` / `strict` |
| `--filter` | `-f` | Filter characters | `-f 病死穷` |
| `--output` | `-o` | Output to file | `-o result.txt` |

---

## 🔌 Web API

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/generate` | Async name generation (returns task_id) |
| `GET` | `/api/task/:id` | Query task status |
| `GET` | `/api/names/:taskId` | Paginated name list |
| `GET` | `/api/name-detail` | Detailed name analysis |
| `GET` | `/api/explore/:taskId` | Random name exploration |

---

## 🧩 Library Usage

```bash
go get github.com/babyname/fate/v4
```

```go
package main

import (
    "fmt"
    "time"

    "github.com/babyname/fate/v4"
    "github.com/babyname/fate/v4/config"
)

func main() {
    cfg := config.DefaultConfig()
    f, _ := fate.New(cfg)

    filter := fate.NewFilter(fate.FilterOption{
        CharacterFilter:     true,
        MinStroke:           3,
        MaxStroke:           18,
        RegularFilter:       true,
        DaYanFilter:         true,
        WuXingFilter:        true,
        AvoidCharacters:     []string{"病", "死"},
    })

    s := f.NewSessionWithFilter(filter)
    born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")
    input := &fate.Input{
        Last: [2]string{"Zhang", ""},
        Born: born,
        Sex:  fate.SexBoy,
    }

    s.Start(input)
    s.Wait()

    output := input.Output()
    for _, nr := range output.TopNames() {
        fmt.Printf("%s - %.1f (%s)\n", nr.FullName, nr.Score, nr.Grade)
    }

    // ExcellentTable: streaming Top-N access
    table := output.GetExcellentTable()
    for _, e := range table.TopN(10) {
        fmt.Printf("%s%s - %.1f has_poetry=%v\n", e.Char1, e.Char2, e.Score, e.HasPoetry)
    }
}
```

---

## 🏛️ Architecture

```
┌─────────────────────────────────────────────┐
│              fate CLI / Server               │
├─────────────────────────────────────────────┤
│  cmd/console     │  cmd/server               │
│  (CLI tool)      │  (HTTP API + embedded UI) │
├─────────────────────────────────────────────┤
│           internal/http (API routes)         │
├─────────────────────────────────────────────┤
│              fate Core Engine                │
│  ┌──────────┬──────────┬──────────────────┐ │
│  │ Session  │ Filter   │ ExcellentTable   │ │
│  │ Async    │ Multi    │ Streaming        │ │
│  │ Gen      │ Stage    │ Top-N            │ │
│  ├──────────┼──────────┼──────────────────┤ │
│  │ Bazi     │ Wuge     │ Zhouyi           │ │
│  └──────────┴──────────┴──────────────────┘ │
├─────────────────────────────────────────────┤
│          chronos/v2 (Bazi engine)            │
│          yi (Zhouyi engine)                  │
│          Ent ORM (data layer)                │
├─────────────────────────────────────────────┤
│          SQLite3 (data storage)              │
│          character.json (character DB)       │
│          chinese-poetry (poetry corpus)      │
└─────────────────────────────────────────────┘
```

---

## 📐 5-Dimension Scoring

| Dimension | Weight | Criteria |
|-----------|:------:|----------|
| 🎨 **Cultural** | 20% | Common usage, regular script, semantic richness |
| 🔥 **Wuxing** | 25% | Name Wuxing vs Day Master Xi-Yong match |
| 🐉 **Zodiac** | 10% | Zodiac Wuxing vs name Wuxing interaction |
| 📏 **Wuge** | 25% | Tian/Ren/Di/Wai/Zong Ge fortune scores |
| 🎵 **Phonetics** | 20% | Syllable tone harmony, ping-ze balance |

---

## 🔀 Branch Strategy

| Branch | Purpose | go.mod module |
|--------|---------|---------------|
| `main` | v4.x active development | `github.com/babyname/fate/v4` |
| `v3` | v3.x bugfix maintenance | `github.com/babyname/fate` |

> ⚠️ **Historical**: v3 branch lacks `/v3` suffix (pre-dates stricter Go module conventions).
> From v4 onward, `/v4` suffix is consistently used. Future v5 will use `/v5`.

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.22+ |
| Frontend | React 19 + Tailwind CSS 3 + Rsbuild |
| State | Zustand |
| UI Kit | Radix UI + Lucide Icons |
| ORM | Ent |
| Storage | SQLite3 (modernc, CGo-free) |
| Bazi | chronos/v2 |
| Zhouyi | yi |

---

## 🧪 Development

```bash
# Frontend (hot reload)
cd web && bun install && bun run dev

# Backend
go run ./cmd/server

# Production build
cd web && bun run build
go build -o fate-server ./cmd/server

# Tests
go test ./...

# Lint
golangci-lint run
```

---

## 📄 License

[MIT](LICENSE) © babyname
