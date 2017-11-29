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

var True, False number

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

func (n number) toString() string {
	if n.denominator == 1 {
		return strconv.Itoa(n.numerator)
	} else {
		return fmt.Sprintf("%d/%d", n.numerator, n.denominator)
	}
}

func (n *number) toInt() int {
	return n.numerator / n.denominator
}

func gcd(a, b int) int {
	for b > 0 {
		a, b = b, a%b
	}
	return a
}

func (self *number) reduce() {
	if self.numerator == 0 {
		if self.denominator != 0 {
			self.denominator = 1
		}
		return
	}
	g := gcd(self.denominator, self.numerator)
	self.denominator /= g
	self.numerator /= g
	if self.denominator < 0 {
		self.numerator = -self.numerator
		self.denominator = -self.denominator
	}
}

func equal(n, n_ number) number {
	if n.numerator == n_.numerator && n.denominator == n_.denominator {
		return True
	}
	return False
}

func lt(n, n_ number) number {
	if n.numerator*1.0/n.denominator < n_.numerator*1.0/n_.denominator {
		return True
	}
	return False
}

func le(n, n_ number) number {
	if n.numerator*1.0/n.denominator <= n_.numerator*1.0/n_.denominator {
		return True
	}
	return False
}

func gt(n, n_ number) number {
	if n.numerator*1.0/n.denominator > n_.numerator*1.0/n_.denominator {
		return True
	}
	return False
}

func ge(n, n_ number) number {
	if n.numerator*1.0/n.denominator >= n_.numerator*1.0/n_.denominator {
		return True
	}
	return False
}

func add(n, n_ number) (res number) {
	g := n.denominator * n_.denominator / gcd(n.denominator, n_.denominator)
	res = number{n.numerator*g/n.denominator + n_.numerator*g/n_.denominator, g}
	res.reduce()
	return
}

func sub(n, n_ number) (res number) {
	g := n.denominator * n_.denominator / gcd(n.denominator, n_.denominator)
	res = number{n.numerator*g/n.denominator - n_.numerator*g/n_.denominator, g}
	res.reduce()
	return
}

func mul(n, n_ number) (res number) {
	res = number{n.numerator * n_.numerator, n.denominator * n_.denominator}
	res.reduce()
	return
}

func div(n, n_ number) (res number) {
	res = number{n.numerator * n_.denominator, n.denominator * n_.numerator}
	res.reduce()
	return
}

func init() {
	True = number{1, 1}
	False = number{0, 1}
}
