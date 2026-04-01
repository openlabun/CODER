package getbysessionid

import "strings"

func MapPath(id string) PathDTO {
	id = strings.TrimSpace(id)
	if id == "" {
		return PathDTO{UserID: nil}
	}
	return PathDTO{UserID: &id}
}
