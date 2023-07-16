package prompt

func allowedType(fieldType string) bool {
	switch fieldType {
	case "int":
		return true
	case "string":
		return true
	case "bool":
		return true
	case "float":
		return true
	default:
		return false
	}
}
