# Lucas-Lehmer
An implementation of the [Lucas-Lehmer](https://en.wikipedia.org/wiki/Lucas–Lehmer_primality_test) primality test in go (with a few optimizations).

### Optimizations
 * For all Mersenne primes (of the form 2<sup>p</sup>-1) p has to be a prime number, so before doing the more costly Lucas-lehmer test we make sure p is prime.
 * The modulus in the Lucas-Lehmer test is made a lot faster in the fastMod function which uses bitwise operators instead of division, moving the slowest part of the test to the multiplication of s*s.
