package converter

func Error(v error) *error {
	return &v
}

func ErrorValue(v *error) error {
	if v != nil {
		return *v
	}
	return nil
}
