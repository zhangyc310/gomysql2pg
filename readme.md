# gomysql2pg

## 一、工具特性以及环境要求
### 1.1 功能特性

  支持MySQL数据库一键迁移到postgresql内核类型的目标数据库，如postgresql数据库、海量数据库vastbase、华为postgresql、电信telepg、人大金仓Kingbase V8R6等

- 无需繁琐部署，开箱即用，小巧轻量化

- 在线迁移MySQL到目标数据库的表、视图、索引、外键、自增列等对象

- 多个goroutine并发迁移数据，充分利用CPU多核性能

- 支持迁移源库部分表功能

- 记录迁移日志，转储表、视图等DDL对象创建失败的sql语句

- 一键迁移MySQL到postgresql，方便快捷，轻松使用


### 1.2 环境要求
在运行的客户端PC需要同时能连通源端MySQL数据库以及目标数据库

支持Windows、Linux、MacOS

### 1.3 如何安装
解压之后即可运行此工具

若在Linux环境下请使用unzip解压，例如：


`[root@localhost opt]# unzip gomysql2pg.zip`

## 二、使用方法

以下为Windows平台示例，其余操作系统命令行参数一样

`注意:`在`Windows`系统请在`CMD`运行此工具，如果是在`MacOS`或者`Linux`系统，请在有读写权限的目录运行

### 2.1 编辑yml配置文件

编辑`example.cfg`文件，分别输入源库跟目标数据库信息

```yaml
src:
  host: 192.168.1.3
  port: 3306
  database: test
  username: root
  password: 11111
dest:
  host: 192.168.1.200
  port: 5432
  database: test
  username: test
  password: 11111
pageSize: 100000
maxParallel: 30
tables:
  test1:
    - select * from test1
  test2:
    - select * from test2
exclude:
  operalog1
  operalog2
  operalog3

```


### 2.2 全库迁移

迁移全库表结构、行数据，视图、索引约束、自增列等对象

gomysql2pg.exe  --config 配置文件
```
示例
gomysql2pg.exe --config example.yml
```

### 2.3 查看迁移摘要

全库迁移完成之后会生成迁移摘要，观察下是否有失败的对象，通过查询迁移日志可对迁移失败的对象进行分析


### 2.4 其他迁移模式

除了迁移全库之外，工具还支持迁移部分数据库对象，如部分表结构，视图，自增列，索引等对象


#### 2.4.1 全库迁移

迁移全库表结构、行数据，视图、索引约束、自增列等对象

gomysql2pg.exe  --config 配置文件

```
示例
gomysql2pg.exe --config example.yml
```

#### 2.4.2 自定义SQL查询迁移

不迁移全库数据，只迁移部分表，根据配置文件中自定义查询语句迁移表结构和表数据到目标库
gomysql2pg.exe  --config 配置文件 -s

```
示例
gomysql2pg.exe  --config example.yml -s
```

#### 2.4.3 迁移全库所有表结构

仅在目标库创建所有表的表结构

gomysql2pg.exe  --config 配置文件 -t

```
示例
gomysql2pg.exe  --config example.yml -t
```

#### 2.4.4 迁移自定义表的表结构

仅在目标库创建自定义的表

gomysql2pg.exe  --config 配置文件 -s -t

```
示例
gomysql2pg.exe  --config example.yml -s -t
```


#### 2.4.5 迁移全库表数据

只迁移全库表行数据到目标库，仅行数据，不包括表结构

gomysql2pg.exe  --config 配置文件 onlyData
```
示例
gomysql2pg.exe  --config example.yml onlyData
```

#### 2.4.6 迁移自定义表数据

只迁移yml配置文件中自定义查询sql，仅行数据，不包括表结构

gomysql2pg.exe  --config 配置文件 onlyData -s

```
示例
gomysql2pg.exe  --config example.yml onlyData -s
```

#### 2.4.7 迁移自增列到目标序列形式

只迁移MySQL的自增列转换为目标数据库序列
gomysql2pg.exe  --config 配置文件 seqOnly

```
示例
gomysql2pg.exe  --config example.yml seqOnly
```

#### 2.4.8 迁移索引等约束

只迁移MySQL的主键、索引这类对象到目标数据库
gomysql2pg.exe  --config 配置文件 idxOnly

```
示例
gomysql2pg.exe  --config example.yml idxOnly
```

#### 2.4.9 迁移视图

只迁移MySQL的视图到目标数据库

gomysql2pg.exe  --config 配置文件 viewOnly

```
示例
gomysql2pg.exe  --config example.yml viewOnly
```

## change history

### v0.1.7
2023-07-11

使用多个goroutine并发创建表，迁移摘要信息优化

### v0.1.6
2023-07-10

Add Makefile and output config info


### v0.1.5
2023-07-07

增加全局变量通道处理迁移行数据失败的计数，会在迁移摘要中展示

### v0.1.4
2023-06-30

修复只能迁移linux pg库，在Windows下迁移失败的问题，创建表的方法目前改成了单线程

### v0.1.3
2023-06-28

增加单独迁移表行数据的命令，迁移摘要优化，错误信息转储到日志文件优化

### v0.1.2
2023-06-27

增加迁移摘要，完善创建有外键的约束

### v0.1.1
2023-06-26

增加创建视图、外键、触发器到目标数据库


### v0.1.0
2023-06-16

增加创建索引、主键、等约束

### v0.0.9
2023-06-14

新增创建序列


### v0.0.8
2023-06-13

使用多个goroutine并发生成每个表的迁移任务、创建表，其余优化

### v0.0.7
2023-06-12

修复prepareSqlStr中没有行数据被漏掉创建的表,迁移数据前会查询下目标表是否存在,其余优化

### v0.0.6
2023-06-09

增加创建基本表的功能

### v0.0.5
2023-06-06

增加标题字符图，显示版本信息,彩色文字显示输出

### v0.0.4
2023-06-05

在遇到Ctrl+c输入后主动关闭数据库正在运行的sql,输出格式简化,转储迁移失败的表数据到日志目录

### v0.0.3
2023-06-02

config文件增加端口设定,自定义sql外面包了一层select * from (自定义sql) where 1=0 用于获取列字段，避免查询全表数据,在copy方法的exec刷buffer之前，再一次主动使用row.close关闭数据库连接

### v0.0.2
2023-05-24

增加排除表参数，以及config yml文件配置异常检查

### v0.0.1
2023-05-23

log方法打印调用文件以及方法源行数，增加日志重定向到平面文件