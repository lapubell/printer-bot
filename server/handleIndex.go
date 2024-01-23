package server

import (
	_ "embed"
	"encoding/base64"
	"net/http"
	"strings"
)

//go:embed assets/index.html
var indexHTML string

//go:embed assets/style.css
var stylesheet string

//go:embed assets/printer.jpg
var printerBot []byte

func (s *server) handleIndex() http.HandlerFunc {
	printerBotByteString := base64.StdEncoding.EncodeToString(printerBot)

	return func(w http.ResponseWriter, r *http.Request) {
		indexHTML := strings.ReplaceAll(indexHTML, "%%PRINTERBOT%%", printerBotByteString)
		indexHTML = strings.ReplaceAll(indexHTML, "%%CSS%%", "<style>"+stylesheet+"</style>")

		w.Write([]byte(indexHTML))
	}
}
