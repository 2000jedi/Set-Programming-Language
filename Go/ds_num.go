package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type number struct {
	numerator   int
	denominator int
}

var True, False storage

func (n *number) construct(v string) {
	raw := strings.Split(v, ".")
	switch len(raw) {
	case 1:
		temp, err := strconv.Atoi(raw[0])
		if err != nil {
			panic(err)
		}
		n.numerator = temp
		n.denominator = 1
	case 2:
		temp, err := strconv.Atoi(raw[0])
		if err != nil {
			panic(err)
		}
		n.numerator = temp
		temp, err = strconv.Atoi(raw[1])
		if err != nil {
			panic(err)
		}
		n.denominator = int(math.Pow10(int(math.Ceil(math.Log10(float64(temp + 1))))))
		n.numerator = n.numerator*n.denominator + temp
	default:
		panic("Invalid number: " + v + "\n")
	}
	n.reduce()
}

func (n *number) toString() string {
	if n.denominator == 1 {
		return fmt.Sprint(n.numerator)
	} else {
		return fmt.Sprintf("%d/%d", n.numerator, n.denominator)
	}
}

func (n *number) toInt() int {
	return n.numerator / n.denominator
}

func (n *number) print() {
	fmt.Print(n.toString())
}

func init() {
	True = storage{VAR_NUMBER, number{1, 1}}
	False = storage{VAR_NUMBER, number{0, 0}}
}
