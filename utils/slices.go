// Package utils provides utility functions for general-purpose functionality.
package utils

// Filter filters the elements of an array based on a predicate.
func Filter[TInput any](inputArray []TInput, predicate func(TInput) bool) []TInput {
	outputArray := make([]TInput, 0)
	for i := range inputArray {
		if predicate(inputArray[i]) {
			outputArray = append(outputArray, inputArray[i])
		}
	}
	return outputArray
}

// Map maps the elements of an array to another array based on a mapping function.
func Map[TInput, TOutput any](inputArray []TInput, mapFunction func(TInput) TOutput) []TOutput {
	outputArray := make([]TOutput, len(inputArray))
	for i := range inputArray {
		outputArray[i] = mapFunction(inputArray[i])
	}
	return outputArray
}

// FlatMap maps the elements of an array to another array based on a mapping function that returns a slice.
func FlatMap[TInput, TOutput any](inputArray []TInput, mapFunction func(TInput) []TOutput) []TOutput {
	outputArray := make([]TOutput, 0)
	for i := range inputArray {
		outputArray = append(outputArray, mapFunction(inputArray[i])...)
	}
	return outputArray
}

// Flatten flattens a 2D array to a 1D array.
func Flatten[T any](inputArray [][]T) []T {
	outputArray := make([]T, 0)
	for i := range inputArray {
		outputArray = append(outputArray, inputArray[i]...)
	}
	return outputArray
}

// Some checks if any element of an array satisfies a predicate.
func Some[TInput any](inputArray []TInput, predicate func(TInput) bool) bool {
	for _, e := range inputArray {
		if predicate(e) {
			return true
		}
	}
	return false
}

// All checks if all elements of an array satisfy a predicate.
func All[TInput any](inputArray []TInput, predicate func(TInput) bool) bool {
	for _, e := range inputArray {
		if !predicate(e) {
			return false
		}
	}
	return true
}
