package utils

import "errors"


const MaxUploadSize = 5 << 20 // 5MB

var AllowedCommentExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}
var ErrInvalidImageType = errors.New("only JPEG, PNG, and GIF images are allowed")