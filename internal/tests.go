package internal

import (
	"fmt"

	"github.com/mannemsolutions/pgsectest/pkg/pg"
)

type Tests []Test

type Test struct {
	Name     string     `yaml:"name"`
	Dividend string     `yaml:"dividend"`
	Divisor  string     `yaml:"divisor"`
	Url      string     `yaml:"url"`
	Advice   string     `yaml:"advice"`
	Results  pg.Results `yaml:"results"`
	Score    TestScore  `yaml:"score"`
}

func (t *Test) Validate() (err error) {
	if t.Name == "" {
		return fmt.Errorf("a defined test is missing the query and name arguments")
	}
	return t.Score.Validate()
}
