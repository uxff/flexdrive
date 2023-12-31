package tplfuncs

import (
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
)

var tplFuncMap = make(template.FuncMap, 0)

func init() {
	loadFuncMap()
}

func GetFuncMap() template.FuncMap {
	return tplFuncMap
}

func loadFuncMap() {
	tplFuncMap["i18nja"] = func(format string, args ...interface{}) string {
		return "" //i18n.Tr("ja-JP", format, args...)
	}
	//"i18n": i18n.Tr,
	tplFuncMap["datenow"] = func(format string) string {
		return time.Now().Add(time.Duration(9) * time.Hour).Format(format)
	}
	tplFuncMap["dateformatJst"] = func(in time.Time) string {
		in = in.Add(time.Duration(9) * time.Hour)
		return in.Format("2006/01/02 15:04")
	}

	tplFuncMap["qescape"] = func(in string) string {
		return url.QueryEscape(in)
	}
	tplFuncMap["nl2br"] = func(in string) string {
		return strings.Replace(in, "\n", "<br>", -1)
	}

	tplFuncMap["tostr"] = func(in interface{}) string {
		return fmt.Sprintf("%d", in) //convert.ToStr(reflect.ValueOf(in).Interface())
	}

	tplFuncMap["first"] = func(in interface{}) interface{} {
		return reflect.ValueOf(in).Index(0).Interface()
	}

	tplFuncMap["last"] = func(in interface{}) interface{} {
		s := reflect.ValueOf(in)
		return s.Index(s.Len() - 1).Interface()
	}

	tplFuncMap["truncate"] = func(in string, length int) string {
		return runewidth.Truncate(in, length, "...")
	}

	tplFuncMap["noname"] = func(in string) string {
		if in == "" {
			return "(未入力)"
		}
		return in
	}

	tplFuncMap["cleanurl"] = func(in string) string {
		return strings.Trim(strings.Trim(in, " "), "　")
	}

	tplFuncMap["append"] = func(data map[interface{}]interface{}, key string, value interface{}) template.JS {
		if _, ok := data[key].([]interface{}); !ok {
			data[key] = []interface{}{value}
		} else {
			data[key] = append(data[key].([]interface{}), value)
		}
		return template.JS("")
	}

	tplFuncMap["appendmap"] = func(data map[interface{}]interface{}, key string, name string, value interface{}) template.JS {
		v := map[string]interface{}{name: value}

		if _, ok := data[key].([]interface{}); !ok {
			data[key] = []interface{}{v}
		} else {
			data[key] = append(data[key].([]interface{}), v)
		}
		return template.JS("")
	}
	tplFuncMap["urlfor"] = func(endpoint string, values ...interface{}) string {
		return endpoint
	}
	tplFuncMap["captchaUrl"] = func() string {
		return fmt.Sprintf("/captcha?t=%d", time.Now().Unix())
	}
	tplFuncMap["mgrStatus"] = func(status int) string {
		return base.StatusMap[status]
	}
	// 所有的空间单位必须是int64 基础为kB
	tplFuncMap["space4Human"] = func(space int64) string {
		if space < 1024 {
			return fmt.Sprintf("%d kB", space)
		}
		if space < 1024*1024 {
			return fmt.Sprintf("%.01f MB", float32(space)/1024)
		}
		if space < 1024*1024*1024 {
			return fmt.Sprintf("%.01f GB", float32(space)/1024/1024)
		}
		return fmt.Sprintf("%.02f TB", float32(space)/1024/1024/1024)
	}
	// 所有的大小单位必须是int64 基础为Byte
	tplFuncMap["size4Human"] = func(space int64) string {
		if space < 1024 {
			return fmt.Sprintf("%d B", space)
		}
		if space < 1024*1024 {
			return fmt.Sprintf("%.01f kB", float32(space)/1024)
		}
		if space < 1024*1024*1024 {
			return fmt.Sprintf("%.01f MB", float32(space)/1024/1024)
		}
		return fmt.Sprintf("%.02f GB", float32(space)/1024/1024/1024)
	}
	// 所有的空间单位必须是int64
	tplFuncMap["spaceRate"] = func(used int64, quota int64) string {
		return fmt.Sprintf("%d", int(float32(used)/float32(quota)*100))
	}
	tplFuncMap["orderStatus"] = func(orderStatus int) string {
		return dao.OrderStatusMap[orderStatus]
	}
	tplFuncMap["amount4Human"] = func(amount int) string {
		return fmt.Sprintf("%.02f", float32(amount)/100)
	}
	tplFuncMap["timeSmell"] = func(t time.Time) string {
		diff := time.Now().Sub(t) / time.Second
		if diff < 30 {
			return "#4edd1c" // green
		} else if diff < 120 {
			return "#dcdd1c" // yellow
		} else if diff < 300 {
			return "#ddc31c" // orange
		} else if diff < 3000 {
			return "#dd7c1c" // #dd1c1c // red
		} else {
			return "gray"
		}
	}
	tplFuncMap["orderStatus"] = func(status int) string {
		return dao.OrderStatusMap[status]
	}
}
