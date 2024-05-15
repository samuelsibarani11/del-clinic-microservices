package utils

import (
	"fmt"
	"path/filepath"
)

func GenerateImageFile(filename string, ext string) string {
	typ := filepath.Ext(ext)
	filenames := fmt.Sprintf("Image_%s%s", filename, typ)
	return filenames
}
