package getpublic

func MapQuery(tag, difficulty string) QueryDTO {
	var t *string
	var d *string
	if tag != "" {
		t = &tag
	}
	if difficulty != "" {
		d = &difficulty
	}
	return QueryDTO{Tag: t, Difficulty: d}
}
