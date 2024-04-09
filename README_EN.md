# 命运起名 (Fate)

![FATE](docs/fate.png)

![Go](https://github.com/babyname/fate/workflows/Go/badge.svg)
[![GoDoc](https://godoc.org/github.com/babyname/fate?status.svg)](http://godoc.org/github.com/babyname/fate)
[![license](https://img.shields.io/github/license/babyname/fate.svg)](https://github.com/babyname/fate/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/babyname/fate)](https://goreportcard.com/report/github.com/babyname/fate)

[[ENGLISH](README_EN.md)][[中文](README.md)]

## A Modern Science-based Chinese Name Creation Tool

The first open-source Chinese name creation project on `Github`.

## Table of Contents

[TOR]

### Introduction

This program is suitable for generating two-character names with either a single or double surname (e.g., 独孤**, 李张**, 张**, 王**).
A good name accompanies one throughout their life; Fate helps you choose an auspicious and meaningful name.

### Versions

Specific versions are released separately, with pre-release binary files generated for each commit going forward.
The latest version uses Sqlite3 database, eliminating the need to import database files manually.
Simply download the Sqlite3 database below and use it locally.

- **[v3.5.5 Download](https://github.com/babyname/fate/releases/tag/v3.5.5)**
- **[Sqlite3 Database](https://github.com/babyname/fate/releases/download/v3.5.4/fate_sqlite3_database.zip)**
- **[Latest Self-Compiled Version](https://github.com/babyname/fate/releases/tag/auto_build)**
- **[Older SQL Database File: 20200331](https://github.com/babyname/fate/releases/download/v3.5.1/fate_db_200331.7z)**

### Usage Methods

#### Generating Names via Go Code and API Calls

```go
// Import the database data before using (test character library is mostly complete; for peace of mind, you can verify generated names on some naming websites later).
// Load configuration (refer to example/create_a_name for specific parameters).
cfg := config.Default()
// Birthdate
born := chronos.New("2020/01/23 11:31")
// Surname: lastName := "张"
// First parameter: surname
// Second parameter: birthdate
f := fate.NewFate(lastName, born.Solar().Time(), fate.ConfigOption(cfg))
e := f.MakeName(context.Background())
if e != nil {
t.Fatal(e)
}

}
```

#### Generating Names Using Precompiled Binary Files

```shell
# 生成配置文件，可修改数据库，及一些基本参数
fate.exe init
# 输出姓名
fate.exe name -l 张 -b "2020/02/06 15:04"
```

### Common Issues

#### Error: count total error: The system cannot find the path specified

- Time conversion failure due to missing zoneinfo (usually occurs in Windows environments). Download the zoneinfo file above and place it in the same directory as the executable file to resolve the issue.
- The latest version no longer requires the zoneinfo file.

#### How to Import Data (MySQL)

1. Connect to the MySQL database: mysql -u username -p password
2. Create the database: CREATE SCHEMA fate DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
3. Use the fate database: USE fate;
4. Import the database file: SOURCE /path/to/sql/file;
5. It is recommended to use tools like Navicat for faster import speeds.

#### Database Configuration: Replace Relevant Parts in config.json

**MYSQL Configuration**

```json
"database": { "host": "127.0.0.1", "port": "3306", "user": "root", "pwd": "111111", "name": "fate", "max_idle_con": 0, "max_open_con": 0, "driver": "mysql", "file": "", "dsn": "", "show_sql": false, "show_exec_time": false },
```

**SQLITE3 Configuration**

```json
"database": { "name": "fate", "driver": "sqlite3", },
```

## Version Plans

### Version 1

Primarily manual work; now deprecated.

### Version 2

Automatically generates name characters with manual screening.

### Version 3 (In Development)

1. Adds API query functionality (may require web-related work in the future; interested parties can sign up).
2. Refines and simplifies the dictionary, migrating it to `Entgo`.
3. Defines more comprehensive query rules.

### Version 4 (Planned)

Optimizes algorithms, adjusts interfaces, databases, enhances documentation, and fixes bugs.

### Version 7 (Planned)

Utilizes AI and big data matching algorithms to generate even better names.

### About Fate

Fate employs the following algorithms to automatically generate names based on dictionary rules:

- I Ching Hexagrams
- Da Yan's Numbers
- Three Talents and Five Elements
- Auspicious Deities (Balancing Deities)
- Zodiac Characters
- Eight Characters' Auspiciousness

Currently, Fate calculates names using a comprehensive approach based on six schools:

- Stroke Count School: Believes that all auspicious strokes ensure a prosperous life. Accuracy: 12.5%.
- Three Talents School: Completely disregards stroke auspiciousness, believing that harmony among the five elements associated with Heaven, Earth, and Man ensures a prosperous life. Accuracy: 56.6%.
- Augmenting Eight Characters: Disregards stroke auspiciousness, asserting that names compensating for deficiencies in one's innate eight-character birth chart guarantee prosperity. In reality, this method's accuracy is very low.
- Hexagram Image School: Ignores stroke auspiciousness, considering only the beauty of the hexagram derived from the name as indicative of a fortunate life. Accuracy: 40.26%.
- Heavenly Fortune School: Pays no attention to stroke auspiciousness, maintaining that avoiding being subdued by the five elements of the year of birth in one's name leads to a prosperous life. Accuracy: 25.32%.
- Zodiac Sign School: Disregards stroke auspiciousness, asserting that using the correct character shapes associated with one's zodiac sign ensures a fortunate life. Accuracy: 27.55%.

Currently used libraries:

- BaZi calculation (for determining birth charts)::
  <https://github.com/godcong/chronos>
- Dictionary data (a crawler tool populating the dictionary database):
  <https://github.com/godcong/excavator>
  If you know of any better alternatives, please let us know.

### 资料查询

1. Compilation of websites for nationwide and provincial name duplication checks

Provided by a netizen：`https://zhuanlan.zhihu.com/p/89654568` (**Please exercise caution when accessing non-official sites**)
[Repository address](./docs/chinese_name_query.md)

### Contributors
