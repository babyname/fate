# 命运(Fate)

![Go](https://github.com/godcong/fate/workflows/Go/badge.svg)
[![GoDoc](https://godoc.org/github.com/godcong/fate?status.svg)](http://godoc.org/github.com/godcong/fate)
[![license](https://img.shields.io/github/license/godcong/fate.svg)](https://github.com/godcong/fate/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/godcong/fate)](https://goreportcard.com/report/github.com/godcong/fate)

## 现代科学取名工具(A modern science chinese name create tool)

Github第一个开源的中文取名项目(The first chinese name create tool in github)

## 简介 ##

本程序适用于单个姓或双个姓，起2个名的情况。（如：独孤**，李张**，张**，王**）  
一个好名字伴随人的一生，FATE让你取一个好名字。

## 全国及各省重名查询网址汇总

原地址：`https://zhuanlan.zhihu.com/p/89654568`

[本仓库地址](./docs/chinese_name_query.md)

## 关于版本：

除非稳定版本会单独出release,以后每次提交都会生成二进制文件的pre_release提供下载.  
【[最新自编译版本](https://github.com/godcong/fate/releases/tag/auto_build)】
【[最新数据库文件:20200331](https://github.com/godcong/fate/releases/download/v3.5.1/fate_db_200331.7z)】
【[v3.5.2下载](https://github.com/godcong/fate/releases/tag/v3.5.2)】

## 关于起名算法 ##

FATE使用了以下算法,按照每种算法的准确度,使用程度也有高有低,不会一概而否,也不会偏向单独某种算法.

```
周易卦象  
大衍之数  
三才五格  
喜用神（平衡用神）  
生肖用字  
八字吉凶  
```  

## 接口调用生成姓名 ##

```   
      使用前请导入database的数据（测试字库已基本完善，保险起见生成姓名后可以去一些测名网站验证下）
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

### 针对没有安装Go环境的用户,使用二进制文件在运行前务必把zoneinfo.zip下载并和二进制文件放在一起(不要解压),不然会报错.

### [zoneinfo文件](https://github.com/godcong/fate/blob/master/zoneinfo.zip)

## 二进制可执行文件生成姓名 ##

```   
       //没有安装go环境的请下载master下的zoneinfo文件和fate二进制文件放一起
       //生成配置文件(可修改数据库，及一些基本参数)：
       fate.exe init
       //输出姓名：
       fate.exe name -l 张 -b "2020/02/06 15:04"
```

## 常见问题:

```
1. Q: count total error:The system cannot find the path specified
   A: zoneinfo缺失导致的时间转换失败问题(一般发生在windows环境下),
        下载上面的zoneinfo文件并放到执行文件相同的目录下即可解决.
        最新版会检查根目录,已无需重新init.
        地址:https://github.com/godcong/fate/blob/master/zoneinfo.zip

2. Q: 如何导入数据
   A: 
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

## 版本履历:

第一版:
大部分是手动工作,现已废弃

第二版:
可自动生成名字字符 + 手工筛选

第三版(开发中):
每次生成一个名字,并目标喜用神,生肖喜忌,月历转换,生成八字等功能  
八字计算: https://github.com/godcong/chronos  
字典数据: https://github.com/godcong/excavator  
数据库重新切回mysql,mongo虽然插入简单,检索语法太繁琐了...

第四版(计划中):
优化算法,调整接口,数据库,完善文档以及修复一些bug.

第五版(计划中):
图形界面UI,api接口调用.

第六版(计划中):
占坑备用

第七版(计划中):
通过AI,大数据匹配算法,取出更好更佳的名字.

Fate涵盖以下六大派作为起名的基准:
• 笔划派： 认为笔划全吉，人生就大吉。其实准确度仅12.5 %   
• 三才派： 完全不管笔划吉凶，只认为天地人三才五行相生，人生就大吉。其实准确度仅56.6 %。  
• 补八字： 完全不管笔划吉凶，只认为名字补到先天八字命盘欠缺，人生就大吉。其实准确度非常低。  
• 卦象派： 完全不管笔划吉凶，只认为名字求出卦象漂亮，人生就大吉。其实准确度仅40.26 %。  
• 天运派： 完全不管笔划吉凶，只认为名字不要被出生年天运五行所剋，人生就大吉。其实准确度仅25.32 %。  
• 生肖派： 完全不管笔划吉凶，只认为生肖用对字形，人生就大吉。其实准确度仅27.55 %。

## PS ##

在过去的几年中虽然Fate经过了好几个版本的改进，但是仍然有许多不足之处。 其中包括生成的名字太多不容易筛选，  
还有些用户遇到了一些和Go相关的问题, 一些用户不知道如何导入数据库等。 这些问题都只能慢慢想办法去解决。
