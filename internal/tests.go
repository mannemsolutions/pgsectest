package internal

import (
	"fmt"

	"github.com/mannemsolutions/pgsectest/pkg/pg"
)

type Tests []Test

type Test struct {
	Name    string     `yaml:"name"`
	Query   string     `yaml:"query"`
	Results pg.Results `yaml:"results"`
	Score   TestScore  `yaml:"score"`
}

func (t *Test) Validate() (err error) {
	if t.Name == "" {
		t.Name = t.Query
	} else if t.Query == "" {
		// Let's hope it is a query
		t.Query = t.Name
	}
	if t.Name == "" {
		return fmt.Errorf("a defined test is missing the query and name arguments")
	}
	return t.Score.Validate()
}
