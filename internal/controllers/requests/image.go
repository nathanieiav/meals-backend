package requests

import (
	"fmt"
	"io"
	"mime/multipart"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/h2non/filetype"
)

type (
	CreateImage struct {
		Image *multipart.FileHeader `json:"image" form:"image" binding:"required"`
	}

	UpdateImage struct {
		ID    uuid.UUID             `json:"id" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Image *multipart.FileHeader `json:"image" form:"image" binding:"required"`
	}
)

func (req *CreateImage) IsImage() error {
	// * convert self to *multipart.Fileheader
	// * open the fileheader to read the file
	file, err := req.Image.Open()
	if err != nil {
		utlogger.LogError(err)
		return err
	}
	defer file.Close()

	// * convert the file to bytes
	filebytes, err := io.ReadAll(file)
	if err != nil {
		utlogger.LogError(err)
		return err
	}

	// * check if the file is an image
	if !filetype.IsImage(filebytes) {
		err := fmt.Errorf("file is not an image")
		utlogger.LogError(err)
		return err
	}

	return nil
}
