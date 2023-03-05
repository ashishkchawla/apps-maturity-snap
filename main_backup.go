package main

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

/*type Service struct {
	ID   string
	Name string
}

var services []*Service*/

func main_bacjkup() {
	fmt.Println("welcome to go")
	graphqlClient := graphql.NewClient("https://api.opslevel.com/graphql")
	token := "X6YVXLwIRjceklWoQOrXD5RHa2VmDrJKMnnf"
	fetchServicesRequest := graphql.NewRequest(`
	{
		account {
		  services (ownerAlias: "knp_-_content_ingestion") {
			nodes {
			  name
			  id  
			}
		  }
		}
	  }
	`)
	fetchServicesRequest.Header.Set("Authorization", "Bearer "+token)
	fetchServicesRequest.Header.Set("Content-Type", "application/json")
	fetchServicesRequest.Header.Set("Accept", "application/json")

	var serviceListResponse interface{}

	/*{
		account map;


	}*/

	if err := graphqlClient.Run(context.Background(), fetchServicesRequest, &serviceListResponse); err != nil {
		panic(err)
	}

	fmt.Println(serviceListResponse)

	/*for serviceMap := range serviceListResponse["account"]["services"]["nodes"]{
		fmt.Println(serviceMap["id"])

	}*/

}
