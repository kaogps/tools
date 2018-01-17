package sqlhandler

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/kdada/tinygo/sql"
)

type updateHandler struct {
	DB             *sql.DB
	condition      []string
	values         []interface{}
	isWhereAlready bool
	whereCondition []string
}

func (this *updateHandler) Set(value interface{}, key string) *updateHandler {
	if !this.isWhereAlready {
		this.condition = append(this.condition, key)
		this.values = append(this.values, value)
	}

	return this
}

func (this *updateHandler) Where(value interface{}, key string) *updateHandler {
	this.whereCondition = append(this.whereCondition, key)
	this.values = append(this.values, value)
	this.isWhereAlready = true
	return this
}

func (this *updateHandler) Exec(tableName string) (int, error) {
	if !this.isWhereAlready {
		return 0, errors.New("没有条件值！")
	}
	var sqlStr = "update " + tableName + " set "
	for _, v := range this.condition {
		sqlStr += " " + v + "=?,"
	}
	sqlStr = sqlStr[0:len(sqlStr)-1] + " where 1=1 "
	for _, v := range this.whereCondition {
		sqlStr += " and " + v + "= ? "
	}
	debug("sql语句：", sqlStr)
	return this.DB.Exec(sqlStr+" limit 1 ", this.values...).RowsAffected()
}

//// Update 单表修改 0值跳过
//func (this *updateHandler) Update(model interface{}, whereId string) bool {
//	var value = indirectModel(model)

//	var sql = "update `" + getDbName(value.Type().Name()) + "` set "
//	this.parseUpdate(value)
//	var columns string
//	if len(this.condition) > 0 {
//		for _, v := range this.condition {
//			columns += v[1] + ","
//		}
//	}
//	sql += columns[:len(columns)-1]
//	var id = value.FieldByName(whereId)
//	if id.IsValid() {
//		sql += " where " + whereId + "= " + strconv.Itoa(int(id.Int()))
//	} else {
//		fmt.Println(whereId, "条件字段错误！")
//		return false
//	}
//	var e = this.DB.Exec(sql, this.GetValues()...).Error()
//	if e != nil {
//		tools.Debug("sql错误:", e)
//		return false
//	}
//	return true

//}
//func (this *updateHandler) parseUpdate(t reflect.Value, whereId string) {
//	if t.Kind() != reflect.Struct {
//		return
//	}
//	for i := 0; i < t.NumField(); i++ {
//		if t.Type().Field(i).Anonymous {
//			this.parseUpdate(t.Field(i))
//			continue
//		}
//		var key = getDbName(t.Type().Field(i).Name)
//		if strings.ToLower(key) == whereId { //当为判断条件字段时，跳过
//			continue
//		}
//		var field = t.Field(i)
//		switch field.Kind() {
//		case reflect.Bool:
//		case reflect.Struct:
//			if field.Type().Name() == "Time" {
//				var tStr = fmt.Sprintf("%v", field.Interface())
//				var tm, err = time.Parse("2006-01-02", strings.Split(tStr, " ")[0])
//				if err != nil || tm.IsZero() {
//					continue
//				}
//				this.setCondtion(key, "`"+key+"`=?")
//				this.setValue(field.Interface())
//			}
//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//			if field.Int() != 0 {
//				this.setCondtion(key, "`"+key+"`=?")
//				this.setValue(field.Int())
//			}
//		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//			if field.Uint() != 0 {
//				this.setCondtion(key, "`"+key+"`=?")
//				this.setValue(field.Uint())
//			}
//		case reflect.Float32, reflect.Float64:
//			if field.Float() != 0 {
//				this.setCondtion(key, "`"+key+"`=?")
//				this.setValue(field.Float())
//			}
//		case reflect.Interface:
//			this.parseUpdate(field.Field(i))
//		case reflect.String:
//			if field.String() != "" {
//				this.setCondtion(key, "`"+key+"`=?")
//				this.setValue(field.String())
//			}
//		}
//	}
//}

//// UpdateModels 批量更新(0值字段跳过)
////func (this *BaseService) UpdateModels(models interface{}) {
////	//	var value = Indirect(models)
////	//	var columns = GetColumns(value)
////	//	_ = columns

////}

//// UpdateModels 批量更新(指定字段 字段名小写蛇形)

//func (this *updateHandler) UpdateWithTable(model interface{}, tableName, whereId string, columns ...string) {
//	// 判断column更新字段名是否为空
//	var colLen = len(columns)
//	if columns <= 0 {
//		return
//	}

//	// 设置条件和更新值

//}

//func (this *updateHandler) parse(t reflect.Value, columns []string, val map[string][]interface{}) {

//}

// UpdateColumnWithTable 批量更新(指定字段)
func (this *updateHandler) UpdateColumnWithTable(models interface{}, tableName string, whereId string, columns ...string) (int, error) {
	if len(columns) < 1 {
		return 0, errors.New("无指定字段")
	}
	var length = 0
	var value = indirectModel(models)
	if value.Kind() == reflect.Slice {
		length = value.Len()
		//tableName = "`" + SnakeName(value.Index(0).Type().Name()) + "`"
	} else if value.Kind() == reflect.Struct {
		length = 1
		//		tableName = "`" + SnakeName(value.Type().Name()) + "`"
	} else {
		return 0, errors.New("类型错误")
	}

	if length < 1 {
		return 0, errors.New("数组不能为空")
	}
	var sql = "update " + tableName + " SET "

	var setStr = " CASE " + whereId + " " + strings.Repeat(" WHEN ? THEN ? ", length) + " END "
	var colStr []string
	for i := 0; i < len(columns); i++ {
		var cc = "`" + columns[i] + "`=" + setStr
		colStr = append(colStr, cc)
	}
	sql += strings.Join(colStr, ",") + " WHERE " + whereId + " in (" +
		strings.Join(strings.Split(strings.Repeat("?", length), ""), ",") + ")"
	var vals = make(map[string][]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		vals[columns[i]] = nil
	}

	getValueMap(whereId, value, vals)
	var sqlVal []interface{}
	for _, v := range columns {
		for i := 0; i < length; i++ {
			var v1, v2 interface{}
			if vals[v] == nil {
				v2 = ""
			} else {
				v2 = vals[v][i]
			}
			if vals[whereId] == nil {
				v1 = nil
			} else {
				v1 = vals[whereId][i]
			}
			sqlVal = append(sqlVal, v1, v2)
			//sqlVal = append(sqlVal, vals[whereId][i], vals[v][i])
		}
	}
	fmt.Println(vals)
	//将where条件的id拼上
	sqlVal = append(sqlVal, vals[whereId]...)
	fmt.Println(sql)
	fmt.Println(sqlVal)
	return this.DB.Exec(sql, sqlVal...).RowsAffected()
}
func getValueMap(whereId string, t reflect.Value, vals map[string][]interface{}) {
	if t.Kind() == reflect.Slice {
		for i := 0; i < t.Len(); i++ {
			getValueMap(whereId, t.Index(i), vals)
		}
	}
	if t.Kind() == reflect.Struct {
		var n = t.NumField()
		for i := 0; i < n; i++ {
			var key = SnakeName(t.Type().Field(i).Name)
			if key != whereId {
				//检查key是否为指定字段
				if _, ok := vals[key]; !ok {
					continue
				}
			}
			if t.Type().Field(i).Anonymous {
				getValueMap(whereId, t.Field(i), vals)
			} else {
				vals[key] = append(vals[key], t.Field(i).Interface())
			}
		}
	}
}
