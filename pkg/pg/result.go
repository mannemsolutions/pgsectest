package pg

import (
	"fmt"
	"regexp"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map"
)

type Result struct {
	kv *orderedmap.OrderedMap
}

func NewResultFromByteArrayArray(cols []string, values []interface{}) (r Result, err error) {
	om := orderedmap.New()
	if len(cols) != len(values) {
		return r, fmt.Errorf("number of cols different then number of values")
	}
	for i, col := range cols {
		om.Set(col, NewResultValue(values[i]))
	}
	r.kv = om
	return r, nil
}

func (r Result) OneField() (rv ResultValue, err error) {
	if oldest := r.kv.Oldest(); oldest == nil {
		return rv, fmt.Errorf("There should be at least one column in the one row to get OneField")
	} else if v, ok := oldest.Value.(ResultValue); !ok {
		return rv, fmt.Errorf("Values should ResultValues")
	} else {
		return v, nil
	}
}

func (r Result) String() (s string) {
	var results []string
	for pair := r.kv.Oldest(); pair != nil; pair = pair.Next() {
		key := pair.Key.(string)
		value := pair.Value.(ResultValue)
		results = append(results, fmt.Sprintf("%s: %s",
			FormattedString(key),
			value.Formatted()))
	}
	return fmt.Sprintf("{ %s }", strings.Join(results, ", "))
}

func (r Result) Columns() (cols []string) {
	for pair := r.kv.Oldest(); pair != nil; pair = pair.Next() {
		key := pair.Key.(string)
		cols = append(cols, key)
	}
	return cols
}

func (r Result) Values() (vals []string) {
	for pair := r.kv.Oldest(); pair != nil; pair = pair.Next() {
		value := pair.Value.(ResultValue)
		vals = append(vals, value.AsString())
	}
	return vals
}

func (r Result) KeyValueStrings() (vals []string) {
	repl := regexp.MustCompile("([ ='])")
	for pair := r.kv.Oldest(); pair != nil; pair = pair.Next() {
		key := pair.Key.(string)
		key = repl.ReplaceAllString(key, "\\$1")
		value := pair.Value.(ResultValue)
		sVal := repl.ReplaceAllString(value.AsString(), "\\$1")
		vals = append(vals, fmt.Sprintf("'%s'='%s'", key, sVal))
	}
	return vals
}

func (r Result) Compare(other Result) (err error) {
	if r.kv.Len() != other.kv.Len() {
		return fmt.Errorf("number of columns different between row %v and compared row %v",
			r.Columns(), other.Columns())
	}
	for pair := r.kv.Oldest(); pair != nil; pair = pair.Next() {
		if key, ok := pair.Key.(string); !ok {
			return fmt.Errorf("my Keys should be strings")
		} else if value, ok := pair.Value.(ResultValue); !ok {
			return fmt.Errorf("my Values should ResultValues")
		} else if otherValue, exists := other.kv.Get(key); !exists {
			return fmt.Errorf("column row (%s) not in other row", FormattedString(key))
		} else if otherValue, ok := otherValue.(ResultValue); !ok {
			return fmt.Errorf("others Values should ResultValues")
		} else if matched, err := regexp.MatchString(otherValue.AsString(), value.AsString()); err != nil {
			if value != otherValue {
				return fmt.Errorf("comparedrow is not an re, and column %s differs between row (%s), and comparedrow (%s)",
					FormattedString(key),
					value.Formatted(),
					otherValue.Formatted())
			}
		} else if !matched {
			return fmt.Errorf("column %s value (%s) does not match with regular expression (%s)",
				FormattedString(key),
				value.Formatted(),
				otherValue.Formatted())
		}
	}
	return nil
}

type Results []Result

func (rs Results) RowsKeyValues() []string {
	var arr []string
	for _, result := range rs {
		arr = append(arr, strings.Join(result.KeyValueStrings(), ","))
	}
	return arr
}

func (rs Results) String() (s string) {
	if len(rs) == 0 {
		return "[ ]"
	}
	var arr []string
	for _, result := range rs {
		arr = append(arr, result.String())
	}
	return fmt.Sprintf("[ %s ]", strings.Join(arr, ", "))
}

func (rs Results) Compare(other Results) (err error) {
	if len(rs) != len(other) {
		return fmt.Errorf("different result (%s) then expected (%s)", rs.String(),
			other.String())
	}
	for i, result := range rs {
		err = result.Compare(other[i])
		if err != nil {
			return fmt.Errorf("different %d'th result: %s", i, err.Error())
		}
	}
	return nil
}

func (rs Results) OneField() (rv ResultValue, err error) {
	if len(rs) != 1 {
		return rv, fmt.Errorf("There should be exactly one row to get OneField")
	}
	return rs[0].OneField()
}
