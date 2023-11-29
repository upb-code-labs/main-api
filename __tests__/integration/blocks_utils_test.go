package integration

import (
	"fmt"
	"net/http"
)

func UpdateMarkdownBlockContent(cookie *http.Cookie, blockUUID string, payload map[string]interface{}) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/blocks/markdown_blocks/%s/content", blockUUID)
	w, r := PrepareRequest("PUT", endpoint, payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}
