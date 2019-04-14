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
	// fmt.Println("nn.Z1", len(nn.Z1), len(nn.Z1[0]))
	// fmt.Println(nn.Z1)
	// fmt.Println("nn.W1", len(nn.W1), len(nn.W1[0]))
	// fmt.Println(nn.W1)
	dotValue, err := nn.dot(nn.Z1, nn.W1)
	if err != nil {
		fmt.Println("Dot error")
	}
	// fmt.Println("dotValue", dotValue)

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
	// fmt.Println("nn.Z2", len(nn.Z2), nn.Z2)
	// fmt.Println("nn.Wout", len(nn.Wout), len(nn.Wout[0]))
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

	// fmt.Println("nn.Wout", len(nn.Wout), len(nn.Wout[0]))
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

	// fmt.Println("gradU1", len(gradU1), len(gradU1[0]))
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
	nn.W1 = [][]float64{
		{1.76405235, 0.40015721, 0.97873798, 2.2408932, 1.86755799, -0.97727788, 0.95008842, -0.15135721, -0.10321885, 0.4105985, 0.14404357, 1.45427351, 0.76103773, 0.12167502, 0.44386323, 0.33367433, 1.49407907, -0.20515826, 0.3130677, -0.85409574, -2.55298982, 0.6536186, 0.8644362, -0.74216502, 2.26975462, -1.45436567, 0.04575852, -0.18718385, 1.53277921, 1.46935877, 0.15494743, 0.37816252, -0.88778575, -1.98079647, -0.34791215, 0.15634897, 1.23029068, 1.20237985, -0.38732682, -0.30230275, -1.04855297, -1.42001794, -1.70627019, 1.9507754, -0.50965218, -0.4380743, -1.25279536, 0.77749036, -1.61389785, -0.21274028, -0.89546656, 0.3869025, -0.51080514, -1.18063218, -0.02818223, 0.42833187, 0.06651722, 0.3024719, -0.63432209, -0.36274117, -0.67246045, -0.35955316, -0.81314628, -1.7262826},
		{0.17742614, -0.40178094, -1.63019835, 0.46278226, -0.90729836, 0.0519454, 0.72909056, 0.12898291, 1.13940068, -1.23482582, 0.40234164, -0.68481009, -0.87079715, -0.57884966, -0.31155253, 0.05616534, -1.16514984, 0.90082649, 0.46566244, -1.53624369, 1.48825219, 1.89588918, 1.17877957, -0.17992484, -1.07075262, 1.05445173, -0.40317695, 1.22244507, 0.20827498, 0.97663904, 0.3563664, 0.70657317, 0.01050002, 1.78587049, 0.12691209, 0.40198936, 1.8831507, -1.34775906, -1.270485, 0.96939671, -1.17312341, 1.94362119, -0.41361898, -0.74745481, 1.92294203, 1.48051479, 1.86755896, 0.90604466, -0.86122569, 1.91006495, -0.26800337, 0.8024564, 0.94725197, -0.15501009, 0.61407937, 0.92220667, 0.37642553, -1.09940079, 0.29823817, 1.3263859, -0.69456786, -0.14963454, -0.43515355, 1.84926373}}

	nn.B1 = []float64{0.67229476, 0.40746184, 0.76991607, 0.53924919, 0.67433266, 0.03183056, -0.63584608, 0.67643329, 0.57659082, 0.20829876, 0.39600671, 1.09306151, -1.49125759, 0.4393917, 0.1666735, 0.63503144, 2.38314477, 0.94447949, -0.91282223, 1.11701629, 1.31590741, 0.4615846, -0.06824161, 1.71334272, -0.74475482, 0.82643854, 0.09845252, 0.66347829, 1.12663592, 1.07993151, -1.14746865, 0.43782004, 0.49803245, 1.92953205, 0.94942081, 0.08755124, -1.22543552, 0.84436298, 1.00021535, 1.5447711, 1.18802979, 0.31694261, 0.92085882, 0.31872765, 0.85683061, 0.65102559, 1.03424284, 0.68159452, -0.80340966, 0.68954978, 0.4555325, 0.01747916, 0.35399391, 1.37495129, -0.6436184, -2.22340315, 0.62523145, 1.60205766, 1.10438334, 0.05216508, -0.739563, 1.5430146, -1.29285691, 0.26705087}

	nn.Wout = [][]float64{
		{-0.03928282}, {-1.1680935}, {0.52327666}, {-0.17154633}, {0.77179055}, {0.82350415}, {2.16323595}, {1.33652795}, {-0.36918184}, {-0.23937918}, {1.0996596}, {0.65526373}, {0.64013153}, {-1.61695604}, {-0.02432612}, {-0.73803091}, {0.2799246}, {-0.09815039}, {0.91017891}, {0.31721822}, {0.78632796}, {-0.4664191}, {-0.94444626}, {-0.41004969}, {-0.01702041}, {0.37915174}, {2.25930895}, {-0.04225715}, {-0.955945}, {-0.34598178}, {-0.46359597}, {0.48148147}, {-1.54079701}, {0.06326199}, {0.15650654}, {0.23218104}, {-0.59731607}, {-0.23792173}, {-1.42406091}, {-0.49331988}, {-0.54286148}, {0.41605005}, {-1.15618243}, {0.7811981}, {1.49448454}, {-2.06998503}, {0.42625873}, {0.67690804}, {-0.63743703}, {-0.39727181}, {-0.13288058}, {-0.29779088}, {-0.30901297}, {-1.67600381}, {1.15233156}, {1.07961859}, {-0.81336426}, {-1.46642433}, {0.52106488}, {-0.57578797}, {0.14195316}, {-0.31932842}, {0.69153875}, {0.69474914}}

	nn.Bout = []float64{-0.72559738}
	x := [][]float64{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	t := [][]float64{{0}, {1}, {1}, {0}}

	fmt.Println("nn.W1")
	for _, row := range nn.W1 {
		fmt.Println(row)
	}
	fmt.Println("nn.B1")
	for _, row := range nn.B1 {
		fmt.Println(row)
	}
	fmt.Println("nn.Wout")
	for _, row := range nn.Wout {
		fmt.Println(row)
	}
	fmt.Println("nn.Bout")
	fmt.Println(nn.Bout)

	// fmt.Println(nn)
	// train
	for i := 0; i < 1000; i++ {
		nn.Forward(x)
		nn.Train(x, t)
	}

	// test
	for j := 0; j < 4; j++ {
		xx := x[j]
		// tt := t[j]
		fmt.Println("in:", xx, "pred:", nn.Forward([][]float64{xx}))
	}

}
