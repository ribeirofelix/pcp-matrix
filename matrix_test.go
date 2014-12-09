package matrix

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const verbose = true

func TestParallelProduct(t *testing.T) {

	h := 400

	rand.Seed(time.Now().UnixNano())

	for w := 1000; w < 100000; w += 1000 {
		A := Normals(h, w)
		B := Normals(w, h)

		start := time.Now().UnixNano()
		C := ParallelProduct(A, B)

		end := time.Now().UnixNano()
		if verbose {
			fmt.Printf("%d,%d,%f,parallel\n", w, h, float64(end-start)/1000000000)
		}

		start = time.Now().UnixNano()
		Cs := A.Times(B)

		end = time.Now().UnixNano()
		if verbose {
			fmt.Printf("%d,%d,%f,serial\n", w, h, float64(end-start)/1000000000)
		}

		if !Equals(Cs, C) {
			t.Fail()
		}
	}

}
