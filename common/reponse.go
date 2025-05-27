package common

type ResponseNormal struct {
	Message string      `json:"message"`
	Email   string      `json:"email,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Result  bool 	  	`json:"result,omitempty"` 
	Paging  interface{}	`json:"paging,omitempty"`
}

type ResponseLogin struct {
	UserID 	  	 string `json:"user_id"`
	Role		 string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ResponseAccessToken struct {
	AccessToken string `json:"access_token"`
}

func NewResponseNormal(message string, data interface{}) *ResponseNormal {
	return &ResponseNormal{
		Message: message,
		Data:    data,
	}
}

func NewResponseLogin(userID, role, accessToken, refreshToken string) *ResponseLogin {
	return &ResponseLogin{
		UserID:       userID,
		Role:         role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func NewResponseRegister(message, email string) *ResponseNormal {
	return &ResponseNormal{
		Message: message,
		Email:   email,
	}
}

func NewResponseResult(message string, result bool) *ResponseNormal {
	return &ResponseNormal{
		Message: message,
		Result:  result,
	}
}


func NewResponseForgotPassword(message, email string, result bool) *ResponseNormal {
	return &ResponseNormal{
		Message: message,
		Email:   email,
		Result:  result,
	}
}

func NewResponseAccessToken(accessToken string) *ResponseAccessToken {
	return &ResponseAccessToken{
		AccessToken: accessToken,
	}
}

func NewResponsePaging(message string, data interface{}, paging interface{}) *ResponseNormal {
	return &ResponseNormal{
		Message: message,
		Data:    data,
		Paging:  paging,
	}
}