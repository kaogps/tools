package sqlhandler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kdada/tinygo/sql"
)

type insertUpdateSqlHandler struct {
	DB *sql.DB
}

// InsertOrUpdateAdd 插入或者更新增加某些字段
func (this *insertUpdateSqlHandler) InsertOrUpdate(models interface{}, columns ...string) (int, bool) {
	var value = reflect.ValueOf(models)
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return this.insertOrUpdateModels(value, columns...)
}

// insertOrUpdateModels 新增数组
func (this *insertUpdateSqlHandler) insertOrUpdateModels(models reflect.Value, cols ...string) (int, bool) {
	if models.Kind() == reflect.Ptr {
		models = models.Elem()
	}
	if models.Len() < 1 {
		fmt.Println("结构体数组为空！")
		return 0, false
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
		fmt.Println("表名获取失败！")
		return 0, false
	}
	var sql = "insert into `" + SnakeName(tableName) + "`(" + strings.Join(columns, ",") + ")values " + strings.Join(vs, ",")

	if len(cols) > 0 {
		sql += "  on duplicate key update"
		var updateStrs = make([]string, len(cols))
		for i := range cols {
			updateStrs[i] = " " + cols[i] + "=values(" + cols[i] + ")"
		}
		sql += strings.Join(updateStrs, ",")
	}
	var res = this.DB.Exec(sql, values...)
	fmt.Println("[tinysql]  ", sql, values)
	if res.Error() != nil {
		fmt.Println(res.Error())
		return 0, false
	}
	var id, err1 = res.LastInsertId()
	if err1 != nil {
		fmt.Println(err1)
		return 0, false
	}
	return id, true
}
