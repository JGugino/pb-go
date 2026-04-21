package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JGugino/pb-go/services"
	"github.com/joho/godotenv"
)

const (
	BASE_URL = "http://localhost:8090"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env")
	}

	pb := services.Pocketbase{}

	pb.Init(BASE_URL)

	_, err := pb.Auth.AuthWithPasswordForCollection("_superusers", "", "", os.Getenv("PB_IDENTITY"), os.Getenv("PB_PASSWORD"))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	records, err := pb.Record.ListRecords("testing_collection", pb.Auth.AuthToken, services.PocketBaseListOptions{
		Page:    1,
		PerPage: 30,
		Filter:  fmt.Sprintf("id='%s'", "ibopxmpxt3dap2o"),
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	updatedRecord, err := pb.Record.UpdateRecord("testing_collection", records.Items[0]["id"].(string), pb.Auth.AuthToken, map[string]any{
		"text": "This is some updated example text",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pb.Collection.ImportCollections(pb.Auth.AuthToken, []map[string]any{
		{
			"name": "example_collection",
			"type": services.BaseCollection,
			"fields": []map[string]any{
				{"text": "text"},
			},
		},
	}, true)

	fmt.Println(updatedRecord["id"])
}
