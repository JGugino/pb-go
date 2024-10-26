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

		log.Println("Get admin record")
		getRes, err := adminApi.GetList(0, 0, "-created", "", "", false, res.Token)

		if ok := services.HandlePocketBaseError(err); !ok {
			return
		}

		log.Println(getRes)
	}
}
