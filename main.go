package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"os/signal"
	"syscall"
)

// Essentially constants used in loop body of Lucas-Lehmer
var one *big.Int = big.NewInt(1)
var two *big.Int = big.NewInt(2)

func main() {
	fmt.Println("Starting to look for primes")
	primes := "" // We'll use this to store all the primes we find

	// On exit (or ctrl+c) print all found primes
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nPrimes found: ", primes)
		os.Exit(1)
	}()

	i := uint(1)

	for {
		// Mersenne primes are always raised to a prime power
		// So, first make sure that i (our power) is prime
		isIPrime := true

		for j := uint(2); j<uint(math.Sqrt(float64(i))); j++ {
			if i % j == 0 {
				isIPrime = false
				break
			}
		}

		if isIPrime && LucasLehmer(i) {
			fmt.Printf("2^%d-1 is prime\n", i)
			// Some wack way to concat an int (well uint)
			primes += fmt.Sprint(i) + ", "
		} else {
			fmt.Printf("2^%d-1 is not prime\n", i)
		}

		// Primes can't be even, no need to check
		i += 2
	}
}

// p here is the same p as in 2^p-1
func LucasLehmer(p uint) (isPrime bool) {
	// Lsh is like << but for math/big.Int (big bois)
	// m is the mersenne prime in it's full form
	m := big.NewInt(0)
	m = m.Sub(m.Lsh(one, p), one)

	s := big.NewInt(4)

	for i := 0; i<int(p)-2; i++  {
		// This is the obvious way to perform the check, however...
		// The modulus can be done quicker using the fastMod method
		//s.Mod(s.Sub(s.Mul(s, s), two), m)

		// This is the much quicker method, and according to wikipedia
		// The bottleneck will be in multiplying s*s
		s = fastMod(s.Sub(s.Mul(s, s), two), m, p)
	}

	return s.Cmp(big.NewInt(0)) == 0
}

// s here is the same as k in this wikipedia pages:
// https://en.wikipedia.org/wiki/Lucas%E2%80%93Lehmer_primality_test#Time_complexity
func fastMod(s *big.Int, m *big.Int, p uint) *big.Int {
	// So why are there dummy variables?
	// Well, math/big needs a *big.Int to perform the methods on as the
	// implicit param and I couldn't really find a way around it, so...
	var dummy1, dummy2 big.Int

	for s.Cmp(m) > 0 {
		// And is big's logical and, Rsh is right shift (>>)
		s.Add(dummy1.And(s, m), dummy2.Rsh(s, p))
	}

	if(s.Cmp(m) == 0) {
		return big.NewInt(0)
	} else {
		return s
	}
}
