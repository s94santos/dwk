package routers

import (
	"github.com/gin-gonic/gin"
	"os"
	"log"
	"time"
	"net/http"
	"io"

)

const (
	imageURL  = "https://picsum.photos/1200"
	cacheDir  = "cache"
	imagePath = cacheDir + "/image.jpg"
	ttl       = 10 * time.Minute
)

func SetRoutes() *gin.Engine {

	r := gin.Default()

	r.GET("/", landingPage)
	r.GET("/image", imageHandler)

	r.LoadHTMLGlob("public/*")

	return r

}


func landingPage(c *gin.Context) {
	
	handler()
	
	todos := []string{
			"Learn Go",
			"Build a web app",
			"Deploy with Docker",
	}

	c.HTML(200, "index.html", gin.H{
		"todos": todos,
	})
}

func handler() {
	if needsRefresh() {
		log.Println("Downloading new image")
		_ = downloadImage()

	}
}

func imageHandler(c *gin.Context) {
	if needsRefresh() {
		if err := downloadImage(); err != nil {
			c.String(http.StatusInternalServerError, "Failed to load image")
			return
		}
	}

	c.File(imagePath)
}

func needsRefresh() bool {
	info, err := os.Stat(imagePath)
	if err != nil {
		return true // file does not exist
	}
	return time.Since(info.ModTime()) > ttl
}

func downloadImage() error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
