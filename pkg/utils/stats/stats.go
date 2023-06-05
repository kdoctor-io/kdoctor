// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package stats

import (
	"errors"
	"sort"
)

var (
	EmptyInputErr = errors.New("Input must not be empty.")
	BoundsErr     = errors.New("Input is outside of range.")
)

// Float32Data is a named type for []float32 with helper methods
type Float32Data []float32

func (f Float32Data) Get(i int) float32 { return f[i] }

// Len returns length of slice
func (f Float32Data) Len() int { return len(f) }

// Less returns if one number is less than another
func (f Float32Data) Less(i, j int) bool { return f[i] < f[j] }

// Swap switches out two numbers in slice
func (f Float32Data) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

// Min returns the minimum number in the data
func (f Float32Data) Min() (float32, error) { return Min(f) }

// Max returns the maximum number in the data
func (f Float32Data) Max() (float32, error) { return Max(f) }

// Sum returns the total of all the numbers in the data
func (f Float32Data) Sum() (float32, error) { return Sum(f) }

// CumulativeSum returns the cumulative sum of the data
func (f Float32Data) CumulativeSum() ([]float32, error) { return CumulativeSum(f) }

// Mean returns the mean of the data
func (f Float32Data) Mean() (float32, error) { return Mean(f) }

// Max finds the highest number in a slice
func Max(input Float32Data) (max float32, err error) {

	// Return an error if there are no numbers
	if input.Len() == 0 {
		return 0, EmptyInputErr
	}

	// Get the first value as the starting point
	max = input.Get(0)

	// Loop and replace higher values
	for i := 1; i < input.Len(); i++ {
		if input.Get(i) > max {
			max = input.Get(i)
		}
	}

	return max, nil
}

// Min finds the lowest number in a set of data
func Min(input Float32Data) (min float32, err error) {

	// Get the count of numbers in the slice
	l := input.Len()

	// Return an error if there are no numbers
	if l == 0 {
		return 0, EmptyInputErr
	}

	// Get the first value as the starting point
	min = input.Get(0)

	// Iterate until done checking for a lower value
	for i := 1; i < l; i++ {
		if input.Get(i) < min {
			min = input.Get(i)
		}
	}
	return min, nil
}

// Sum adds all the numbers of a slice together
func Sum(input Float32Data) (sum float32, err error) {

	if input.Len() == 0 {
		return 0, EmptyInputErr
	}

	// Add em up
	for _, n := range input {
		sum += n
	}

	return sum, nil
}

// Mean gets the average of a slice of numbers
func Mean(input Float32Data) (float32, error) {

	if input.Len() == 0 {
		return 0, EmptyInputErr
	}

	sum, _ := input.Sum()

	return sum / float32(input.Len()), nil
}

// CumulativeSum calculates the cumulative sum of the input slice
func CumulativeSum(input Float32Data) ([]float32, error) {

	if input.Len() == 0 {
		return Float32Data{}, EmptyInputErr
	}

	cumSum := make([]float32, input.Len())

	for i, val := range input {
		if i == 0 {
			cumSum[i] = val
		} else {
			cumSum[i] = cumSum[i-1] + val
		}
	}

	return cumSum, nil
}

// Percentile finds the relative standing in a slice of floats
func Percentile(input Float32Data, percent float32) (percentile float32, err error) {
	length := input.Len()
	if length == 0 {
		return 0, EmptyInputErr
	}

	if length == 1 {
		return input[0], nil
	}

	if percent <= 0 || percent > 100 {
		return 0, BoundsErr
	}

	// Start by sorting a copy of the slice
	sort.Sort(input)

	// Multiply percent by length of input
	index := (percent / 100) * float32(len(input))

	// Check if the index is a whole number
	if index == float32(int64(index)) {

		// Convert float to int
		i := int(index)

		// Find the value at the index
		percentile = input[i-1]

	} else if index > 1 {

		// Convert float to int via truncation
		i := int(index)

		// Find the average of the index and following values
		percentile, _ = Mean(Float32Data{input[i-1], input[i]})

	} else {
		return 0, BoundsErr
	}

	return percentile, nil

}
