package seeds

import (
	"math/rand"

	"github.com/bxcodec/faker/v3"
	"github.com/sing3demons/ambassador/src/database"
	"github.com/sing3demons/ambassador/src/models"
)

func CreateProduct() {
	database.Connect()

	for i := 0; i < 30; i++ {
		product := models.Product{
			Title:       faker.Username(),
			Description: faker.Username(),
			Image:       faker.URL(),
			Price:       float64(rand.Intn(90) + 10),
		}

		database.GetDB().Create(&product)
	}

}
