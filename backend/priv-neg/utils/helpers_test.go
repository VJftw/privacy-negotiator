package utils_test

import (
	"testing"

	"reflect"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
)

func TestArrayUnion(t *testing.T) {
	type unionTest struct {
		A   []string
		B   []string
		res []string
	}

	spec := []unionTest{
		{
			A:   []string{"a", "b", "c"},
			B:   []string{"a", "b"},
			res: []string{"a", "b"},
		},
		{
			A:   []string{"a", "b", "c"},
			B:   []string{"d", "e"},
			res: []string{},
		},
	}

	for _, s := range spec {
		actualRes := utils.ArrayUnion(s.A, s.B)

		if !reflect.DeepEqual(actualRes, s.res) {
			t.Fail()
		}
	}

}

func TestIsSubset(t *testing.T) {

	type subsetTest struct {
		A   []string
		B   []string
		res bool
	}

	spec := []subsetTest{
		{
			A:   []string{"a"},
			B:   []string{"a", "b"},
			res: true,
		},
		{
			A:   []string{"a", "b"},
			B:   []string{"a", "b", "c"},
			res: true,
		},
		{
			A:   []string{"a", "b", "c"},
			B:   []string{"a", "b"},
			res: false,
		},
		{
			A:   []string{"a", "b", "c"},
			B:   []string{"a", "b", "d"},
			res: false,
		},
		{
			A:   []string{"a", "b", "c"},
			B:   []string{"a", "c", "b"},
			res: true,
		},
		{
			A:   []string{"a", "b", "c"},
			B:   []string{"b", "c", "a", "d"},
			res: true,
		},
		{
			A:   []string{"a", "b", "e"},
			B:   []string{"b", "c", "a", "d"},
			res: false,
		},
	}

	for _, s := range spec {
		actualRes := utils.IsSubset(s.A, s.B)

		if actualRes != s.res {
			t.Fail()
		}
	}
}
