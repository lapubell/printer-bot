package server

import (
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func (s *server) handleUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20) // up to 32MB

		file, handler, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"error\":\"bad request\"}"))
			return
		}
		defer file.Close()

		fileName := strconv.Itoa(int(time.Now().Unix())) + handler.Filename

		localFile, err := os.Create(fileName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"error\":\"bad request\"}"))
			return
		}
		defer localFile.Close()

		io.Copy(localFile, file)

		go func() {
			time.Sleep(time.Second * 1)

			stdout, err := exec.Command("/usr/bin/lpr", "-P", "EPSON_XP_7100_Series_USB", fileName).Output()
			if err != nil {
				s.logger.Error(err)
				return
			}
			s.logger.Info(stdout)
		}()

		go func() {
			time.Sleep(time.Second * 30)
			os.Remove(fileName)
		}()

		output := map[string]string{
			"message": "yay!",
		}

		bytes, _ := json.Marshal(output)

		w.Write(bytes)
	}
}
