package getme

import "strings"

func ResolveUserIDFromRequest(queryUserID, headerUserEmail string) string {
	if strings.TrimSpace(queryUserID) != "" {
		return strings.TrimSpace(queryUserID)
	}
	return strings.TrimSpace(headerUserEmail)
}
