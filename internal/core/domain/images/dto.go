package images

type ListImagesRequest struct {
	Page  int `query:"page" validate:"required,gte=1"`
	Limit int `query:"limit" validate:"required,gte=1,lte=100"`
}

type ListImagesResponse struct {
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
	Count int     `json:"count"`
	Data  []Image `json:"data"`
}

type CreateImageRequest struct {
	ImageURL        string   `json:"image_url" validate:"required,url"`
	Transformations []string `json:"transformations" validate:"required,dive,oneof=resize crop rotate grayscale"`
}

type CreateImageResponse struct {
	ID string `json:"id"`
}

type DeleteImageRequest struct {
	ID string `validate:"required,uuid"`
}

type ProcessImageRequest struct {
	ID               string `validate:"required,uuid"`
	OriginalImageURL string `validate:"required,url"`
	StorageKey       string
	Transformations  []string
}

type UpdateImageRequest struct {
	ID                  string              `validate:"required,uuid"`
	ScaleTransformation ScaleTransformation `validate:"required"`
}
