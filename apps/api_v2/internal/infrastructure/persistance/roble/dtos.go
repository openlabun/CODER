package roble_infrastructure

type RobleLoginResponse struct {
	AccessToken  string            `json:"accessToken"`
	RefreshToken string            `json:"refreshToken"`
	User         *RobleUserPayload `json:"user,omitempty"`
}

type RobleRefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RobleUserPayload struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	AvatarURL   *string `json:"avatarUrl"`
	PhoneNumber *string `json:"phoneNumber"`
	BirthDate   *string `json:"birthDate"`
	Gender      *string `json:"gender"`
	Country     *string `json:"country"`
	City        *string `json:"city"`
	Address     *string `json:"address"`
	Role        string  `json:"role"`
}

type RobleVerifyTokenResponse struct {
	Valid bool `json:"valid"`
}

type RobleMessageResponse struct {
	Message string `json:"message"`
}

type RobleInsertRequest struct {
	TableName string           `json:"tableName"`
	Records   []map[string]any `json:"records"`
}

type RobleUpdateRequest struct {
	TableName string         `json:"tableName"`
	IDColumn  string         `json:"idColumn"`
	IDValue   string         `json:"idValue"`
	Updates   map[string]any `json:"updates"`
}

type RobleDeleteRequest struct {
	TableName string `json:"tableName"`
	IDColumn  string `json:"idColumn"`
	IDValue   string `json:"idValue"`
}
