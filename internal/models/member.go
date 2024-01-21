package models

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	Member struct {
		helper.Model
		UserID         uuid.UUID         `json:"user_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User           User              `json:"user"`
		CaregiverID    *uuid.UUID        `json:"caregiver_id,omitempty" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4" default:"null"`
		Caregiver      *Caregiver        `json:"caregiver,omitempty"`
		OrganizationID *uuid.UUID        `json:"organization_id,omitempty" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4" default:"null"`
		Organization   *Organization     `json:"organization,omitempty"`
		Illness        []*MemberIllness  `json:"illness,omitempty"`
		Allergy        []*MemberAllergy  `json:"allergy,omitempty"`
		Height         float64           `json:"height" gorm:"not null" binding:"required" example:"100"`
		Weight         float64           `json:"weight" gorm:"not null" binding:"required" example:"150"`
		BMI            float64           `json:"bmi" gorm:"not null;type:decimal(10,2)" binding:"required" example:"19"`
		FirstName      string            `json:"first_name" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName       string            `json:"last_name" gorm:"not null" binding:"required" example:"Vince"`
		Gender         consttypes.Gender `json:"gender" gorm:"not null; type:gender_enum" binding:"required" example:"Male"`
		DateOfBirth    customs.CDT_DATE  `json:"date_of_birth" gorm:"not null" binding:"required" example:"2000-10-20"`
	}

	MemberIllness struct {
		helper.Model
		MemberID  uuid.UUID `json:"member_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		IllnessID uuid.UUID `json:"illness_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Illness   Illness   `json:"illness"`
	}

	MemberAllergy struct {
		helper.Model
		MemberID  uuid.UUID `json:"member_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		AllergyID uuid.UUID `json:"allergy_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Allergy   Allergy   `json:"allergy"`
	}
)

func (m *Member) ToResponse() *responses.MemberResponse {
	mres := responses.MemberResponse{}

	if err := copier.Copy(&mres, &m); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &mres
}
