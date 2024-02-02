package utils

func Map[TInput, TOutput any](inputArray []TInput, mapFunction func(TInput) TOutput) []TOutput {
	outputArray := make([]TOutput, len(inputArray))
	for i := range inputArray {
		outputArray[i] = mapFunction(inputArray[i])
	}
	return outputArray
}
