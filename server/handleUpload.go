package server

import (
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
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

		fileName := "/tmp/" + strconv.Itoa(int(time.Now().Unix())) + "." + handler.Filename

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
			s.logger.Info("gonna try and print " + fileName)

			mimeType, err := mimetype.DetectFile(fileName)
			if err != nil {
				s.logger.Error("bummer! " + err.Error())
				return
			}

			if mimeType.Is("image/jpeg") {
				s.logger.Info("Now printing an image!")

				// check if needs rotation
				stdout, err := exec.Command("/usr/bin/identify", "-verbose", fileName).Output()
				if err != nil {
					s.logger.Error("bad info: " + err.Error())
					return
				}
				lines := strings.Split(string(stdout), "\n")
				geometryLine := ""

				for _, line := range lines {
					if len(line) < 13 {
						continue
					}

					if line[0:12] == "  Geometry: " {
						geometryLine = strings.Replace(line, "  Geometry: ", "", 1)
						geometryLine = strings.Replace(geometryLine, "+0+0", "", 1)
						break
					}
				}

				dimentions := strings.Split(geometryLine, "x")
				if len(dimentions) != 2 {
					s.logger.Error("bad dimentions!", dimentions, geometryLine)
					return
				}
				s.logger.Info("Dimentions: ", dimentions)

				width, err := strconv.Atoi(dimentions[0])
				if err != nil {
					s.logger.Error("width ", err.Error())
					return
				}
				height, err := strconv.Atoi(dimentions[1])
				if err != nil {
					s.logger.Error("height ", err.Error())
					return
				}

				if width > height {
					s.logger.Info("wider than tall, ROTATE!")
					exec.Command("/usr/bin/convert", fileName, "-rotate", "-90", fileName).Output()
				}

				// force the image to 4x6
				resizedFilename := strings.ReplaceAll(fileName, "/tmp/", "/tmp/resized-")
				croppedFilename := strings.ReplaceAll(fileName, "/tmp/", "/tmp/cropped-")
				printerFilename := strings.ReplaceAll(fileName, "/tmp/", "/tmp/print-")
				printerFilename = strings.ReplaceAll(printerFilename, ".jpg", ".pdf")

				// resize file and delete og image
				s.logger.Info("resizing...", fileName, resizedFilename)
				exec.Command("/usr/bin/convert", fileName, "-resize", "1920x2880^", resizedFilename).Output()

				s.logger.Info("cropping...", resizedFilename, croppedFilename)
				// crop file and delete resized image
				exec.Command("/usr/bin/convert", resizedFilename, "-gravity", "center", "-crop", "1920x2880+0+0", "+repage", croppedFilename).Output()

				// convert to pdf and delete cropped image
				s.logger.Info("converting...", printerFilename, croppedFilename)
				exec.Command("/usr/local/bin/pdfcpu", "import", printerFilename, croppedFilename).Output()

				// print this thing!
				s.logger.Info("printing!" + printerFilename)
				stdout, err = exec.Command("/usr/bin/lpr", "-P", "EPSON_XP_7100_Series_USB", "-o", "InputSlot=Photo", "-o", "PageSize=4x6.Borderless", printerFilename).Output()
				if err != nil {
					s.logger.Error(err)
					return
				}
				s.logger.Info(stdout)
				return
			}

			if mimeType.Is("application/pdf") {
				s.logger.Info("Now printing a PDF!")
				stdout, err := exec.Command("/usr/bin/lpr", "-P", "EPSON_XP_7100_Series_USB", fileName).Output()
				if err != nil {
					s.logger.Error(err)
					return
				}
				s.logger.Info(stdout)
				return
			}

			s.logger.Error("not printing anything, invalid mime type :(", mimeType)
		}()

		// go func() {
		// 	time.Sleep(time.Second * 30)
		// 	os.Remove(fileName)
		// }()

		output := map[string]string{
			"message": "yay!",
		}

		outputBytes, _ := json.Marshal(output)

		w.Write(outputBytes)
	}
}
