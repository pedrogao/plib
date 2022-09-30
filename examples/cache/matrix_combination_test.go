package cache

import (
	"sync"
	"testing"
)

var matrixLength = 512

// 6400:
// BenchmarkMatrixCombination
// BenchmarkMatrixCombination-12    	      12	  96176104 ns/op
// 512:
// BenchmarkMatrixCombination
// BenchmarkMatrixCombination-12    	    3226	    385806 ns/op
func BenchmarkMatrixCombination(b *testing.B) {
	matrixA := createMatrix(matrixLength)
	matrixB := createMatrix(matrixLength)

	for n := 0; n < b.N; n++ {
		for i := 0; i < matrixLength; i++ {
			for j := 0; j < matrixLength; j++ {
				matrixA[i][j] = matrixA[i][j] + matrixB[i][j]
			}
		}
	}
}

// 6400:
// BenchmarkMatrixReversedCombination
// BenchmarkMatrixReversedCombination-12    	       1	1289900702 ns/op
// 512:
// BenchmarkMatrixReversedCombination
// BenchmarkMatrixReversedCombination-12    	    1389	   1151074 ns/op
func BenchmarkMatrixReversedCombination(b *testing.B) {
	matrixA := createMatrix(matrixLength)
	matrixB := createMatrix(matrixLength)

	for n := 0; n < b.N; n++ {
		for i := 0; i < matrixLength; i++ {
			for j := 0; j < matrixLength; j++ {
				matrixA[i][j] = matrixA[i][j] + matrixB[j][i]
			}
		}
	}
}

// 512:
// BenchmarkMatrixReversedCombinationPerBlock
// BenchmarkMatrixReversedCombinationPerBlock-12    	    1802	    812322 ns/op
func BenchmarkMatrixReversedCombinationPerBlock(b *testing.B) {
	matrixA := createMatrix(matrixLength)
	matrixB := createMatrix(matrixLength)
	blockSize := 8

	for n := 0; n < b.N; n++ {
		for i := 0; i < matrixLength; i += blockSize {
			for j := 0; j < matrixLength; j += blockSize {
				// 64 = 8 * 8
				// 每次计算一个block，一个block是8个int64，刚好一个cache line
				for ii := i; ii < i+blockSize; ii++ {
					for jj := j; jj < j+blockSize; jj++ {
						matrixA[ii][jj] = matrixA[ii][jj] + matrixB[jj][ii]
					}
				}
			}
		}
	}
}

const (
	M                = 1_000_000
	CacheLinePadSize = 64
)

type SimpleStruct struct {
	n int
}

// BenchmarkStructureFalseSharing
// BenchmarkStructureFalseSharing-12    	     561	   2151351 ns/op
func BenchmarkStructureFalseSharing(b *testing.B) {
	structA := SimpleStruct{}
	structB := SimpleStruct{}
	wg := sync.WaitGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func() {
			for j := 0; j < M; j++ {
				structA.n += j
			}
			wg.Done()
		}()
		go func() {
			for j := 0; j < M; j++ {
				structB.n += j
			}
			wg.Done()
		}()
		wg.Wait()
	}
}

type PaddedStruct struct {
	n int
	_ CacheLinePad
}

type CacheLinePad struct {
	_ [CacheLinePadSize]byte
}

// BenchmarkStructurePadding
// BenchmarkStructurePadding-12    	     771	   1658719 ns/op
func BenchmarkStructurePadding(b *testing.B) {
	structA := PaddedStruct{}
	structB := SimpleStruct{}
	wg := sync.WaitGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func() {
			for j := 0; j < M; j++ {
				structA.n += j
			}
			wg.Done()
		}()
		go func() {
			for j := 0; j < M; j++ {
				structB.n += j
			}
			wg.Done()
		}()
		wg.Wait()
	}
}
