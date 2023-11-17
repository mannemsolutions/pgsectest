package pg

import (
	"fmt"
	"strings"
	"time"

	"reflect"

	"github.com/jackc/pgtype"
)

type ResultValue struct {
	org  interface{}
	str  *string
	fstr *string
	flt  *float64
}

func NewResultValue(org interface{}) ResultValue {
	return ResultValue{
		org: org,
	}
}

func (rv *ResultValue) Org() interface{} {
	return rv.org
}

func FormattedString(s string) string {
	return fmt.Sprintf("'%s'", strings.Replace(s, "'", "\\'", -1))
}

func (rv *ResultValue) Formatted() string {
	if rv.fstr == nil {
		fstr := FormattedString(rv.AsString())
		rv.fstr = &fstr
	}
	return *rv.fstr
}

func (rv *ResultValue) AsString() string {
	if rv.str == nil {
		var s string
		switch v := rv.org.(type) {
		case string:
			s = v
		case float32, float64:
			s = fmt.Sprintf("%f", v)
		case bool:
			s = fmt.Sprintf("%t", v)
		case time.Duration:
			s = v.String()
		case time.Time:
			s = v.String()
		case int:
			s = fmt.Sprintf("%d", v)
		case int8:
			s = fmt.Sprintf("%d", v)
		case int16:
			s = fmt.Sprintf("%d", v)
		case int32:
			s = fmt.Sprintf("%d", v)
		case int64:
			s = fmt.Sprintf("%d", v)
		case uint:
			s = fmt.Sprintf("%d", v)
		case uint8:
			s = fmt.Sprintf("%d", v)
		case uint16:
			s = fmt.Sprintf("%d", v)
		case uint32:
			s = fmt.Sprintf("%d", v)
		case uint64:
			s = fmt.Sprintf("%d", v)
		case []byte:
			s = fmt.Sprintf("%d", v)
		case pgtype.Float4Array:
			var l []string
			for _, e := range v.Elements {
				l = append(l, fmt.Sprintf("%f", e.Float))
			}
			s = fmt.Sprintf("[%s]", strings.Join(l, ","))
		case nil:
			s = "nil"
		default:
			s = fmt.Sprintf("unknown datatype %v", v)
		}
		rv.str = &s
	}
	return *rv.str
}

var floatType = reflect.TypeOf(float64(0))

func (rv *ResultValue) AsFloat() (float64, error) {
	if rv.flt == nil {
		v := reflect.ValueOf(rv.org)
		v = reflect.Indirect(v)
		if !v.Type().ConvertibleTo(floatType) {
			return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
		}
		fv := v.Convert(floatType)
		f := fv.Float()
		rv.flt = &f
	}
	return *rv.flt, nil
}
