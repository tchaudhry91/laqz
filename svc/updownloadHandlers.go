package svc

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *QServer) UploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(10 << 20)
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, handler, err := r.FormFile("asset")
		if err != nil {
			s.respond(w, r, nil, http.StatusInternalServerError, nil)
			return
		}
		s.logger.Log("msg", "uploading file", "name", handler.Filename)
		defer file.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			s.respond(w, r, nil, http.StatusInternalServerError, nil)
		}

		rand.Seed(time.Now().UnixNano())
		asset_name := randSeq(15)

		err = ioutil.WriteFile(filepath.Join(s.fileUploadDirectory, asset_name, filepath.Ext(handler.Filename)), fileBytes, 0644)
		if err != nil {
			s.respond(w, r, nil, http.StatusInternalServerError, nil)
			return
		}

		// return that we have successfully uploaded our file!
		s.respond(w, r, map[string]string{"url": fmt.Sprintf("%s/uploads/%s", s.externalURL, asset_name)}, http.StatusOK, nil)
	}
}
