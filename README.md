# 启命宝 fate

<p align="center">
  <strong>fate — 八字五行起名算法引擎</strong><br>
  <em>基于八字五行·三才五格·周易卦象的智能起名算法引擎与命令行工具</em>
</p>

<p align="center">
  <a href="README_EN.md">English</a> | 中文
</p>

---

## 功能特性

- 八字计算 — 四柱八字、五行强弱、调候用神
- 双算法喜用神 — 平衡用神法 + 格局用神法（10种格局）
- 三才五格 — 81数大衍、阴阳五行、O(1)查找
- 周易卦象 — 64卦完整解读
- 五维评分 — 文化印象/五行八字/生肖/五格数理/音韵

---

## 关于 fate 与 qiming

**fate** 是开源的起名算法引擎，提供核心的八字计算、喜用神分析、五格筛选和名字生成能力，以命令行工具和 Go 库两种形式交付。

**qiming** 是基于 fate 构建的商业起名服务，提供 Web 界面、诗词取名等面向终端用户的功能。

| | fate | qiming |
|---|---|---|
| 定位 | 算法引擎 + CLI | 商业服务 |
| 八字/喜用神/五维评分 | 包含 | 包含 |
| 诗词取名 | 不包含 | 包含 |
| Web 界面 | 不包含 | 包含 |
| 开源 | 是 | 否 |

---

## 快速开始

### 环境要求

- Go 1.22+

### 安装

```bash
git clone https://github.com/babyname/fate.git
cd fate
go mod download
```

---

## 两种使用方式

### 方式一：命令行起名

适合快速生成名字列表，输出简洁高效：

```bash
# 生成名字
go run ./cmd/console name -s 张 -b "2024/06/15 10:30" -g boy

# 查看单个名字的详细分析
go run ./cmd/console name detail 峰 瑞 -s 张 -b "2024/06/15 10:30" -g boy

# 查看所有参数
go run ./cmd/console name -h
```

**参数说明：**

| 参数 | 说明 | 示例 |
|------|------|------|
| `-s, --surname` | 姓氏 | `-s 张` |
| `-b, --born` | 出生日期 | `-b "2024/06/15 10:30"` |
| `-g, --gender` | 性别 | `-g boy` 或 `-g girl` |
| `--xiyong` | 喜用神算法 | `--xiyong balance` 或 `--xiyong geju` |
| `--strictness` | 五格筛选严格度 | `--strictness moderate` |
| `-f, --filter` | 从结果中过滤指定汉字 | `-f 病死` |
| `-o, --output` | 输出到文件 | `-o result.txt` |

---

### 方式二：代码导入使用

适合集成到其他 Go 项目中：

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
    fmt.Printf("共生成 %d 个名字\n", output.Total())

    for _, nr := range output.TopNames() {
        fmt.Printf("  %s - 评分: %.1f\n", nr.FullName, nr.Score)
    }
}
```

**核心 API：**

| 类型 | 说明 |
|------|------|
| `fate.New(cfg)` | 创建 Fate 实例 |
| `fate.NewSessionWithFilter(filter)` | 创建带筛选条件的会话 |
| `fate.NewFilter(option)` | 根据选项创建筛选器 |
| `fate.FilterOption{}` | 筛选选项（笔画范围/五行/大衍/性别等） |
| `fate.Input{}` | 输入参数（姓氏/生日/性别） |
| `fate.Output` | 输出结果（TopNames/AllNames/Total） |
| `fate.Session` | 会话接口（Start/Stop/Wait） |

---

## 喜用神算法

### 平衡用神法

基于日主强弱，通过同党/异党力量对比取用神：

| 日主状态 | 用神 | 喜神 | 忌神 |
|----------|------|------|------|
| 日主强 | 克我者（官杀） | 我克者+生我者 | 同我者（比劫） |
| 日主弱 | 生我者（印星） | 同我者（比劫） | 克我者+我生者 |

### 格局用神法

先定格局（正官/七杀/食神/伤官/正财/偏财/建禄/阳刃等10种），再根据格局取用神。

---

## 五维评分体系

| 维度 | 权重 | 说明 |
|------|------|------|
| 文化印象 | 20% | 常用字、正体字、字义丰富度 |
| 五行八字 | 25% | 名字五行与喜用神匹配度 |
| 生肖 | 10% | 生肖与名字五行生克关系 |
| 五格数理 | 25% | 天/人/地/外/总格吉凶 |
| 音韵 | 20% | 姓名音韵和谐度 |

---

## 技术栈

- **Go 1.22+** — 主语言
- **Ent ORM** — 数据库 ORM
- **SQLite3** — 数据存储
- **chronos/v2** — 八字计算
- **yi** — 周易卦象

---

## License

MIT License
