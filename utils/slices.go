package utils

func Filter[TInput any](inputArray []TInput, predicate func(TInput) bool) []TInput {
	outputArray := make([]TInput, 0)
	for i := range inputArray {
		if predicate(inputArray[i]) {
			outputArray = append(outputArray, inputArray[i])
		}
	}
	return outputArray
}

func Map[TInput, TOutput any](inputArray []TInput, mapFunction func(TInput) TOutput) []TOutput {
	outputArray := make([]TOutput, len(inputArray))
	for i := range inputArray {
		outputArray[i] = mapFunction(inputArray[i])
	}
	return outputArray
}

func Some[TInput any](inputArray []TInput, predicate func(TInput) bool) bool {
	for _, e := range inputArray {
		if predicate(e) {
			return true
		}
	}
	return false
}

func All[TInput any](inputArray []TInput, predicate func(TInput) bool) bool {
	for _, e := range inputArray {
		if !predicate(e) {
			return false
		}
	}
	return true
}
