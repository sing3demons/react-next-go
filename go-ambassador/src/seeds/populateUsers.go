package seeds

import (
	"github.com/bxcodec/faker/v3"
	"github.com/sing3demons/ambassador/src/database"
	"github.com/sing3demons/ambassador/src/models"
)

func CreateAmbassador() {
	database.Connect()

	for i := 0; i < 30; i++ {
		ambassador := models.User{
			FirstName:    faker.FirstName(),
			LastName:     faker.LastName(),
			Email:        faker.Email(),
			IsAmbassador: true,
		}

		ambassador.EncryptPassword("1234")

		database.GetDB().Create(&ambassador)
	}
}
