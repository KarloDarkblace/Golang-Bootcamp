package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func ReadNumbers() []int {
	var numbers []int
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()

		if text == "x" {
			break
		}

		number, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Пожалуйста, введите число или 'x' для завершения программы.")
			continue
		}

		numbers = append(numbers, number)
	}

	return numbers
}

func CountMean(numbers []int) float64 {
	if len(numbers) == 0 {
		return 0
	}
	var mean float64 = 0.0

	for _, number := range numbers {
		mean += float64(number)
	}

	mean /= float64(len(numbers))

	return mean
}

func CountMedian(numbers []int) float64 {
	n := len(numbers)
	if n == 0 {
		return 0.0
	}

	sort.Ints(numbers)

	var median float64 = 0.0
	if (n % 2) == 1 {
		median = float64(numbers[n/2])
	} else {
		left_num := numbers[(n/2)-1]
		right_num := numbers[n/2]
		median = float64(left_num+right_num) / 2.0
	}
	return median
}

func CountMode(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}

	sort.Ints(numbers)

	mode := numbers[0]
	maxCount := 1
	currentCount := 1

	for i := 1; i < len(numbers); i++ {
		if numbers[i] == numbers[i-1] {
			currentCount++
		} else {
			currentCount = 1
		}

		if currentCount > maxCount {
			maxCount = currentCount
			mode = numbers[i]
		}
	}

	return mode
}

func CountSD(numbers []int) float64 {
	if len(numbers) == 0 {
		return 0
	}

	mean := CountMean(numbers)
	var sum float64 = 0.0
	n := float64(len(numbers))

	for _, number := range numbers {
		sum += math.Pow(float64(number)-mean, 2)
	}

	sd := math.Sqrt(sum / n)
	return sd
}

func main() {
	fmt.Println("Для завершения ввода напишите 'x'")
	numbers := ReadNumbers()

	mean := math.Round(CountMean(numbers)*100) / 100
	median := math.Round(CountMedian(numbers)*100) / 100
	mode := CountMode(numbers)
	sd := math.Round(CountSD(numbers)*100) / 100

	fmt.Println("Mean:", mean)
	fmt.Println("Median:", median)
	fmt.Println("Mode:", mode)
	fmt.Println("SD:", sd)
}
