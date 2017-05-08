package s3manager_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mastertinner/s3manager"
	"github.com/stretchr/testify/assert"
)

func TestGetObjectHandler(t *testing.T) {
	assert := assert.New(t)

	cases := map[string]struct {
		s3                   s3manager.S3
		bucketName           string
		objectName           string
		expectedStatusCode   int
		expectedBodyContains string
	}{
		"s3 error": {
			s3: &s3Mock{
				Err: errors.New("mocked S3 error"),
			},
			bucketName:           "testBucket",
			objectName:           "testObject",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: http.StatusText(http.StatusInternalServerError),
		},
	}

	for tcID, tc := range cases {
		r := mux.NewRouter()
		r.
			Methods(http.MethodGet).
			Path("/buckets/{bucketName}/objects/{objectName}").
			Handler(s3manager.GetObjectHandler(tc.s3))

		ts := httptest.NewServer(r)
		defer ts.Close()

		url := fmt.Sprintf("%s/buckets/%s/objects/%s", ts.URL, tc.bucketName, tc.objectName)
		resp, err := http.Get(url)
		assert.NoError(err, tcID)
		defer func() {
			err = resp.Body.Close()
			assert.NoError(err, tcID)
		}()

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(err, tcID)

		assert.Equal(tc.expectedStatusCode, resp.StatusCode, tcID)
		assert.Contains(string(body), tc.expectedBodyContains, tcID)
	}
}
