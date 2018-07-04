package main

import (
	"fmt"
	"log"
)

func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}
func extendeGcd(x, y int) (gcd, u, v int) {
	u, old_u := 0, 1
	v, old_v := 1, 0
	r, old_r := y, x

	for r != 0 {
		quotient := int(old_r / r)
		old_r, r = r, old_r-quotient*r
		old_u, u = u, old_u-quotient*u
		old_v, v = v, old_v-quotient*v
	}

	return old_r, old_u, old_v
}

func multInvser(x, p int) int {
	gcd, u, _ := extendeGcd(x, p)
	if gcd != 1 {
		log.Println("Only relative prime can inverse")
		return 0
	} else {
		t := u % p
		if t < 0 {
			return p + t
		} else {
			return u % p
		}
	}
}
func main() {
	z := gcd(15, 25)
	fmt.Println(z)

	u, w, s := extendeGcd(15, 25)
	fmt.Println(u, w, s)

	x := multInvser(2, 5)
	fmt.Println(x)
}
