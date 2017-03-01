package main

func gcd_(a, b int) int {
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
	gcd := gcd_(self.denominator, self.numerator)
	self.denominator /= gcd
	self.numerator /= gcd
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
	gcd := n.denominator * n_.denominator / gcd_(n.denominator, n_.denominator)
	res = number{n.numerator*gcd/n.denominator + n_.numerator*gcd/n_.denominator, gcd}
	res.reduce()
	return
}

func sub(n, n_ number) (res number) {
	gcd := n.denominator * n_.denominator / gcd_(n.denominator, n_.denominator)
	res = number{n.numerator*gcd/n.denominator - n_.numerator*gcd/n_.denominator, gcd}
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
	if lambda.vartype == var_fsm["function"] {
		f := lambda.data.(function)
		return f.function(argc, variable)
	} else {
		f := lambda.data.(inherit)
		return f.function(argc, variable)
	}
}
