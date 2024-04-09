# 命运起名 (Fate)

![FATE](docs/fate.png)

![Go Version](https://img.shields.io/badge/go%20version-%3E=1.22.1-blue.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/babyname/fate?status.svg)](http://godoc.org/github.com/babyname/fate)
[![license](https://img.shields.io/github/license/babyname/fate.svg)](https://github.com/babyname/fate/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/babyname/fate)](https://goreportcard.com/report/github.com/babyname/fate)

[[ENGLISH](README_EN.md)][[中文](README.md)]

## 现代科学取名工具 (A modern science Chinese name create tool)

Github 上第一个开源的中文取名项目 (The first Chinese name create tool on `github`)

## 目录

[TOR]

### 简介

本程序适用于单个姓或双个姓，起 2 个名的情况。（如：独孤**, 李张**, 张**, 王**）
一个好名字伴随人的一生，FATE 让你取一个好名字。

### 关于版本

特定版本会单独出 release，以后每次提交都会生成二进制文件的 pre_release 提供下载。
最新版使用 Sqlite3 数据库，不再需要导入数据库文件了。直接下载下面的 Sqlite3 数据库到本地就能使用。

- **[v3.5.5 下载](https://github.com/babyname/fate/releases/tag/v3.5.5)**
- **[Sqlite3 数据库](https://github.com/babyname/fate/releases/download/v3.5.4/fate_sqlite3_database.zip)**
- **[最新自编译版本](https://github.com/babyname/fate/releases/tag/auto_build)**
- **[旧版 SQL 数据库文件:20200331](https://github.com/babyname/fate/releases/download/v3.5.1/fate_db_200331.7z)**

### 使用方法

#### 编写运行 Go 代码，接口调用生成姓名

```go
// 使用前请导入 database 的数据（测试字库已基本完善，保险起见生成姓名后可以去一些测名网站验证下）
// 加载配置（具体参数参考 example/create_a_name）
cfg := config.Default()
// 生日
born := chronos.New("2020/01/23 11:31")
// 姓氏 lastName := "张"
// 第一参数：姓氏
// 第二参数：生日
f := fate.NewFate(lastName, born.Solar().Time(), fate.ConfigOption(cfg))
e := f.MakeName(context.Background())
if e != nil {
t.Fatal(e)
}
```

#### 使用预编译二进制文件生成姓名

```shell
# 生成配置文件，可修改数据库，及一些基本参数
fate.exe init
# 输出姓名
fate.exe name -l 张 -b "2020/02/06 15:04"
```

### 常见问题

#### 报错: count total error: The system cannot find the path specified

- zoneinfo 缺失导致的时间转换失败问题（一般发生在 Windows 环境下）， 下载上面的 zoneinfo 文件并放到执行文件相同的目录下即可解决。
- 最新版已不需要 zoneinfo 文件。

#### 如何导入数据 (MySQL)

1. 链接到 MySQL 数据库 mysql -u 用户名 -p 密码
2. 创建数据库 CREATE SCHEMA fate DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
3. 使用 fate 数据库 USE fate;
4. 导入数据库文件 SOURCE /path/to/sql/file;
5. 建议使用 Navicat 等工具导入，导入速度较快

#### 数据库配置，替换 config.json 中相关部分

**MYSQL 配置:**

```json
"database": {
  "host": "127.0.0.1",
  "port": "3306",
  "user": "root",
  "pwd": "111111",
  "name": "fate",
  "max_idle_con": 0,
  "max_open_con": 0,
  "driver": "mysql",
  "file": "",
  "dsn": "",
  "show_sql": false,
  "show_exec_time": false
}
```

**SQLITE3 配置:**

```json
"database": { "name": "fate", "driver": "sqlite3" }
```

## 版本计划

### 第一版

大部分是手动工作，现已废弃

### 第二版

可自动生成名字字符 + 手工筛选

### 第三版 (开发中)

1. 添加 API 接口查询（后期可能需要一些 WEB 方面的工作，如果有兴趣的可以报名）。
2. 完善精简字典库，并迁移到 EntORM。
3. 更完善的查询规则定义。

### 第四版 (计划中)

优化算法，调整接口，数据库，完善文档以及修复一些 bug。

### 第七版 (计划中)

通过 AI，大数据匹配算法，取出更好更佳的名字。

### 关于 FATE

FATE 使用了以下算法，查询字典库自动生成匹配规则的名字。
按照每种算法的准确度，使用程度也有高有低，不会一概而否，也不会偏向单独某种算法。
不会按照个人喜好做出选择。

- 周易卦象
- 大衍之数
- 三才五格
- 喜用神（平衡用神）
- 生肖用字
- 八字吉凶

目前 Fate 以六大派为基准综合计算生成名字:

- 笔划派: 认为笔划全吉，人生就大吉。准确度 12.5 %
- 三才派: 完全不管笔划吉凶，只认为天地人三才五行相生，人生就大吉。准确度 56.6 %。
- 补八字: 完全不管笔划吉凶，只认为名字补到先天八字命盘欠缺，人生就大吉。其实准确度非常低。
- 卦象派: 完全不管笔划吉凶，只认为名字求出卦象漂亮，人生就大吉。准确度 40.26 %。
- 天运派: 完全不管笔划吉凶，只认为名字不要被出生年天运五行所剋，人生就大吉。准确度 25.32 %。
- 生肖派: 完全不管笔划吉凶，只认为生肖用对字形，人生就大吉。准确度 27.55 %。

目前使用到的一些库:

- 八字计算（用于计算生辰）:
  <https://github.com/godcong/chronos>
- 字典数据（一个爬虫工具填充字典数据库）:
  <https://github.com/godcong/excavator>
  如果谁有更好用的可以告诉我。

### 资料查询

1. 全国及各省重名查询网址汇总

网友提供：`https://zhuanlan.zhihu.com/p/89654568` (**请谨慎访问非本站点地址**)
[本仓库地址](./docs/chinese_name_query.md)

### 贡献者

<table>
<tr>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/godcong>
            <img src=https://avatars.githubusercontent.com/u/2727298?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=godcong/>
            <br />
            <sub style="font-size:14px"><b>godcong</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/Z-fly>
            <img src=https://avatars.githubusercontent.com/u/10470892?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Z-fly/>
            <br />
            <sub style="font-size:14px"><b>Z-fly</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/fesiong>
            <img src=https://avatars.githubusercontent.com/u/9912496?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Sinclair/>
            <br />
            <sub style="font-size:14px"><b>Sinclair</b></sub>
        </a>
    </td>
</tr>
</table>
