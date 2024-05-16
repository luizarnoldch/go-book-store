package lib

import (
	"log"
	appError "main/utils/error"
	"mime"
	"strings"
)

func GetBoundaryFromMultipart(header string) (string, *appError.Error) {
	// Separate "Content-Type": {"multipart/mixed; boundary=foo"},
	// mediaType <- "multipart/mixed" params <- "boundary=foo"
	mediaType, params, err := mime.ParseMediaType(header)
	if err != nil {
		log.Printf("Error parsing media type: %v", err)
		return "", appError.NewBadRequestError("Error parsing media type")
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		log.Printf("Error media type is not multipart")
		return "", appError.NewBadRequestError("Error media type is not multipart")
	}

	return params["boundary"], nil
}
