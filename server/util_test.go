package server

import (
	"fmt"
	"testing"
)

func TestIsValidQuery(t *testing.T) {
	querys := []string{
		`C1==test1`,
		`C1=="test1"`,
		`C2 != "k1"`,
		`C3 &= "K2"`,
		`C4 $= "k3"`,
		`C1 = "k"`,
		`C2 * "kl"`,
		`C3==`,
		`C1 == "A" or C2 == "B"`,
		`C1 = "A" or C2 & "B"`,
		`C1!="a1" and C3=="u2"`,
		`C1 & "12"`,
	}
	for _, query := range querys {
		fmt.Printf("query: %s, %v\n", query, isValidQuery(query))
	}
}

func TestIsValidModify(t *testing.T) {
	modifys := []string{
		`INSERT a1,a2,a3`,
		`INSERT a1,a2`,
		`INSERT a1`,
		`DELETE a1`,
		`DELETE a1,a2`,
		`DELETE a1,a2,a3`,
		`UPDATE a1,C2,b2`,
		`UPDATE a1,a2,C3,d3`,
	}

	for _, modify := range modifys {
		fmt.Printf("modify: %s, %v\n", modify, isValidModify(modify))
	}
}
