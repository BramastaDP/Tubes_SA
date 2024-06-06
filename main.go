package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// BRUTEFORCE
func knapsackBruteForce(weights, values []int, capacity int) (int, []bool) {
	n := len(weights)
	maxValue := 0
	bestSelected := make([]bool, n)

	var knapsackRecursive func(int, int, int, []bool)
	knapsackRecursive = func(index, currentWeight, currentValue int, selected []bool) {
		if index == n {
			if currentWeight <= capacity && currentValue > maxValue {
				maxValue = currentValue
				copy(bestSelected, selected)
			}
			return
		}

		knapsackRecursive(index+1, currentWeight, currentValue, selected)

		if currentWeight+weights[index] <= capacity {
			selected[index] = true
			knapsackRecursive(index+1, currentWeight+weights[index], currentValue+values[index], selected)
			selected[index] = false
		}
	}

	knapsackRecursive(0, 0, 0, make([]bool, n))
	return maxValue, bestSelected
}

func printSelectedItemsBruteForce(weights, values []int, capacity int) {
	result, bestSelected := knapsackBruteForce(weights, values, capacity)

	fmt.Print("Kombinasi BruteForce: {")
	first := true
	for i, selected := range bestSelected {
		if selected {
			if !first {
				fmt.Print(", ")
			}
			fmt.Print(i + 1)
			first = false
		}
	}
	fmt.Println("}")

	totalWeight := 0
	for i, selected := range bestSelected {
		if selected {
			totalWeight += weights[i]
		}
	}
	fmt.Println("Total Weight:", totalWeight)
	fmt.Println("Brute Force Profit: $", result)
}

// BRANCH N BOUND
type Node struct {
	level    int
	value    int
	weight   int
	bound    float64
	selected []bool
}

func calculateBound(u Node, n, capacity int, weights, values []int) float64 {
	profitBound := float64(u.value)
	j := u.level + 1
	totalWeight := u.weight

	for j < n && totalWeight+weights[j] <= capacity {
		totalWeight += weights[j]
		profitBound += float64(values[j])
		j++
	}

	if j < n {
		profitBound += float64(capacity-totalWeight) * float64(values[j]) / float64(weights[j])
	}

	return profitBound
}

func knapsackBranchAndBound(weights, values []int, capacity int) (int, []bool) {
	n := len(weights)
	Q := []Node{}
	initialSelected := make([]bool, n)
	u := Node{-1, 0, 0, 0, initialSelected}
	u.bound = calculateBound(u, n, capacity, weights, values)
	Q = append(Q, u)

	maxValue := 0
	bestSelected := make([]bool, n)

	for len(Q) > 0 {
		v := Q[0]
		Q = Q[1:]

		if v.bound > float64(maxValue) {
			u.level = v.level + 1

			if u.level < n {
				u.weight = v.weight + weights[u.level]
				u.value = v.value + values[u.level]
				u.selected = append([]bool(nil), v.selected...)
				u.selected[u.level] = true

				if u.weight <= capacity && u.value > maxValue {
					maxValue = u.value
					copy(bestSelected, u.selected)
				}

				u.bound = calculateBound(u, n, capacity, weights, values)

				if u.bound > float64(maxValue) {
					Q = append(Q, u)
				}

				u.weight = v.weight
				u.value = v.value
				u.selected = append([]bool(nil), v.selected...)
				u.selected[u.level] = false
				u.bound = calculateBound(u, n, capacity, weights, values)

				if u.bound > float64(maxValue) {
					Q = append(Q, u)
				}
			}
		}
	}

	return maxValue, bestSelected
}

func printSelectedItemsBranchAndBound(weights, values []int, capacity int) {
	result, selected := knapsackBranchAndBound(weights, values, capacity)

	fmt.Print("Kombinasi Branch and Bound: {")
	first := true
	for i, selectd := range selected {
		if selectd {
			if !first {
				fmt.Print(", ")
			}
			fmt.Print(i + 1)
			first = false
		}
	}
	fmt.Println("}")

	totalWeight := 0
	for i, selectd := range selected {
		if selectd {
			totalWeight += weights[i]
		}
	}
	fmt.Println("Total Weight:", totalWeight)
	fmt.Println("Branch and Bound Profit: $", result)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	file, err := os.Create("execution_times.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"N", "Brute Force (s)", "Branch and Bound (s)"})

	for i := 1; i <= 15; i++ {
		N := 10 * i
		weights := make([]int, N)
		values := make([]int, N)
		capacity := 150

		for j := 0; j < N; j++ {
			weights[j] = rand.Intn(81) + 20
			values[j] = rand.Intn(101) + 50
		}

		fmt.Println("--------------------------------------------------------------")
		fmt.Printf("| %-6s | %-20s | %-26s |\n", "Barang", "Berat (Harga Barang)", "Nilai Kebahagiaan (Profit)")
		fmt.Println("--------------------------------------------------------------")

		for k := 0; k < N; k++ {
			fmt.Printf("| %-6d | %-20d | %-26d |\n", k+1, weights[k], values[k])
		}

		fmt.Println("--------------------------------------------------------------")

		start := time.Now()
		bruteForceResult, _ := knapsackBruteForce(weights, values, capacity)
		bruteForceDuration := time.Since(start).Seconds()

		start = time.Now()
		branchAndBoundResult, _ := knapsackBranchAndBound(weights, values, capacity)
		branchAndBoundDuration := time.Since(start).Seconds()

		printSelectedItemsBruteForce(weights, values, capacity)
		fmt.Printf("Brute Force Profit: $%d\n", bruteForceResult)
		fmt.Printf("Brute Force Execution Time: %.4f s\n", bruteForceDuration)
		fmt.Println()
		printSelectedItemsBranchAndBound(weights, values, capacity)
		fmt.Printf("Branch and Bound Profit: $%d\n", branchAndBoundResult)
		fmt.Printf("Branch and Bound Execution Time: %.4f s\n", branchAndBoundDuration)

		writer.Write([]string{fmt.Sprintf("%d", N), fmt.Sprintf("%.4f", bruteForceDuration), fmt.Sprintf("%.4f", branchAndBoundDuration)})
	}
}