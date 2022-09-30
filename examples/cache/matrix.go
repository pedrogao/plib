package cache

func createMatrix(size int) [][]int64 {
	matrix := make([][]int64, size)
	for i := 0; i < size; i++ {
		matrix[i] = make([]int64, size)
		for j := 0; j < size; j++ {
			matrix[i][j] = 2
		}
	}
	return matrix
}
