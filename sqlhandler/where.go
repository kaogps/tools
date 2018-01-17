package sqlhandler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/kdada/tinygo/sql"
)

type whereHandler struct {
	DB        *sql.DB
	condition [][2]string
	values    []interface{}
}

// NewQuery 新建查询 传入表名
func (this *whereHandler) NewQuery(table string) *whereHandler {
	return this.setCondtion("sql", "select * from `"+table+"`")
}

// Query 解析查询条件，结果赋予models
func (this *whereHandler) Query(models interface{}) (int, error) {
	return this.DB.Query(this.GetCondition(), this.GetValues()...).Scan(models)
}

// Limit limit条件
func (this *whereHandler) Limit(num ...int) *whereHandler {
	var limit string
	if len(num) == 1 {
		limit = ` limit ` + strconv.Itoa(num[0])
	} else if len(num) == 2 {
		limit = ` limit ` + strconv.Itoa(num[0]) + `,` + strconv.Itoa(num[1])
	}
	return this.AddCondition(limit)
}

// LimitPP limit分页，按page和pagesize
func (this *whereHandler) LimitPP(page, pagesize int) *whereHandler {
	var limit = ` limit ` + strconv.Itoa((page-1)*pagesize) + `,` + strconv.Itoa(pagesize)
	return this.AddCondition(limit)
}

// AddCondition 添加条件 无需key值,无value值
func (this *whereHandler) AddCondition(condition string, values ...interface{}) *whereHandler {
	//检查是否存在此key,有则添加
	for k, v := range this.condition {
		if v[0] == "condition" {
			this.condition[k][1] += " " + condition
			return this
		}
	}
	this.condition = append(this.condition, [2]string{"condition", condition})
	this.values = append(this.values, values...)
	return this
}

// Set 默认等于
func (this *whereHandler) Set(key string, value interface{}) *whereHandler {
	return this.SetEq(key, value)
}

// SetEq 等于
func (this *whereHandler) SetEq(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " and `"+key+"` = ? ")
	this.setValue(value)
	return this
}

// SetNotEq 不等于
func (this *whereHandler) SetNotEq(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " and `"+key+"` != ? ")
	this.setValue(value)
	return this
}

// SetNull 为空
func (this *whereHandler) SetNull(key string) *whereHandler {
	this.setCondtion(key, " and `"+key+"` is null ")
	return this
}

// SetNotNull 不为空
func (this *whereHandler) SetNotNull(key string) *whereHandler {
	this.setCondtion(key, " and `"+key+"` is not null ")
	return this
}

// SetGt 大于
func (this *whereHandler) SetGt(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " and `"+key+"` > ? ")
	this.setValue(value)
	return this
}

// SetGe 大于等于
func (this *whereHandler) SetGe(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " and `"+key+"` >=? ")
	this.setValue(value)
	return this
}

// SetLt 小于
func (this *whereHandler) SetLt(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " and `"+key+"` < ? ")
	this.setValue(value)
	return this
}

// SetLe 小于等于
func (this *whereHandler) SetLe(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " and `"+key+"` <= ? ")
	this.setValue(value)
	return this
}

// SetBetween 两者之间
func (this *whereHandler) SetBetween(key string, value1 interface{}, value2 interface{}) *whereHandler {
	this.setCondtion(key, " and `"+key+"` between ? and ? ")
	this.values = append(this.values, value1, value2)
	return this
}

// SetOr 或者等于
func (this *whereHandler) SetOr(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " or `"+key+"` = ? ")
	this.setValue(value)
	return this
}

// SetOrLike 或者相似
func (this *whereHandler) SetOrLike(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " or `"+key+"` like ? ")
	this.setValue(value)
	return this
}

// SetLike 相似
func (this *whereHandler) SetLike(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " and `"+key+"` like ?")
	this.setValue(value)
	return this
}

// SetDate 设置日期
func (this *whereHandler) SetDate(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " and to_days(`"+key+"`) = to_days(?) ")
	this.setValue(value)
	return this
}

// SetDateGe 设置起始日期
func (this *whereHandler) SetDateGe(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " and to_days(`"+key+"`) >= to_days(?) ")
	this.setValue(value)
	return this
}

// SetDateLe 设置结束日期
func (this *whereHandler) SetDateLe(key string, value interface{}) *whereHandler {
	this.setCondtion(key, " and to_days(`"+key+"`) <= to_days(?) ")
	this.setValue(value)
	return this
}

// SetIn In条件
func (this *whereHandler) SetIn(key string, value ...interface{}) *whereHandler {
	var length = len(value)
	if length == 0 {
		return this
	}
	var tmp = strings.Split(strings.Repeat("?", length), "")
	this.setCondtion(key, " and `"+key+"` in ("+strings.Join(tmp, ",")+")")
	this.setValue(value...)
	return this
}

// SetInSlice SetIn数组数据(可以为[]int,[]string等)
func (this *whereHandler) SetInSlice(key string, value interface{}) *whereHandler {
	//	var k = reflect.TypeOf(&value).Kind()
	//	fmt.Println("kind:", k)
	//	if k == reflect.Ptr {
	//		k = reflect.TypeOf(&value).Elem().Kind()
	//	}
	//	fmt.Println("kind:", k)
	//	if k != reflect.Slice {
	//		return this
	//	}
	var s []interface{}
	var d, _ = json.Marshal(value)
	var er = json.Unmarshal(d, &s)
	if er != nil {
		fmt.Println("SetInSlice---Error:", er)
		return this
	}
	return this.SetIn(key, s...)
}

