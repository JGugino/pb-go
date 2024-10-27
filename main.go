package main

import (
	"log"

	"github.com/JGugino/pb-go/services"
)

func main() {
	baseUrl := "http://localhost:8090"

	healthApi := services.HealthAPI{BaseURL: baseUrl}
	adminApi := services.AdminsAPI{BaseURL: baseUrl}

	health := healthApi.CheckAPIHealth()

	if health.Code == 200 {
		log.Println(health.Message)

		log.Println("Authenticating admin with password")
		res, err := adminApi.AuthWithPassword(services.AdminAuthQuery{Email: "gugino.joshua@gmail.com", Password: "Password123", Fields: ""})

		if ok := services.HandlePocketBaseError(err); !ok {
			return
		}

		log.Println("Create admin")
		createRes, err := adminApi.CreateAdmin("", "support@qrify.info", "Password!234", "Password!234", 3, res.Token, "")
		if ok := services.HandlePocketBaseError(err); !ok {
			return
		}

		log.Println(createRes)
	}
}
