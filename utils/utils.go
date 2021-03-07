package utils

func IsImageExt(extName string) bool {
	var supportExtNames = map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".ico": true, ".svg": true, ".bmp": true, ".gif": true,
	}
	return supportExtNames[extName]
}

