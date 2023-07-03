package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fossology/LicenseDb/pkg/db"
	"github.com/fossology/LicenseDb/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	db.Connect()
	exitcode := m.Run()
	os.Exit(exitcode)
}

func makeRequest(method, path string, body interface{}, isAuthanticated bool) *httptest.ResponseRecorder {
	reqBody, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if isAuthanticated {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("avinal:avinal")))
	}
	w := httptest.NewRecorder()
	Router().ServeHTTP(w, req)
	return w
}
func TestGetSingleLicense(t *testing.T) {
	expectLicense := models.License{
		Shortname:     "MIT",
		Fullname:      "MIT License",
		Text:          "MIT License\n\nCopyright (c) <year> <copyright holders>\n\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the \"Software\"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.\n",
		Url:           "https://opensource.org/licenses/MIT",
		TextUpdatable: "f",
		DetectorType:  "1",
		Active:        "t",
		Flag:          "1",
		Marydone:      "f",
	}
	w := makeRequest("GET", "/api/license/MIT", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var res models.LicenseResponse
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, expectLicense, res.Data[0])

}
