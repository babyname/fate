# 命运起名(Fate)

![FATE](docs/fate.png)

![Go](https://github.com/babyname/fate/workflows/Go/badge.svg)
[![GoDoc](https://godoc.org/github.com/babyname/fate?status.svg)](http://godoc.org/github.com/babyname/fate)
[![license](https://img.shields.io/github/license/babyname/fate.svg)](https://github.com/babyname/fate/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/babyname/fate)](https://goreportcard.com/report/github.com/babyname/fate)

## 现代科学取名工具(A modern science chinese name create tool)

Github上第一个开源的中文取名项目(The first chinese name create tool in `github`)

## 目录

[TOR]

### 简介

本程序适用于单个姓或双个姓, 起2个名的情况. （如：独孤**, 李张**, 张**, 王**）  
一个好名字伴随人的一生, FATE让你取一个好名字.

### 关于版本

特定版本会单独出release,以后每次提交都会生成二进制文件的pre_release提供下载.  
最新版使用Sqlite3数据库,不在需要导入数据库文件了. 直接下载下面的Sqlite3数据库到本地就能使用.

【[v3.5.5下载](https://github.com/babyname/fate/releases/tag/v3.5.5)】
【[Sqlite3数据库](https://github.com/babyname/fate/releases/download/v3.5.4/fate_sqlite3_database.zip)】

【[最新自编译版本](https://github.com/babyname/fate/releases/tag/auto_build)】
【[旧版SQL数据库文件:20200331](https://github.com/babyname/fate/releases/download/v3.5.1/fate_db_200331.7z)】

### 使用方法

1. 编写运行go代码,接口调用生成姓名

    ```go
    //使用前请导入database的数据（测试字库已基本完善, 保险起见生成姓名后可以去一些测名网站验证下）
    //加载配置（具体参数参考example/create_a_name）
    cfg := config.Default()
    //生日：
    born := chronos.New("2020/01/23 11:31")
    //姓氏：
    lastName := "张"
    //第一参数：姓氏
    //第二参数：生日 
    f := fate.NewFate(lastName, born.Solar().Time(), fate.ConfigOption(cfg))
    
    e := f.MakeName(context.Background())
    if e != nil {
      t.Fatal(e)
    }
    ```

2. 使用预编译二进制文件生成姓名

    ```shell
    #没有安装go环境的请下载master下的zoneinfo文件和fate二进制文件放一起
    #生成配置文件, 可修改数据库, 及一些基本参数
    fate.exe init
    #输出姓名
    fate.exe name -l 张 -b "2020/02/06 15:04"
    ```

3. ~~针对没有安装Go环境的用户,使用二进制文件在运行前务必把zoneinfo.zip下载并和二进制文件放在一起(不要解压),不然会报错.~~
    ~~[zoneinfo文件](https://github.com/babyname/fate/blob/master/zoneinfo.zip)~~
    最新编译的版本使用了Go新版编译, 已经不再需要手动下载`zoneinfo.zip`文件了.

### 常见问题

1. 报错: count total error:The system cannot find the path specified

    ```docs
    zoneinfo缺失导致的时间转换失败问题(一般发生在windows环境下),
    下载上面的zoneinfo文件并放到执行文件相同的目录下即可解决.
    最新版会检查根目录,已无需重新init.
    地址:https://github.com/babyname/fate/blob/master/zoneinfo.zip
    ```

2. 如何导入数据(Mysql)

    ```docs
    //链接到mysql数据库
    mysql -u用户名 -p密码
    //创建数据库
    CREATE schema `fate` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
    //使用fate数据库
    use fate;
    //导入数据库文件
    source /path/to/sql/file;
    PS:建议使用Navicat等工具导入,导入速度较快
    ```

3. 数据库配置, 替换config.json中相关部分
   - MYSQL配置:
     - host填写mysql数据库的地址
     - user填写mysql数据库的用户名
     - pwd填写mysql数据库的密码
     - name填写mysql数据库的库名

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
    },   
   ```

   - SQLITE3配置:
     - name填写本地sqlite的数据库名字, 放在fate同一目录下

   ```json
    "database": {
      "name": "fate",
      "driver": "sqlite3",
    },
   ```

## 版本计划

第一版:
大部分是手动工作,现已废弃

第二版:
可自动生成名字字符 + 手工筛选

第三版(开发中):

1. 添加API接口查询(后期可能需要一些WEB方面的工作, 如果有兴趣的可以报名).
2. 完善精简字典库, 并迁移到EntORM.
3. 更完善的查询规则定义.

第四版(计划中):
优化算法,调整接口,数据库,完善文档以及修复一些bug.

第七版(计划中):
通过AI,大数据匹配算法,取出更好更佳的名字.

### 关于FATE

FATE使用了以下算法,查询字典库自动生成匹配规则的名字.
按照每种算法的准确度,使用程度也有高有低,不会一概而否,也不会偏向单独某种算法.
不会按照个人喜好做出选择.

- 周易卦象  
- 大衍之数  
- 三才五格  
- 喜用神（平衡用神）  
- 生肖用字  
- 八字吉凶  

目前Fate以六大派为基准综合计算生成名字:

- 笔划派: 认为笔划全吉, 人生就大吉. 准确度12.5 %
- 三才派: 完全不管笔划吉凶, 只认为天地人三才五行相生, 人生就大吉. 准确度56.6 %.
- 补八字: 完全不管笔划吉凶, 只认为名字补到先天八字命盘欠缺, 人生就大吉. 其实准确度非常低.
- 卦象派: 完全不管笔划吉凶, 只认为名字求出卦象漂亮, 人生就大吉. 准确度40.26 %.
- 天运派: 完全不管笔划吉凶, 只认为名字不要被出生年天运五行所剋, 人生就大吉. 准确度25.32 %.
- 生肖派: 完全不管笔划吉凶, 只认为生肖用对字形, 人生就大吉. 准确度27.55 %.

目前使用到的一些库:

- 八字计算(用于计算生辰): <https://github.com/godcong/chronos>  
- 字典数据(一个爬虫工具填充字典数据库): <https://github.com/godcong/excavator>
如果谁有更好用的可以告诉我.

### 资料查询

1. 全国及各省重名查询网址汇总

    网友提供：`https://zhuanlan.zhihu.com/p/89654568`(**请谨慎访问非本站点地址**)
    [本仓库地址](./docs/chinese_name_query.md)

## 一些废话

  在过去的几年中虽然Fate经过了好几个版本的改进, 但是仍然有许多不足之处.
  包括生成的名字太多不容易筛选,
  有些用户遇到了一些和Go相关的问题,
  一些用户不知道如何导入数据库等.
  这些问题都只能慢慢想办法去解决.

  还有些用户因为字典库生成的名字中有些字的寓意不好, 你可以手动删掉你不喜欢字, 却来恶意中伤作者.
  我想说的是这个字也不是我造的, 你如果有问题可以去找造那个字的人.
  如果觉得这个工具不好你可以不用.

  最近一年中因为作者个人原因导致Fate更新缓慢, 在这里向大家道个歉.
  大家也知道现在国内的IT环境, 毕竟我也要生活, 生活所迫没有太多时间放在业余的项目上.
  我只能尽量抽出时间来完善Fate的规则和代码. 
  在这里同样要感谢支持我的朋友们, 大家的出发点我相信是一样的.
  用这个工具目的都是为了给孩子取一个好名字.