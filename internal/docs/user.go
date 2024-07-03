package docs

type SignupResponse struct {
	// The status of the request.
	//
	// Required: true
	// Example: false
	Ok bool
	// The response message.
	//
	// Required: true
	// Example: signed up successfully
	Message string
	// The response data.
	//
	// Required: true
	// Example: success
	Data string
}

//	error response for signup requests.
//
// swagger:response signupErrorResponse
type SignupErrorResponse struct {
	// in:body
	Ok bool `json:"ok"`
	// in:body
	Message string `json:"message"`
	// in:body
}

//	Request paramerters for signup requests.
//
// swagger:parameters signupRequest
type SignupRequest struct {
	// in:body
	Email string `json:"email" validate:"required,email"`
	// in:body
	Password string `json:"password" validate:"required,min=8,max=20"`
}

// Response for signup request
// swagger:response signupSuccessResponse
type SignupSuccessWrapper struct {
	// The success response
	// in:body
	Body SignupResponse
}
