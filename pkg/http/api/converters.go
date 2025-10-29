package api

import (
	"fmt"

	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
)

// ConvertAPITransformationsToDomain converts API TransformationRequest to domain TransformationRequest
func ConvertAPITransformationsToDomain(apiTransformations []TransformationRequest) ([]images.TransformationRequest, error) {
	domainTransformations := make([]images.TransformationRequest, 0, len(apiTransformations))

	for i, apiTrans := range apiTransformations {
		discriminator, err := apiTrans.Discriminator()
		if err != nil {
			return nil, fmt.Errorf("failed to get discriminator for transformation %d: %w", i, err)
		}

		config := make(map[string]interface{})

		switch discriminator {
		case "resize":
			resize, err := apiTrans.AsResizeTransformation()
			if err != nil {
				return nil, fmt.Errorf("failed to parse resize transformation %d: %w", i, err)
			}
			config["width"] = resize.Config.Width
			config["height"] = resize.Config.Height

		case "grayscale":
			// Grayscale has no config
			// Empty config map is fine

		case "trim":
			trim, err := apiTrans.AsTrimTransformation()
			if err != nil {
				return nil, fmt.Errorf("failed to parse trim transformation %d: %w", i, err)
			}
			config["threshold"] = trim.Config.Threshold

		case "blur":
			blur, err := apiTrans.AsBlurTransformation()
			if err != nil {
				return nil, fmt.Errorf("failed to parse blur transformation %d: %w", i, err)
			}
			config["sigma"] = blur.Config.Sigma

		case "rotate":
			rotate, err := apiTrans.AsRotateTransformation()
			if err != nil {
				return nil, fmt.Errorf("failed to parse rotate transformation %d: %w", i, err)
			}
			config["angle"] = int(rotate.Config.Angle)

		default:
			return nil, fmt.Errorf("unknown transformation type: %s", discriminator)
		}

		domainTransformations = append(domainTransformations, images.TransformationRequest{
			Name:   discriminator,
			Config: config,
		})
	}

	return domainTransformations, nil
}

// ConvertDomainTransformationsToAPI converts domain TransformationRequest to API TransformationRequest
func ConvertDomainTransformationsToAPI(domainTransformations []images.TransformationRequest) ([]TransformationRequest, error) {
	apiTransformations := make([]TransformationRequest, 0, len(domainTransformations))

	for i, domainTrans := range domainTransformations {
		var apiTrans TransformationRequest

		switch domainTrans.Name {
		case "resize":
			width, ok := domainTrans.Config["width"].(int)
			if !ok {
				if widthFloat, ok := domainTrans.Config["width"].(float64); ok {
					width = int(widthFloat)
				} else {
					return nil, fmt.Errorf("resize transformation %d: missing or invalid width", i)
				}
			}

			height, ok := domainTrans.Config["height"].(int)
			if !ok {
				if heightFloat, ok := domainTrans.Config["height"].(float64); ok {
					height = int(heightFloat)
				} else {
					return nil, fmt.Errorf("resize transformation %d: missing or invalid height", i)
				}
			}

			err := apiTrans.FromResizeTransformation(ResizeTransformation{
				Name: ResizeTransformationNameResize,
				Config: ResizeConfig{
					Width:  width,
					Height: height,
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create resize transformation %d: %w", i, err)
			}

		case "grayscale":
			err := apiTrans.FromGrayscaleTransformation(GrayscaleTransformation{
				Name:   GrayscaleTransformationNameGrayscale,
				Config: GrayscaleConfig{},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create grayscale transformation %d: %w", i, err)
			}

		case "trim":
			threshold, ok := domainTrans.Config["threshold"].(float64)
			if !ok {
				return nil, fmt.Errorf("trim transformation %d: missing or invalid threshold", i)
			}

			err := apiTrans.FromTrimTransformation(TrimTransformation{
				Name: TrimTransformationNameTrim,
				Config: TrimConfig{
					Threshold: threshold,
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create trim transformation %d: %w", i, err)
			}

		case "blur":
			sigma, ok := domainTrans.Config["sigma"].(float64)
			if !ok {
				return nil, fmt.Errorf("blur transformation %d: missing or invalid sigma", i)
			}

			err := apiTrans.FromBlurTransformation(BlurTransformation{
				Name: BlurTransformationNameBlur,
				Config: BlurConfig{
					Sigma: sigma,
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create blur transformation %d: %w", i, err)
			}

		case "rotate":
			angle, ok := domainTrans.Config["angle"].(int)
			if !ok {
				if angleFloat, ok := domainTrans.Config["angle"].(float64); ok {
					angle = int(angleFloat)
				} else {
					return nil, fmt.Errorf("rotate transformation %d: missing or invalid angle", i)
				}
			}

			err := apiTrans.FromRotateTransformation(RotateTransformation{
				Name: RotateTransformationNameRotate,
				Config: RotateConfig{
					Angle: RotateConfigAngle(angle),
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create rotate transformation %d: %w", i, err)
			}

		default:
			return nil, fmt.Errorf("unknown transformation type: %s", domainTrans.Name)
		}

		apiTransformations = append(apiTransformations, apiTrans)
	}

	return apiTransformations, nil
}

// ConvertDomainImageToAPI converts domain Image to API Image
func ConvertDomainImageToAPI(domainImage *images.Image) (*Image, error) {
	apiTransformations, err := ConvertDomainTransformationsToAPI(domainImage.Transformations)
	if err != nil {
		return nil, err
	}

	apiImage := &Image{
		Id:                    &domainImage.ID,
		OriginalImageUrl:      &domainImage.OriginalImageURL,
		ObjectStorageImageKey: &domainImage.ObjectStorageImageKey,
		TransformedImageKey:   &domainImage.TransformedImageKey,
		MimeType:              &domainImage.MimeType,
		Checksum:              &domainImage.Checksum,
		Status:                &domainImage.Status,
		Transformations:       &apiTransformations,
		CreatedAt:             &domainImage.CreatedAt,
		UpdatedAt:             &domainImage.UpdatedAt,
		ErrorMessage: &domainImage.ErrorMessage,
	}

	return apiImage, nil
}

// ConvertAPICreateImageRequestToDomain converts API CreateImageRequest to domain CreateImageRequest
func ConvertAPICreateImageRequestToDomain(apiReq *CreateImageRequest) (*images.CreateImageRequest, error) {
	transformations, err := ConvertAPITransformationsToDomain(apiReq.Transformations)
	if err != nil {
		return nil, err
	}

	return &images.CreateImageRequest{
		ImageURL:        apiReq.ImageUrl,
		Transformations: transformations,
	}, nil
}

// ConvertAPIUpdateImageRequestToDomain converts API UpdateImageRequest to domain UpdateImageRequest
func ConvertAPIUpdateImageRequestToDomain(apiReq *UpdateImageRequest) (*images.UpdateImageRequest, error) {
	transformations, err := ConvertAPITransformationsToDomain(apiReq.Transformations)
	if err != nil {
		return nil, err
	}

	return &images.UpdateImageRequest{
		ID:              apiReq.Id.String(),
		Transformations: transformations,
	}, nil
}
