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
