package utils

func BuildErrors(errs []error) []string {
	result := make([]string, 0)
	for _, el := range errs {
		result = append(result, el.Error())
	}
	return result
}
