package pg

import (
	"fmt"
	"regexp"
	"strings"
)

type Result map[string]ResultValue

func NewResultFromByteArrayArray(cols []string, values []interface{}) (r Result, err error) {
	r = make(Result)
	if len(cols) != len(values) {
		return r, fmt.Errorf("number of cols different then number of values")
	}
	for i, col := range cols {
		r[col] = NewResultValue(values[i])
	}
	return r, nil
}

func (r Result) OneField() (rv ResultValue, err error) {
	// if len(r) != 1 {
	// 	return fmt.Errorf("There should be exactly one row to get OneField")
	// } else
	cols := r.Columns()
	if len(cols) != 1 {
		return rv, fmt.Errorf("There should be exactly one column in the one row to get OneField")
	}
	return r[cols[0]], nil

}

func (r Result) String() (s string) {
	var results []string
	for key, value := range r {
		results = append(results, fmt.Sprintf("%s: %s",
			FormattedString(key),
			value.Formatted()))
	}
	return fmt.Sprintf("{ %s }", strings.Join(results, ", "))
}

func (r Result) Columns() (cols []string) {
	for key := range r {
		cols = append(cols, key)
	}
	return cols
}

func (r Result) Compare(other Result) (err error) {
	if len(r) != len(other) {
		return fmt.Errorf("number of columns different between row %v and compared row %v",
			r.Columns(), other.Columns())
	}
	for key, value := range r {
		otherValue, exists := other[key]
		if !exists {
			return fmt.Errorf("column row (%s) not in compared row", FormattedString(key))
		}
		if matched, err := regexp.MatchString(otherValue.AsString(), value.AsString()); err != nil {
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

func (rs Results) String() (s string) {
	var arr []string
	if len(rs) == 0 {
		return "[ ]"
	}
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
