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

func ParallelProduct(A, B *Matrix) *Matrix {
	if A.cols != B.rows {
		return nil
	}

	C := NewMatrix(A.rows, B.cols)

	in := make(chan int)
	quit := make(chan bool)

	threads := runtime.GOMAXPROCS(0) * 2

	for i := 0; i < threads; i++ {
		go func() {
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
		}()
	}

	for i := 0; i < A.rows; i++ {
		in <- i
	}

	for i := 0; i < threads; i++ {
		quit <- true
	}

	return C
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
