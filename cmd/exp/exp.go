package main

import (
	"fmt"

	"github.com/SuWh1/WebDevGo/models"
)

func main() {
	gs := models.GalleryService{}
	fmt.Println(gs.Images(2))
}
