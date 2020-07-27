package lightweb

type H map[string]interface{}


func joinPath(absolutePath,relativePath string) string {
	 if relativePath == "" {
		return absolutePath
	 }

	 if absolutePath[0] != '/' {
	 	absolutePath = "/" + absolutePath
	 }

	if relativePath[0] != '/' {
		return absolutePath + "/" + relativePath
	}
	return absolutePath + relativePath
}
