package utils

func PtrToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
