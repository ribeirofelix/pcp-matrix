package matrix

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const ε = 0.000001
const verbose = true
const speedTest = true

func TestParallelProduct(t *testing.T) {

	w := 100000
	h := 40

	rand.Seed(time.Now().UnixNano())
	A := Normals(h, w)
	B := Normals(w, h)

	var C *Matrix
	var start, end int64

	start = time.Now().UnixNano()
	Ctrue := A.ParTimes(B)

	end = time.Now().UnixNano()
	if verbose {
		fmt.Printf("%fs for synchronous\n", float64(end-start)/1000000000)
	}

	start = time.Now().UnixNano()
	C = ParallelProduct(A, B)

	end = time.Now().UnixNano()
	if verbose {
		fmt.Printf("%fs for parallel\n", float64(end-start)/1000000000)
	}


	if !Equals(C, Ctrue) {
		t.Fail()
	}

	start = time.Now().UnixNano()
	Cs := A.Times(B)

	end = time.Now().UnixNano()
	if verbose {
		fmt.Printf("%fs for serial\n", float64(end-start)/1000000000)
	}

	if !Equals(Cs, Ctrue) || !Equals(Cs, C) {
		t.Fail()
	}


}
