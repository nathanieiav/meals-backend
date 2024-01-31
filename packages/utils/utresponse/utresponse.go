package utresponse

import (
	"errors"
	"fmt"
	"net/http"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/uttoken"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

/* -------------------------------------------------------------------------- */
/*                              success responses                             */
/* -------------------------------------------------------------------------- */

type (
	SuccessRes struct {
		Status  consttypes.ResponseStatusType `json:"status" default:""`
		Message string                        `json:"message"`
		Data    any                           `json:"data,omitempty"`
		Header  uttoken.TokenHeader           `json:"-"`
	}
)

func SuccessResponse(ctx *gin.Context, code int, res SuccessRes) {
	if res.Header != (uttoken.TokenHeader{}) {
		ctx.Header("refresh-token", res.Header.RefreshToken)
		ctx.Header("refresh-token-expired", res.Header.RefreshTokenExpires.String())
		ctx.Header("Authorization", "Bearer "+res.Header.AccessToken)
		ctx.Header("expired-at", res.Header.AccessTokenExpires.String())
	}
	ctx.JSON(code, res)
}

/* -------------------------------------------------------------------------- */
/*                               error responses                              */
/* -------------------------------------------------------------------------- */

type (
	ValidationErrorMessage struct {
		Namespace string `json:"namespace"`
		Field     string `json:"field"`
		Message   string `json:"message"`
	}

	ErrorRes struct {
		Status  consttypes.ResponseStatusType `json:"status"`
		Message string                        `json:"message"`
		Data    ErrorData                     `json:"data,omitempty"`
	}

	ErrorData struct {
		Debug  error `json:"debug,omitempty"`
		Errors any   `json:"errors"`
	}
)

var (
	// General
	ErrConvertFailed = errors.New("data type conversion failed")

	// Error Field
	ErrFieldIsEmpty             = errors.New("field should not be empty")
	ErrFieldInvalidFormat       = errors.New("field format is invalid")
	ErrFieldInvalidEmailAddress = errors.New("invalid email address format")

	// Token
	ErrTokenExpired      = errors.New("token is expired")
	ErrTokenUnverifiable = errors.New("token is unverifiable")
	ErrTokenMismatch     = errors.New("token is mismatch")
	ErrTokenIsNotTheSame = errors.New("this token is not the same")

	// User
	ErrUserNotFound         = errors.New("user not found")
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrUserIDNotFound       = errors.New("unable to assert user ID")
	ErrUserAlreadyExist     = errors.New("user already exists")
	ErrUserAlreadyConfirmed = errors.New("this user is already confirmed")

	// Email
	ErrSendEmailResetRequest        = errors.New("you already requested a reset password email in less than 5 minutes")
	ErrSendEmailVerificationRequest = errors.New("you already requested a verification message in less than 5 minutes")
)

func ErrorResponse(ctx *gin.Context, code int, res ErrorRes) {
	ctx.JSON(code, res)
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "should be a valid email address"
	case "file":
		return "should be a valid file"
	case "lte":
		return "should be less than " + fe.Param()
	case "gte":
		return "should be greater than " + fe.Param()
	case "len":
		return "should be " + fe.Param() + " character(s) long"
	case "eqfield":
		return "should be equal to " + fe.Param()
	}
	return "unknown error"
}

func ValidationResponse(err error) []ValidationErrorMessage {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		out := make([]ValidationErrorMessage, len(ve))
		for i, fe := range ve {
			out[i] = ValidationErrorMessage{fe.Namespace(), fmt.Sprintf("%s", fe.Field()), getErrorMsg(fe)}
			fmt.Println(err)
		}
		return out
	}

	return nil
}

func GeneralInputRequiredError(message string, ctx *gin.Context, err any) {
	ErrorResponse(ctx, http.StatusUnprocessableEntity, ErrorRes{
		Status:  consttypes.RST_ERROR,
		Message: message,
		Data: ErrorData{
			Debug:  nil,
			Errors: err,
		},
	})
}

func GeneralInternalServerError(message string, ctx *gin.Context, err any) {
	ErrorResponse(ctx, http.StatusInternalServerError, ErrorRes{
		Status:  consttypes.RST_ERROR,
		Message: message,
		Data: ErrorData{
			Debug:  nil,
			Errors: err,
		},
	})
}

func GeneralInvalidRequest(message string, ctx *gin.Context, ve []ValidationErrorMessage, err *error) {
	ErrorResponse(ctx, http.StatusBadRequest, ErrorRes{
		Status:  consttypes.RST_ERROR,
		Message: message,
		Data: ErrorData{
			Debug:  *err,
			Errors: ve,
		},
	})
}
