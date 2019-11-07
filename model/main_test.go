package model_test

import (
	"fmt"
	"github.com/Generalbelly/q-torial-api/model"
	"log"
	"testing"
)

var testCase = []struct {
	op string
	pv string
	up string
}{
	{
		op: model.PathEquals,
		pv: "test/",
		up: "test",
	},
	{
		op: model.PathEquals,
		pv: "test",
		up: "test/",
	},
	{
		op: model.PathStartsWith,
		pv: "te",
		up: "test",
	},
	{
		op: model.PathRegex,
		pv: "^tes",
		up: "test",
	},
	{
		op: model.PathAll,
		pv: "te1",
		up: "test",
	},
}

func TestValidateUrlPathValid(t *testing.T) {
	for _, c := range testCase {
		valid, err := model.ValidateUrlPath(c.op, c.pv, c.up)
		if err != nil {
			log.Fatal(err)
		}
		if valid == false {
			t.Fatal(fmt.Sprintf("op: %s, pv: %s, up: %s", c.op, c.pv, c.up))
		}
	}

}
