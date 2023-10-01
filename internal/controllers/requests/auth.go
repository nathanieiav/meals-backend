package requests

type (
	LoginRequest struct {
		Email    string `json:"email" binding:"required,email" example:"email@email.com"`
		Password string `json:"password" binding:"required" example:"password"`
	}

	RegisterRequest struct {
		Email    string `json:"email" binding:"required,email" example:"email@email.com"`
		Password string `json:"password" binding:"required" example:"password"`
	}

	VerifyTokenRequest struct {
		Token int    `json:"token" binding:"required,number" example:""`
		Email string `json:"email" binding:"required,email"`
	}

	ForgotPasswordRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	ResetPasswordRequest struct {
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
		Token           int    `json:"token" binding:"required"`
	}

	ResetPasswordRedirectRequest struct {
		ResetToken string `uri:"token" binding:"required"`
	}
)
