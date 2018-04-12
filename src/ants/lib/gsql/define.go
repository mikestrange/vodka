package gsql

import (
	"ants/gutil"
	"database/sql"
	"fmt"
)

const MysqlMinOpen = 5
const MysqlMaxOpen = 20

//这里可以添加更多字段, json, byte[]等等
func Int(v interface{}) int {
	switch r := v.(type) {
	case int:
		return r
	case int64:
		return int(r)
	case []byte:
		return gutil.Atoi(string(r))
	case string:
		return gutil.Atoi(r)
	}
	return 0
}

func Int64(v interface{}) int64 {
	switch r := v.(type) {
	case int64:
		return r
	case int:
		return int64(r)
	case []byte:
		return gutil.Atol(string(r))
	case string:
		return gutil.Atol(r)
	}
	return 0
}

func Str(v interface{}) string {
	switch r := v.(type) {
	case string:
		return r
	case []byte:
		return string(r)
	case int:
		return gutil.Itoa(r)
	case int64:
		return gutil.Ltoa(r)
	}
	return ""
}

//static
//方案1(直接返回结果)
func toResult(rows *sql.Rows) []map[string]interface{} {
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("result Err:", err)
		return []map[string]interface{}{}
	}
	count := len(columns)
	var eles []map[string]interface{}
	vals := make([]string, count)
	ptrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			ptrs[i] = &vals[i]
		}
		rows.Scan(ptrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			val := vals[i]
			entry[col] = val
		}
		eles = append(eles, entry)
	}
	return eles
}

//方案2(封装了一下)
func toForm(rows *sql.Rows) IRowsForm {
	columns, err := rows.Columns()
	form := newForm()
	if err == nil {
		count := len(columns)
		vals := make([]string, count)
		ptrs := make([]interface{}, count)
		for rows.Next() {
			for i := 0; i < count; i++ {
				ptrs[i] = &vals[i]
			}
			rows.Scan(ptrs...)
			item := form.next()
			for i, col := range columns {
				item.set(col, vals[i])
			}
		}
	} else {
		fmt.Println("form Err:", err)
	}
	return form
}
