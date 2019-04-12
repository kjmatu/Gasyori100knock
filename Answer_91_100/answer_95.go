package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

type NeuralNetwork struct {
	W1   [][]float64
	B1   []float64
	Wout [][]float64
	Bout []float64
	Lr   float64
	Z1   [][]float64
	Z2   [][]float64
	Out  [][]float64
}

func (nn *NeuralNetwork) dot(matrix1 [][]float64, matrix2 [][]float64) ([][]float64, error) {
	height1 := len(matrix1)
	width1 := len(matrix1[0])

	height2 := len(matrix2)
	width2 := len(matrix2[0])

	if width1 != height2 {
		return nil, errors.New("input matrix does not match size")
	}

	dotArray := make([][]float64, height1)
	for i := range dotArray {
		dotArray[i] = make([]float64, width2)
	}

	for y, row := range dotArray {
		for x := range row {
			// fmt.Println(matrix1[y][:])
			// fmt.Println(matrix2[:][x])
			dotValue := 0.0
			for i := 0; i < len(matrix1[y][:]); i++ {
				dotValue += matrix1[y][i] * matrix2[i][x]
			}
			dotArray[y][x] = dotValue
		}
	}
	// fmt.Println(dotArray)
	return dotArray, nil
}

func (nn *NeuralNetwork) transpose(matrix [][]float64) [][]float64 {

	transposeMatrix := make([][]float64, len(matrix[0]))
	for i := range transposeMatrix {
		transposeMatrix[i] = make([]float64, len(matrix))
	}

	for i, row := range transposeMatrix {
		for j := range row {
			transposeMatrix[i][j] = matrix[j][i]
		}
	}
	return transposeMatrix
}

func (nn *NeuralNetwork) Forward(x [][]float64) [][]float64 {
	nn.Z1 = x

	// Z2
	dotValue, err := nn.dot(nn.Z1, nn.W1)
	if err != nil {
		fmt.Println("Dot error")
	}

	for i, row := range dotValue {
		for j := range row {
			dotValue[i][j] += nn.B1[j]
		}
	}

	nn.Z2 = make([][]float64, len(dotValue))
	for i := range nn.Z2 {
		nn.Z2[i] = make([]float64, len(dotValue[0]))
	}

	for i, row := range nn.Z2 {
		for j := range row {
			nn.Z2[i][j] = sigmoid(dotValue[i][j])
		}
	}

	// Out
	dotValue, err = nn.dot(nn.Z2, nn.Wout)
	if err != nil {
		fmt.Println("Dot error")
	}

	for i, row := range dotValue {
		for j := range row {
			dotValue[i][j] += nn.Bout[j]
		}
	}

	nn.Out = make([][]float64, len(dotValue))
	for i := range nn.Out {
		nn.Out[i] = make([]float64, len(dotValue[0]))
	}

	for i, row := range nn.Out {
		for j := range row {
			nn.Out[i][j] = sigmoid(dotValue[i][j])
		}
	}
	return nn.Out
}