// SetIn In条件
func (this *whereHandler) SetNotIn(key string, value ...interface{}) *whereHandler {
	var length = len(value)
	if length == 0 {
		return nil
	}
	var tmp = strings.Split(strings.Repeat("?", length), "")
	this.setCondtion(key, " and `"+key+"` not in ("+strings.Join(tmp, ",")+")")
	this.setValue(value...)
	return this
}

func (this *whereHandler) SetSql(sql string, value ...interface{}) *whereHandler {
	var length = len(value)
	if length == 0 {
		return nil
	}
	this.setCondtion(sql, " and "+sql+" ")
	this.setValue(value...)
	return this
}

// SetValue 直接添加value值
func (this *whereHandler) SetValue(value ...interface{}) *whereHandler {
	var length = len(value)
	if length == 0 {
		return nil
	}
	this.setValue(value...)
	return this
}

// GetCondition 返回所有条件
func (this *whereHandler) GetCondition() string {
	var sql string
	var condition string
	for _, v := range this.condition {
		if v[0] == "sql" {
			sql = v[1]
			continue
		}
		condition += " " + v[1]
	}
	if sql == "" {
		return condition
	}
	return sql + " where 1=1 " + condition
}

// GetWhere 返回where 1=1条件语句
func (this *whereHandler) GetWhere() string {
	var condition string
	for _, v := range this.condition {

		condition += v[1]
	}

	return " where 1=1 " + condition
}

// GetValues 返回值
func (this *whereHandler) GetValues() []interface{} {
	return this.values
}

// ParseQuery 将结构体解析到查询语句
func (this *whereHandler) ParseQuery(model interface{}, flag ...bool) *whereHandler {
	var value = indirectModel(model)
	if value.Kind() != reflect.Struct {
		//		return errors.New("参数类型错误")
	}
	//	this.setValue(GetValues(value)...)
	if len(flag) == 0 {
		this.parseQuery(value, true)
	} else {
		this.parseQuery(value, flag[0])
	}
	return this
}

// Example 列表请求参数结构体示例
// type Example struct {
// 	//key查询中对应数据库字段的名称,"_"表示忽略此值
// 	// int 默认等于; string 默认like; time 默认天数
// 	GoodsInfo   int       `key:"oi.off_product_id"`
// 	GoodsArea   int       `key:"gi.goods_area"`
// 	GoodsStatus int       `key:"oi.goods_status"`
// 	Page        int       `key:"_"`
// 	PageSize    int       `key:"_"`
// 	GoodsName   string    `key:"gi.name"`
// 	Date        time.Time `key:"o.predict_time"`
// }

func (this *whereHandler) parseQuery(t reflect.Value, flag bool) {
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		var tag = t.Type().Field(i).Tag.Get("key")
		if tag == "_" { //忽略此值
			continue
		}
		if t.Type().Field(i).Anonymous {
			this.parseQuery(t.Field(i), flag)
			continue
		}
		var key = getDbName(t.Type().Field(i).Name)
		if flag {
			if key == "page" || key == "pagesize" || key == "page_size" {
				continue
			}
		}
		if tag != "" {
			key = tag
		}
		var field = t.Field(i)
		switch field.Kind() {
		case reflect.Bool:
		case reflect.Struct:
			if field.Type().Name() == "Time" {
				var tStr = fmt.Sprintf("%v", field.Interface())
				var tm, err = time.Parse("2006-01-02", strings.Split(tStr, " ")[0])
				if err != nil || tm.IsZero() {
					continue
				}
				this.setCondtion(key, " and to_days("+key+")=to_days(?)")
				this.setValue(field.Interface())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				this.setCondtion(key, " and `"+key+"`=?")
				this.setValue(field.Int())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if field.Uint() != 0 {
				this.setCondtion(key, " and `"+key+"`=?")
				this.setValue(field.Uint())
			}
		case reflect.Float32, reflect.Float64:
			if field.Float() != 0 {
				this.setCondtion(key, " and `"+key+"`=?")
				this.setValue(field.Float())
			}
		case reflect.Interface:
			this.parseQuery(field.Field(i), flag)
		case reflect.String:
			if field.String() != "" {
				this.setCondtion(key, " and `"+key+"` like ? ")
				this.setValue("%" + field.String() + "%")
			}
		}
	}

}

func (this *whereHandler) setCondtion(key, condition string) *whereHandler {
	if strings.Contains(key, ".") {
		condition = strings.Replace(condition, "`", "", -1)
	}
	//检查是否存在此key,有则替换
	for k, v := range this.condition {
		if v[0] == key {
			this.condition[k][1] = condition
			return this
		}
	}
	this.condition = append(this.condition, [2]string{key, condition})
	return this
}
func (this *whereHandler) setValue(value ...interface{}) *whereHandler {
	this.values = append(this.values, value...)
	return this
}

// ResetCondition 重置condition
func (this *whereHandler) ResetCondition() *whereHandler {
	this.condition = make([][2]string, 0)
	return this
}

// ResetValue 重置value
func (this *whereHandler) ResetValue() *whereHandler {
	this.values = nil
	return this
}

// Reset 重置condition和value
func (this *whereHandler) Reset() *whereHandler {
	return this.ResetCondition().ResetValue()
}
