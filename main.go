package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
	_ "gocv.io/x/gocv/contrib"
)

var err error

var webcam *gocv.VideoCapture
var stream *mjpeg.Stream

var authUsername = flag.String("user", "username", "Username for http basic authentication")
var authPassword = flag.String("pass", "password", "Password for http basic authentication")

var captureThreshold float64 = 3.0
var minCaptureFrames int = 150

var saveFolder = "/home/jempe/"

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

	lastHash := gocv.NewMat()

	averageHash := contrib.AverageHash{}

	framesToCapture := 0

	saveFile := ""

	var writer *gocv.VideoWriter

	defer writer.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			log.Println("Device Closed")
			return
		}

		if img.Empty() {
			continue
		}

		currentTime := time.Now()

		textColor := color.RGBA{0, 255, 255, 0}

		currentHash := gocv.NewMat()

		averageHash.Compute(img, &currentHash)

		if !currentHash.Empty() {

			if lastHash.Empty() {
				averageHash.Compute(img, &lastHash)
			}

			similarity := averageHash.Compare(lastHash, currentHash)

			gocv.PutText(&img, currentTime.Format("2006-01-02 15:04:05"), image.Pt(10, 35), gocv.FontHersheySimplex, 1.4, textColor, 3)

			if similarity > captureThreshold {
				lastHash = currentHash

				log.Println("Motion Detected")

				if framesToCapture < minCaptureFrames {
					framesToCapture = minCaptureFrames
				}

				if saveFile == "" {
					saveFile = saveFolder + "video" + currentTime.Format("20060102_150405") + ".avi"

					log.Println("Save video to", saveFile)

					writer, err = gocv.VideoWriterFile(saveFile, "MJPG", 25, img.Cols(), img.Rows(), true)
					if err != nil {
						fmt.Printf("error opening video writer device: %v\n", saveFile)
						return
					}
				}
			}

			if framesToCapture > 0 {
				writer.Write(img)

				if framesToCapture == 1 {
					saveFile = ""
				}

				framesToCapture--
			}
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
