# Spy Camera

This repo will help you to turn your PC into a surveillance camera. With this setup, you'll be able to monitor & record activities remotely. When movement is detected, it will save a video to your hard drive.

## Requisites

- Go
- GoCV
- OpenCV

## Installation

1. Install Go from [here](https://golang.org/doc/install).
2. Install OpenCV following the instructions [here](https://docs.opencv.org/master/df/d65/tutorial_table_of_content_introduction.html).
3. Install GoCV by following the instructions [here](https://gocv.io/getting-started/).

## Usage

1. Clone the repository:
    ```sh
    git clone https://github.com/jempe/spy_camera.git
    ```
2. Navigate to the repo directory:
    ```sh
    cd spy_camera
    ```
3. Run the application:
    ```sh
    go run main.go
    ```
4. Access the camera stream from your browser:
    ```
    http://localhost:4000
    ```

## Configuration

You can set the following flags when running the application:

- `-user`: Username for HTTP basic authentication (default: "username")
- `-pass`: Password for HTTP basic authentication (default: "password")
- `-h`: Show help

## Features

- Motion detection using OpenCV
- Save video when motion is detected
- HTTP basic authentication for accessing the stream

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
