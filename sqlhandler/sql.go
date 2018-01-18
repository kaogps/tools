package sqlhandler

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/kdada/tinygo/sql"
)

type SqlHandler struct {
	DB           *sql.DB
	where        *whereHandler
	insert       *insertHandler
	update       *updateHandler
	insertUpdate *insertUpdateSqlHandler
}

func (this *SqlHandler) Open(dbName string) {
	//fmt.Println(this.DB.Exec("use " + dbName).Error())
	this.DB = sql.Open(dbName)
	this.DB.Exec("use " + dbName)
}

// Where 创建查询处理器
func (this *SqlHandler) Where() *whereHandler {
	if this.where == nil {
		this.where = &whereHandler{DB: this.DB}
	}
	return this.where

}

// Insert 创建插入处理器
func (this *SqlHandler) Insert() *insertHandler {
	if this.insert == nil {
		this.insert = &insertHandler{DB: this.DB}
	}
	return this.insert
}

// Updata 创建更新处理器Updata
func (this *SqlHandler) Update() *updateHandler {
	if this.update == nil {
		this.update = &updateHandler{DB: this.DB}
	}
	return this.update
}

// InsertUpdate 创建插入更新处理器
func (this *SqlHandler) InsertUpdate() *insertUpdateSqlHandler {
	if this.insertUpdate == nil {
		this.insertUpdate = &insertUpdateSqlHandler{DB: this.DB}
	}
	return this.insertUpdate
}

// Indirect 初始化model
//func indirectModel(models interface{}) reflect.Value {
//	var value = reflect.ValueOf(models)
//	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
//		return value.Elem()
//	}
//	return value

//}

// Indirect 初始化model
func indirectModel(models interface{}) reflect.Value {
	var value = reflect.ValueOf(models)
	return indirect(value)
}
func indirect(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		return indirect(value.Elem())
	}
	return value
}

// Debug 输出调试信息,包含调用位置
// 如果当前处于发布模式,则不输出logs
func debug(logs ...interface{}) {

	var info = "[DEBUG] " + outputLineInfo()
	var allInfo = info + fmt.Sprintln(logs...)
	log.Print(allInfo)

}

// outputLineInfo 生成行信息
func outputLineInfo() string {
	var _, file, line, _ = runtime.Caller(2)
	var _, fileName = filepath.Split(file)
	return fmt.Sprint(fileName, ":", line, " ")
}

// GetColumns 获取字段名(匿名字段递归解析)
func GetColumns(t reflect.Value) []string {
	var res []string
	if t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		res = append(res, GetColumns(t.Index(0))...)
	}
	if t.Kind() == reflect.Struct {
		var n = t.NumField()
		for i := 0; i < n; i++ {
			if t.Type().Field(i).Anonymous {
				res = append(res, GetColumns(t.Field(i))...)
				continue
			} else {
				res = append(res, "`"+SnakeName(t.Type().Field(i).Name)+"`")
			}
		}
	}
	return res
}

// GetValues 获取要插入的值
func GetValues(t reflect.Value) []interface{} {
	var vals []interface{}
	if t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		for i := 0; i < t.Len(); i++ {
			vals = append(vals, GetValues(t.Index(i))...)
		}
	}
	if t.Kind() == reflect.Struct {
		var n = t.NumField()
		for i := 0; i < n; i++ {
			if t.Type().Field(i).Anonymous {
				vals = append(vals, GetValues(t.Field(i))...)
			} else {
				vals = append(vals, t.Field(i).Interface())
			}
		}
	}
	return vals
}

// SnakeName 驼峰转蛇形
func SnakeName(base string) string {
	var r = make([]rune, 0, len(base))
	var b = []rune(base)
	for i := 0; i < len(b); i++ {
		if i > 0 && b[i] >= 'A' && b[i] <= 'Z' {
			r = append(r, '_', b[i]+32)
			continue
		}
		if i == 0 && b[i] >= 'A' && b[i] <= 'Z' {
			r = append(r, b[i]+32)
			continue
		}
		r = append(r, b[i])
	}
	return string(r)
}