func (nn *NeuralNetwork) Train(x, t [][]float64) {
	// backpropagation output layer
	En := make([][]float64, len(nn.Out))
	for i := range En {
		En[i] = make([]float64, len(nn.Out[0]))
	}

	for i, row := range nn.Out {
		for j := range row {
			En[i][j] = (nn.Out[i][j] - t[i][j]) * nn.Out[i][j] * (1 - nn.Out[i][j])
		}
	}

	// gradEn := append([][]float64{}, En...)

	gradWout, err := nn.dot(nn.transpose(nn.Z2), En)
	if err != nil {
		fmt.Println("Dot error")
	}
	// fmt.Println("gradWout")
	// for _, row := range gradWout {
	// 	fmt.Println(row)
	// }
	// fmt.Println("nn.W1", len(nn.W1), len(nn.W1[0]))
	// for _, row := range nn.W1 {
	// 	fmt.Println(row)
	// }

	ones := make([][]float64, len(En[0]))
	for i := range ones {
		ones[i] = make([]float64, len(En))
	}

	for i, row := range ones {
		for j := range row {
			ones[i][j] = 1.0
		}
	}

	gradBout, err := nn.dot(ones, En)
	if err != nil {
		fmt.Println("Dot error")
	}

	for i, row := range nn.Wout {
		for j := range row {
			nn.Wout[i][j] -= nn.Lr * gradWout[i][j]
		}
	}

	fmt.Println("nn.Wout", len(nn.Wout), len(nn.Wout[0]))
	// for _, row := range nn.Wout {
	// 	fmt.Println(row)
	// }

	nn.Bout[0] -= nn.Lr * gradBout[0][0]

	// backpropagation inter layer
	gradU1, err := nn.dot(En, nn.transpose(nn.Wout))
	// fmt.Println("gradU1", len(gradU1), len(gradU1[0]))
	// fmt.Println("nn.Z2", len(nn.Z2), len(nn.Z2[0]))
	for i, row := range gradU1 {
		for j := range row {
			gradU1[i][j] *= nn.Z2[i][j] * (1 - nn.Z2[i][j])
		}
	}

	// fmt.Println("nn.Wout", len(nn.Wout), len(nn.Wout[0]))
	// for _, row := range nn.Wout {
	// 	fmt.Println(row)
	// }
	// fmt.Println("gradU1", len(gradU1), len(gradU1[0]))
	// for _, row := range gradU1 {
	// 	fmt.Println(row)
	// }

	gradW1, err := nn.dot(nn.transpose(nn.Z1), gradU1)

	fmt.Println("gradU1", len(gradU1), len(gradU1[0]))
	ones = make([][]float64, 1)
	for i := range ones {
		ones[i] = make([]float64, len(gradU1))
	}

	for i, row := range ones {
		for j := range row {
			ones[i][j] = 1.0
		}
	}

	gradB1, err := nn.dot(ones, gradU1)
	// fmt.Println("gradB1", len(gradB1), len(gradB1[0]))

	for i, row := range nn.W1 {
		for j := range row {
			nn.W1[i][j] -= nn.Lr * gradW1[i][j]
		}
	}

	// fmt.Println("nn.B1", len(nn.B1))
	// fmt.Println("gradB1", len(gradB1), len(gradB1[0]))
	for i := range nn.B1 {
		nn.B1[i] -= nn.Lr * gradB1[0][i]
	}
	// for _, elm := range nn.B1 {
	// 	fmt.Println(elm)
	// }

}

func NewNeuralNetwork(index, weight, outd int, lr float64) *NeuralNetwork {
	rand.Seed(0)
	nn := new(NeuralNetwork)

	nn.W1 = make([][]float64, index)
	for i := range nn.W1 {
		nn.W1[i] = make([]float64, weight)
	}

	for i, row := range nn.W1 {
		for j := range row {
			nn.W1[i][j] = rand.NormFloat64()
		}
	}

	nn.B1 = make([]float64, weight)
	for i := range nn.B1 {
		nn.B1[i] = rand.NormFloat64()
	}

	nn.Wout = make([][]float64, weight)
	for i := range nn.Wout {
		nn.Wout[i] = make([]float64, outd)
	}

	for i, row := range nn.Wout {
		for j := range row {
			nn.Wout[i][j] = rand.NormFloat64()
		}
	}

	nn.Bout = make([]float64, outd)
	for i := range nn.Bout {
		nn.Bout[i] = rand.NormFloat64()
	}

	return nn
}

func main() {
	nn := NewNeuralNetwork(2, 64, 1, 0.1)
	x := [][]float64{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	t := [][]float64{{0}, {1}, {1}, {0}}

	for i := 0; i < 1000; i++ {
		nn.Forward(t)
		nn.Train(x, t)
	}

}
