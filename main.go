package main

import (
	"flag"
	"fmt"
	"math/big"
	"sync"
)

var (
	zero *big.Int = big.NewInt(0)
	one  *big.Int = big.NewInt(1)
	two  *big.Int = big.NewInt(2)

	mu         sync.Mutex
	numWorkers = 4
	jobs       = make(chan uint, numWorkers)

	smallPrimes = [...]uint{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101}

	currentP = uint(1)
)

const (
	iterationsPerJob = 128
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	fmt.Println("Starting to look for primes")

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go spawnWorker(wg)
	}

	// We neved call wg.Done, so we'll just keep searching
	wg.Wait()
}

func spawnWorker(wg *sync.WaitGroup) {
	for {
		// Reserve some primes for us to work with
		mu.Lock()
		initialP := currentP
		currentP += iterationsPerJob
		mu.Unlock()

		for i := uint(0); i < iterationsPerJob; i++ {
			myP := initialP + i
			// Little bit of initial factoring because Lucas-Lehmer is expensive
			mightBePrime := true
			for j := uint(0); j < uint(len(smallPrimes)); j++ {
				if myP%smallPrimes[j] == 0 {
					mightBePrime = false
					break
				}
			}

			if mightBePrime && LucasLehmer(myP) {
				fmt.Printf("2^%d-1 is prime\n", myP)
			}
		}
	}
}

// p here is the same p as in 2^p-1
func LucasLehmer(p uint) (isPrime bool) {
	var dummy1, dummy2 big.Int

	s := big.NewInt(4)
	m := big.NewInt(0)
	m = m.Sub(m.Lsh(one, p), one) // = (1 << p) - 1

	for i := 0; i < int(p)-2; i++ {
		// This is the slower but straightforward way
		//s.Mod(s.Sub(s.Mul(s, s), two), m)

		// Or, use this faster way
		// s here is the same as k in this part of the wikipedia page:
		// https://en.wikipedia.org/wiki/Lucas%E2%80%93Lehmer_primality_test#Time_complexity
		s = s.Sub(s.Mul(s, s), two)

		for s.Cmp(m) == 1 {
			// And is big's logical and, Rsh is right shift
			s.Add(dummy1.And(s, m), dummy2.Rsh(s, p))
		}

		if s.Cmp(m) == 0 {
			s = zero
		}
	}

	return s.Cmp(zero) == 0
}
