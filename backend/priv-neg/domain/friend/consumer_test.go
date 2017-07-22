package friend

import (
	"log"
	"os"
	"testing"
)

func TestIsMaximalClique(t *testing.T) {

	rdFriends := []string{"CC", "MV", "DB", "MS", "JY", "RD"}
	graph := map[string][]string{
		"DB": {"RD", "MV", "CC", "DB"},
		"MV": {"RD", "DB", "CC", "MV"},
		"CC": {"RD", "DB", "MV", "CC"},
		"MS": {"RD", "RL", "JY", "JS", "MS"},
		"JY": {"RD", "MS", "RL", "JS", "JY"},
		"RL": {"MS", "JY", "JS", "RL"},
		"JS": {"MS", "JY", "RL", "JS"},
	}

	type isMaximalCliqueTest struct {
		NewClique []string
		res       bool
	}

	spec := []isMaximalCliqueTest{
		{
			NewClique: []string{"RD", "DB", "CC"},
			res:       false,
		},
		{
			NewClique: []string{"RD", "DB", "CC", "MV"},
			res:       true,
		},
	}

	consumer := Consumer{
		logger: log.New(os.Stdout, "[test] ", log.Lshortfile),
	}
	for _, s := range spec {
		actualRes := consumer.isMaximalClique(s.NewClique, rdFriends, graph)

		if actualRes != s.res {
			t.Fail()
		}
	}

}
