package main

import (
	"crypto/sha256"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
	_ "gocv.io/x/gocv/contrib"
)

var webcam *gocv.VideoCapture
var stream *mjpeg.Stream

var authUsername = flag.String("user", "username", "Username for http basic authentication")
var authPassword = flag.String("pass", "password", "Password for http basic authentication")

var showHelp = flag.Bool("h", false, "Show Help")

func main() {
	flag.Parse()

	if *showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	deviceId := 0

	var webcamErr error
	webcam, webcamErr = gocv.OpenVideoCapture(deviceId)
	defer webcam.Close()

	endIfError(webcamErr)

	stream = mjpeg.NewStream()

	go mjpegCapture()

	mux := http.NewServeMux()

	mux.Handle("/", stream)
	log.Fatal(http.ListenAndServe("0.0.0.0:4000", basicAuth(mux)))

}

func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))

			expectedUsernameHash := sha256.Sum256([]byte(*authUsername))
			expectedPasswordHash := sha256.Sum256([]byte(*authPassword))

			if usernameHash == expectedUsernameHash && passwordHash == expectedPasswordHash {

				next.ServeHTTP(w, r)
				return
			} else {
				log.Println("login error")
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func mjpegCapture() {
	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			log.Println("Device Closed")
			return
		}

		if img.Empty() {
			continue
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf.GetBytes())
		buf.Close()
	}
}

func logAndExit(message string) {
	log.Println(message)
	os.Exit(1)
}

func endIfError(err error) {
	if err != nil {
		logAndExit(err.Error())
	}
}
