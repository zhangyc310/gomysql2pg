package cmd

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Database interface {
	TableCreate(logDir string, tableMap map[string][]string)
}

type Table struct {
	columnName             string
	dataType               string
	characterMaximumLength string
	isNullable             string
	columnDefault          string
	numericPrecision       string
	numericScale           string
	datetimePrecision      string
	columnKey              string
	columnComment          string
	ordinalPosition        int
	destType               string
	destNullable           string
	destDefault            string
}

func (tb *Table) TableCreate(logDir string, tableMap map[string][]string) {
	// 声明一个等待组
	var wg sync.WaitGroup
	tableCount := 0
	// 获取tableMap键值对中的表名
	for tblName, _ := range tableMap {
		tableCount += 1
		// 每一个任务开始时, 将等待组增加1
		wg.Add(1)
		// 使用goroutine并发的创建多个表
		go func(tblName string, tb Table, tableCount int) {
			var colTotal int
			// 使用defer, 表示函数完成时将等待组值减1
			defer wg.Done()
			pgCreateTbl := "create table " + tblName + "("
			// 查询当前表总共有多少个列字段
			colTotalSql := fmt.Sprintf("select count(*) from information_schema.COLUMNS  where table_schema=database() and table_name='%s'", tblName)
			err := srcDb.QueryRow(colTotalSql).Scan(&colTotal)
			if err != nil {
				log.Error(err)
			}
			// 查询MySQL表结构
			sql := fmt.Sprintf("select concat('\"',lower(column_name),'\"'),data_type,ifnull(character_maximum_length,'null'),is_nullable,case  column_default when '( \\'user\\' )' then 'user' else ifnull(column_default,'null') end as column_default,ifnull(numeric_precision,'null'),ifnull(numeric_scale,'null'),ifnull(datetime_precision,'null'),ifnull(column_key,'null'),ifnull(column_comment,'null'),ORDINAL_POSITION from information_schema.COLUMNS where table_schema=database() and table_name='%s'", tblName)
			//fmt.Println(sql)
			rows, err := srcDb.Query(sql)
			if err != nil {
				log.Error(err)
			}
			// 遍历MySQL表字段,一行就是一个字段的基本信息
			for rows.Next() {
				if err := rows.Scan(&tb.columnName, &tb.dataType, &tb.characterMaximumLength, &tb.isNullable, &tb.columnDefault, &tb.numericPrecision, &tb.numericScale, &tb.datetimePrecision, &tb.columnKey, &tb.columnComment, &tb.ordinalPosition); err != nil {
					log.Error(err)
				}
				//fmt.Println(columnName,dataType,characterMaximumLength,isNullable,columnDefault,numericPrecision,numericScale,datetimePrecision,columnKey,columnComment,ordinalPosition)
				//适配MySQL字段类型到PostgreSQL字段类型
				// 列字段是否允许null
				switch tb.isNullable {
				case "NO":
					tb.destNullable = "not null"
				default:
					tb.destNullable = "null"
				}
				// 列字段default默认值的处理
				switch {
				case tb.columnDefault != "null": // 默认值不是null并且是字符串类型下面就需要使用fmt.Sprintf格式化让字符串单引号包围，否则这个字符串是没有引号包围的
					if tb.dataType == "varchar" {
						tb.destDefault = fmt.Sprintf("default '%s'", tb.columnDefault)
					} else if tb.dataType == "char" {
						tb.destDefault = fmt.Sprintf("default '%s'", tb.columnDefault)
					} else {
						tb.destDefault = fmt.Sprintf("default %s", tb.columnDefault) // 非字符串类型无需使用单引号包围
					}
				default:
					tb.destDefault = "" // 如果没有默认值，默认值就是空字符串，即目标没有默认值
				}
				// 列字段类型的处理
				switch tb.dataType {
				case "int", "mediumint", "tinyint":
					tb.destType = "int"
				case "varchar":
					tb.destType = "varchar(" + tb.characterMaximumLength + ")"
				case "char":
					tb.destType = "char(" + tb.characterMaximumLength + ")"
				case "text", "tinytext", "mediumtext", "longtext":
					tb.destType = "text"
				case "datetime", "timestamp":
					tb.destType = "timestamp"
				case "decimal", "double", "float":
					if tb.numericScale == "null" {
						tb.destType = "decimal(" + tb.numericPrecision + ")"
					} else {
						tb.destType = "decimal(" + tb.numericPrecision + "," + tb.numericScale + ")"
					}
				case "tinyblob", "blob", "mediumblob", "longblob":
					tb.destType = "bytea"
				// 其余类型，源库使用什么类型，目标库就使用什么类型
				default:
					tb.destType = tb.dataType
				}
				// 在目标库创建的语句
				pgCreateTbl += fmt.Sprintf("%s %s %s %s,", tb.columnName, tb.destType, tb.destNullable, tb.destDefault)
				if tb.ordinalPosition == colTotal {
					pgCreateTbl = pgCreateTbl[:len(pgCreateTbl)-1] + ")" // 最后一个列字段结尾去掉逗号,并且加上语句的右括号
				}
			}
			//fmt.Println(pgCreateTbl) // 打印创建表语句
			// 创建前先删除目标表
			dropDestTbl := "drop table if exists " + tblName + " cascade"
			if _, err = destDb.Exec(dropDestTbl); err != nil {
				log.Error(err)
			}
			// 创建PostgreSQL表结构
			log.Info(fmt.Sprintf("%v ProcessingID %s create table %s", time.Now().Format("2006-01-02 15:04:05.000000"), strconv.Itoa(tableCount), tblName))
			if _, err = destDb.Exec(pgCreateTbl); err != nil {
				log.Error("table ", tblName, " create failed ", err)
				LogError(logDir, "tableCreateFailed", pgCreateTbl, err)
			}
		}(tblName, *tb, tableCount)
	}
	// 等待所有的任务完成
	wg.Wait()
	log.Info("Table structure synced from MySQL to PostgreSQL Table count ", tableCount)
	// 如果指定-t选项，表创建完毕之后就退出程序
	if tableOnly {
		os.Exit(0)
	}
}
