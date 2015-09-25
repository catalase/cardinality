package main

import (
	"math"
)

var NDT = [...][2]float64{
	{0, 0.0},
	{0.01, 0.008},
	{0.02, 0.016},
	{0.03, 0.024},
	{0.04, 0.032},
	{0.05, 0.0398},
	{0.06, 0.0478},
	{0.07, 0.0558},
	{0.08, 0.0638},
	{0.09, 0.0718},
	{0.10, 0.0796},
	{0.11, 0.0876},
	{0.12, 0.0956},
	{0.13, 0.1034},
	{0.14, 0.1114},
	{0.15, 0.1192},
	{0.16, 0.1272},
	{0.17, 0.135},
	{0.18, 0.1428},
	{0.19, 0.1506},
	{0.20, 0.1586},
	{0.21, 0.1664},
	{0.22, 0.1742},
	{0.23, 0.182},
	{0.24, 0.1896},
	{0.25, 0.1974},
	{0.26, 0.2052},
	{0.27, 0.2128},
	{0.28, 0.2206},
	{0.29, 0.2282},
	{0.30, 0.2358},
	{0.31, 0.2434},
	{0.32, 0.251},
	{0.33, 0.2586},
	{0.34, 0.2662},
	{0.35, 0.2736},
	{0.36, 0.2812},
	{0.37, 0.2886},
	{0.38, 0.296},
	{0.39, 0.3034},
	{0.40, 0.3108},
	{0.41, 0.3182},
	{0.42, 0.3256},
	{0.43, 0.3328},
	{0.44, 0.34},
	{0.45, 0.3472},
	{0.46, 0.3544},
	{0.47, 0.3616},
	{0.48, 0.3688},
	{0.49, 0.3758},
	{0.50, 0.383},
	{0.51, 0.39},
	{0.52, 0.397},
	{0.53, 0.4038},
	{0.54, 0.4108},
	{0.55, 0.4176},
	{0.56, 0.4246},
	{0.57, 0.4314},
	{0.58, 0.438},
	{0.59, 0.4448},
	{0.60, 0.4514},
	{0.61, 0.4582},
	{0.62, 0.4648},
	{0.63, 0.4714},
	{0.64, 0.4778},
	{0.65, 0.4844},
	{0.66, 0.4908},
	{0.67, 0.4972},
	{0.68, 0.5036},
	{0.69, 0.5098},
	{0.70, 0.516},
	{0.71, 0.5222},
	{0.72, 0.5284},
	{0.73, 0.5346},
	{0.74, 0.5408},
	{0.75, 0.5468},
	{0.76, 0.5528},
	{0.77, 0.5588},
	{0.78, 0.5646},
	{0.79, 0.5704},
	{0.80, 0.5762},
	{0.81, 0.582},
	{0.82, 0.5878},
	{0.83, 0.5934},
	{0.84, 0.599},
	{0.85, 0.6046},
	{0.86, 0.6102},
	{0.87, 0.6156},
	{0.88, 0.6212},
	{0.89, 0.6266},
	{0.90, 0.6318},
	{0.91, 0.6372},
	{0.92, 0.6424},
	{0.93, 0.6476},
	{0.94, 0.6528},
	{0.95, 0.6578},
	{0.96, 0.663},
	{0.97, 0.668},
	{0.98, 0.673},
	{0.99, 0.6778},
	{1.00, 0.6826},
	{1.01, 0.6876},
	{1.02, 0.6922},
	{1.03, 0.697},
	{1.04, 0.7016},
	{1.05, 0.7062},
	{1.06, 0.7108},
	{1.07, 0.7154},
	{1.08, 0.7198},
	{1.09, 0.7242},
	{1.10, 0.7286},
	{1.11, 0.733},
	{1.12, 0.7372},
	{1.13, 0.7416},
	{1.14, 0.7458},
	{1.15, 0.7498},
	{1.16, 0.754},
	{1.17, 0.758},
	{1.18, 0.762},
	{1.19, 0.766},
	{1.20, 0.7698},
	{1.21, 0.7738},
	{1.22, 0.7776},
	{1.23, 0.7814},
	{1.24, 0.785},
	{1.25, 0.7888},
	{1.26, 0.7924},
	{1.27, 0.796},
	{1.28, 0.7994},
	{1.29, 0.803},
	{1.30, 0.8064},
	{1.31, 0.8098},
	{1.32, 0.8132},
	{1.33, 0.8164},
	{1.34, 0.8198},
	{1.35, 0.823},
	{1.36, 0.8262},
	{1.37, 0.8294},
	{1.38, 0.8324},
	{1.39, 0.8354},
	{1.40, 0.8384},
	{1.41, 0.8414},
	{1.42, 0.8444},
	{1.43, 0.8472},
	{1.44, 0.8502},
	{1.45, 0.853},
	{1.46, 0.8558},
	{1.47, 0.8584},
	{1.48, 0.8612},
	{1.49, 0.8638},
	{1.50, 0.8664},
	{1.51, 0.869},
	{1.52, 0.8714},
	{1.53, 0.874},
	{1.54, 0.8764},
	{1.55, 0.8788},
	{1.56, 0.8812},
	{1.57, 0.8836},
	{1.58, 0.8858},
	{1.59, 0.8882},
	{1.60, 0.8904},
	{1.61, 0.8926},
	{1.62, 0.8948},
	{1.63, 0.8968},
	{1.64, 0.899},
	{1.65, 0.901},
	{1.66, 0.903},
	{1.67, 0.905},
	{1.68, 0.907},
	{1.69, 0.909},
	{1.70, 0.9108},
	{1.71, 0.9128},
	{1.72, 0.9146},
	{1.73, 0.9164},
	{1.74, 0.9182},
	{1.75, 0.9198},
	{1.76, 0.9216},
	{1.77, 0.9232},
	{1.78, 0.925},
	{1.79, 0.9266},
	{1.80, 0.9282},
	{1.81, 0.9298},
	{1.82, 0.9312},
	{1.83, 0.9328},
	{1.84, 0.9342},
	{1.85, 0.9356},
	{1.86, 0.9372},
	{1.87, 0.9386},
	{1.88, 0.9398},
	{1.89, 0.9412},
	{1.90, 0.9426},
	{1.91, 0.9438},
	{1.92, 0.9452},
	{1.93, 0.9464},
	{1.94, 0.9476},
	{1.95, 0.9488},
	{1.96, 0.95},
	{1.97, 0.9512},
	{1.98, 0.9522},
	{1.99, 0.9534},
	{2.00, 0.9544},
	{2.01, 0.9556},
	{2.02, 0.9566},
	{2.03, 0.9576},
	{2.04, 0.9586},
	{2.05, 0.9596},
	{2.06, 0.9606},
	{2.07, 0.9616},
	{2.08, 0.9624},
	{2.09, 0.9634},
	{2.10, 0.9642},
	{2.11, 0.9652},
	{2.12, 0.966},
	{2.13, 0.9668},
	{2.14, 0.9676},
	{2.15, 0.9684},
	{2.16, 0.9692},
	{2.17, 0.97},
	{2.18, 0.9708},
	{2.19, 0.9714},
	{2.20, 0.9722},
	{2.21, 0.9728},
	{2.22, 0.9736},
	{2.23, 0.9742},
	{2.24, 0.975},
	{2.25, 0.9756},
	{2.26, 0.9762},
	{2.27, 0.9768},
	{2.28, 0.9774},
	{2.29, 0.978},
	{2.30, 0.9786},
	{2.31, 0.9792},
	{2.32, 0.9796},
	{2.33, 0.9802},
	{2.34, 0.9808},
	{2.35, 0.9812},
	{2.36, 0.9818},
	{2.37, 0.9822},
	{2.38, 0.9826},
	{2.39, 0.9832},
	{2.40, 0.9836},
	{2.41, 0.984},
	{2.42, 0.9844},
	{2.43, 0.985},
	{2.44, 0.9854},
	{2.45, 0.9858},
	{2.46, 0.9862},
	{2.47, 0.9864},
	{2.48, 0.9868},
	{2.49, 0.9872},
	{2.50, 0.9876},
	{2.51, 0.988},
	{2.52, 0.9882},
	{2.53, 0.9886},
	{2.54, 0.989},
	{2.55, 0.9892},
	{2.56, 0.9896},
	{2.57, 0.9898},
	{2.58, 0.9902},
	{2.59, 0.9904},
	{2.60, 0.9906},
	{2.61, 0.991},
	{2.62, 0.9912},
	{2.63, 0.9914},
	{2.64, 0.9918},
	{2.65, 0.992},
	{2.66, 0.9922},
	{2.67, 0.9924},
	{2.68, 0.9926},
	{2.69, 0.9928},
	{2.70, 0.993},
	{2.71, 0.9932},
	{2.72, 0.9934},
	{2.73, 0.9936},
	{2.74, 0.9938},
	{2.75, 0.994},
	{2.76, 0.9942},
	{2.77, 0.9944},
	{2.78, 0.9946},
	{2.79, 0.9948},
	{2.80, 0.9948},
	{2.81, 0.995},
	{2.82, 0.9952},
	{2.83, 0.9954},
	{2.84, 0.9954},
	{2.85, 0.9956},
	{2.86, 0.9958},
	{2.87, 0.996},
	{2.88, 0.996},
	{2.89, 0.9962},
	{2.90, 0.9962},
	{2.91, 0.9964},
	{2.92, 0.9966},
	{2.93, 0.9966},
	{2.94, 0.9968},
	{2.95, 0.9968},
	{2.96, 0.979},
	{2.97, 0.997},
	{2.98, 0.9972},
	{2.99, 0.9972},
	{3.00, 0.9974},
	{3.01, 0.9974},
	{3.02, 0.9974},
	{3.03, 0.9976},
	{3.04, 0.9976},
	{3.05, 0.9978},
	{3.06, 0.9978},
	{3.07, 0.9978},
	{3.08, 0.998},
	{3.09, 0.998},
}

func SolveQuardratic(a, b, c float64) (float64, float64) {
	sqrtD := math.Sqrt(b*b - 4*a*c)
	return (sqrtD - b) / (2 * a), (b + sqrtD) / (2 * a) * -1
}

func Quess(k, m, R, L float64) (float64, float64) {
	max := func(x, y float64) float64 {
		if x > y {
			return x
		}
		return y
	}

	a, b := SolveQuardratic(R, (m-k)*-L, -k*L*L)
	c, d := SolveQuardratic(R, (m+k)*-L, k*L*L)

	return max(a, b), max(c, d)
}