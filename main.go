package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func getFiles(folder string) []string {
	var result []string

	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			result = append(result, (getFiles(folder + "/" + file.Name()))...)
		} else {
			result = append(result, folder+"/"+file.Name())
		}
	}

	return result
}

func loadPhrases(folder string) []string {
	var phrases []string
	for _, file := range getFiles(folder) {
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
		}

		for _, section := range strings.Split(string(buf), "\n\n") {
			items := strings.Split(strings.TrimSpace(section), "---")
			if !strings.HasPrefix(items[0], "#") {
				phrases = append(phrases, strings.TrimSpace(items[0]))
			}
		}
	}
	return phrases
}

func main() {
	phrases := loadPhrases("./text")
	router := gin.Default()
	router.StaticFile("/", "./index.html")
	router.Static("/assets", "./assets")
	router.GET("/phrase", func(c *gin.Context) {
		index := rand.Intn(len(phrases))
		c.JSON(200, gin.H{
			"phrase": phrases[index],
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
