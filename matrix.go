package matrix

import (
	"math/rand"
	"runtime"
	"time"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type Matrix struct {
	rows, cols, step int
	elems            []float64
}

func NewMatrix(rows, cols int) *Matrix {
	m := new(Matrix)
	m.elems = make([]float64, rows*cols)
	m.rows = rows
	m.cols = cols
	m.step = cols
	return m
}

func Normals(rows, cols int) *Matrix {
	A := NewMatrix(rows, cols)

	for i := 0; i < A.rows; i++ {
		for j := 0; j < A.cols; j++ {
			A.Set(i, j, rand.NormFloat64())
		}
	}

	return A
}

func (self *Matrix) Set(i, j int, v float64) {
	self.elems[i*self.step : i*self.step+self.cols][j] = v

}

func (self *Matrix) Get(i, j int) float64 {
	return self.elems[i*self.step : i*self.step+self.cols][j]

}

func (self *Matrix) ParTimes(scnd *Matrix) *Matrix {
	res := NewMatrix(self.rows, scnd.cols)
	parTimes2(self, scnd, res)
	return res
}

func parTimes2(A, B, C *Matrix) {
	const threshold = 8

	currentGoroutineCount := 1
	maxGoroutines := runtime.GOMAXPROCS(0) + 2

	println("maxGoroutines",maxGoroutines)

	var aux func(sync chan bool, A, B, C *Matrix, rs, re, cs, ce, ks, ke int)
	aux = func(sync chan bool, A, B, C *Matrix, rs, re, cs, ce, ks, ke int) {
		dr := re - rs
		dc := ce - cs
		dk := ke - ks
		switch {
		case currentGoroutineCount < maxGoroutines && dr >= dc && dr >= dk && dr >= threshold:
			sync0 := make(chan bool, 1)
			rm := (rs + re) / 2
			currentGoroutineCount++
			go aux(sync0, A, B, C, rs, rm, cs, ce, ks, ke)
			aux(nil, A, B, C, rm, re, cs, ce, ks, ke)
			<-sync0
			currentGoroutineCount--
		case currentGoroutineCount < maxGoroutines && dc >= dk && dc >= dr && dc >= threshold:
			sync0 := make(chan bool, 1)
			cm := (cs + ce) / 2
			currentGoroutineCount++
			go aux(sync0, A, B, C, rs, re, cs, cm, ks, ke)
			aux(nil, A, B, C, rs, re, cm, ce, ks, ke)
			<-sync0
			currentGoroutineCount--
		case currentGoroutineCount < maxGoroutines && dk >= dc && dk >= dr && dk >= threshold:
			km := (ks + ke) / 2
			aux(nil, A, B, C, rs, re, cs, ce, ks, km)
			aux(nil, A, B, C, rs, re, cs, ce, km, ke)
		default:
			for row := rs; row < re; row++ {
				sums := C.elems[row*C.step : (row+1)*C.step]
				for k := ks; k < ke; k++ {
					for col := cs; col < ce; col++ {
						sums[col] += A.elems[row*A.step+k] * B.elems[k*B.step+col]
					}
				}
			}
		}
		if sync != nil {
			sync <- true
		}
	}

	aux(nil, A, B, C, 0, A.rows, 0, B.cols, 0, A.cols)

	return
}

func ParallelProduct(A, B *Matrix) (C *Matrix) {
	if A.cols != B.rows {
		return nil
	}

	C = NewMatrix(A.rows, B.cols)

	in := make(chan int)
	quit := make(chan bool)

	dotRowCol := func() {
		for {
			select {
			case i := <-in:
				sums := make([]float64, B.cols)
				for k := 0; k < A.cols; k++ {
					for j := 0; j < B.cols; j++ {
						sums[j] += A.Get(i, k) * B.Get(k, j)
					}
				}
				for j := 0; j < B.cols; j++ {
					C.Set(i, j, sums[j])
				}
			case <-quit:
				return
			}
		}
	}

	threads := 2

	for i := 0; i < threads; i++ {
		go dotRowCol()
	}

	for i := 0; i < A.rows; i++ {
		in <- i
	}

	for i := 0; i < threads; i++ {
		quit <- true
	}

	return
}

func Equals(A, B *Matrix) bool {
	if A.rows != B.rows || A.cols != B.cols {
		return false
	}
	for i := 0; i < A.rows; i++ {
		for j := 0; j < A.cols; j++ {
			if A.Get(i, j) != B.Get(i, j) {
				return false
			}
		}
	}
	return true
}

func (A *Matrix) Times(B *Matrix) *Matrix {
	C := NewMatrix(A.rows, B.cols)

	for i := 0; i < A.rows; i++ {
		for j := 0; j < B.cols; j++ {
			sum := float64(0)
			for k := 0; k < A.cols; k++ {
				sum += A.elems[i*A.step+k] * B.Get(k, j)
			}
			C.elems[i*C.step+j] = sum
		}
	}
	return C
}
