package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
	_ "gocv.io/x/gocv/contrib"
)

var webcam *gocv.VideoCapture
var stream *mjpeg.Stream

func main() {
	deviceId := 1

	var webcamErr error
	webcam, webcamErr = gocv.OpenVideoCapture(deviceId)
	defer webcam.Close()

	endIfError(webcamErr)

	stream = mjpeg.NewStream()

	go mjpegCapture()

	mux := http.NewServeMux()

	mux.Handle("/", stream)
	log.Fatal(http.ListenAndServe("0.0.0.0:4000", mux))

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
