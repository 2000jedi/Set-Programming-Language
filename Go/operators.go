package main

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

func equal(n, n_ number) bool {
	return n.numerator == n_.numerator && n.denominator == n_.denominator
}

func lt(n, n_ number) bool {
	return n.numerator*1.0/n.denominator < n_.numerator*1.0/n.denominator
}

func le(n, n_ number) bool {
	return n.numerator*1.0/n.denominator <= n_.numerator*1.0/n.denominator
}

func gt(n, n_ number) bool {
	return n.numerator*1.0/n.denominator > n_.numerator*1.0/n.denominator
}

func ge(n, n_ number) bool {
	return n.numerator*1.0/n.denominator >= n_.numerator*1.0/n.denominator
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

func do_func(lambda storage, argc []storage, variable *Variable) *storage {
	if lambda.vartype == VAR_FUNCTION {
		f := lambda.data.(function)
		return f.function(argc, variable)
	} else if lambda.vartype == VAR_C_FUNCTION {
		f := lambda.data.(inherit)
		return f.function(argc, variable)
	} else {
		panic("not a function")
	}
}
