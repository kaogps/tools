package sqlhandler

import (
	"fmt"

	"reflect"
	"strings"

	"github.com/kdada/tinygo/sql"
)

type insertHandler struct {
	DB *sql.DB
}

// InsertModel 插入数据,表名即为model struct的名称
func (this *insertHandler) InsertModel(model interface{}) int {
	var v = reflect.TypeOf(model).Elem()
	var table = getDbName(v.Name())
	return this.Insert(table, model)
}

// Insert 向指定table插入数据
func (this *insertHandler) Insert(table string, model interface{}) int {
	query := "insert into `" + table + "`"
	value := reflect.ValueOf(model).Elem()
	data := make(map[string]interface{})
	mapStructToMap(value, data)
	keys := " ("
	values := " ("
	params := make([]interface{}, 0, 0)
	for k, v := range data {
		keys += "`" + k + "`,"
		values += "?,"
		params = append(params, v)
	}
	query += keys[:len(keys)-1] + ") values"
	query += values[:len(values)-1] + ")"
	var err error
	fmt.Println("[TinySql]", query)
	var result = this.DB.Exec(query, params...)
	var id int
	id, err = result.LastInsertId()
	if err != nil {
		fmt.Println("数据库错误：", err)
		return -1
	}
	return int(id)
}

// getDbName 驼峰转蛇形
func getDbName(str string) string {
	//	var reg = regexp.MustCompile(`/B[A-Z]`)
	//	return strings.ToLower(reg.ReplaceAllString(str, "_$0"))
	//	var reg = regexp.MustCompile(`([A-Z])+`)
	//	var res = reg.ReplaceAllString(str, "_$1")
	//	res = strings.ToLower(res)
	//	return string(res[1:])
	var r = make([]rune, 0, len(str))
	var b = []rune(str)
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

// mapStructToMap 将一个结构体所有字段(包括通过组合得来的字段)到一个map中
// value:结构体的反射值
// data:存储字段数据的map
func mapStructToMap(value reflect.Value, data map[string]interface{}) {
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			var fieldValue = value.Field(i)
			if fieldValue.CanInterface() {
				var fieldType = value.Type().Field(i)
				if fieldType.Anonymous {
					//匿名组合字段,进行递归解析
					mapStructToMap(fieldValue, data)
				} else {
					//非匿名字段
					var fieldName = fieldType.Tag.Get("db")
					if fieldName == "-" {
						continue
					}
					if fieldName == "" {
						fieldName = getDbName(fieldType.Name)
					}
					data[fieldName] = fieldValue.Interface()
					//fmt.Println(fieldName + ":" + fieldValue.Interface().(string))
				}
			}
		}
	}
}

// InsertBatch 批量插入
func (this *insertHandler) InsertBatch(models interface{}) bool {
	var value = reflect.ValueOf(models)
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return this.insertBatch(value)
}

// insertBatch 新增数组
func (this *insertHandler) insertBatch(models reflect.Value) bool {
	if models.Kind() == reflect.Ptr {
		models = models.Elem()
	}
	if models.Len() < 1 {
		return false
	}
	var data = models.Index(0)
	if data.Kind() == reflect.Interface {
		data = data.Elem()
	}
	var tableName = data.Type().Name()
	var columns = GetColumns(models.Index(0))
	var values = GetValues(models)
	var holders = make([]string, len(columns))
	for i := 0; i < len(columns); i++ {
		holders[i] = "?"
	}
	var vs = make([]string, models.Len())
	for i := 0; i < models.Len(); i++ {
		vs[i] = "(" + strings.Join(holders, ",") + ")"
	}
	if tableName == "" {
		return false
	}
	var sql = "insert into `" + SnakeName(tableName) + "`(" + strings.Join(columns, ",") + ")values " + strings.Join(vs, ",")

	var err = this.DB.Exec(sql, values...).Error()
	fmt.Println("[tinysql]  ", sql, values)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
