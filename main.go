package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type maturityReport struct {
	CategoryBreakdown  []categoryBreakdown  `json:"categoryBreakdown"`
	OverallLevel       overallLevel         `json:"overallLevel"`
	LatestCheckResults []latestCheckResults `json:"latestCheckResults"`
}
type categoryBreakdown struct {
	Category category `json:"category"`
	Level    level    `json:"level"`
}
type category struct {
	Name string `json:"name"`
}

type level struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Alias       string `json:"alias"`
}

type overallLevel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Alias       string `json:"alias"`
}
type latestCheckResults struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type opsLevelService struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	//MaturityReport1 MaturityReport
}

type services struct {
	Nodes []opsLevelService `json:"nodes"`
}

type service struct {
	MaturityReport maturityReport `json:"maturityReport"`
}

type account struct {
	Services services `json:"services"`
	Service  service  `json:"service"`
}

type data struct {
	Account account `json:"account"`
}

type ServicesResponse struct {
	Data data `json:"data"`
}

// Entity

type serviceReportEntity struct {
	ID                   string
	Name                 string
	Description          string
	SecurityLevel        string
	ResiliencyLevel      string
	InfrastructureLevel  string
	QualityLevel         string
	AppArchitectureLevel string
	UnitTestScore        string
	overallLevel         string
	LastUpdated          time.Time
}

type changeLogEntity struct {
	ID                    string
	ChangeSecurity        string
	ChangeResiliency      string
	ChangeInfrastructure  string
	ChangeQuality         string
	ChangeAppArchitecture string
	ChangeUnitTest        string
	LastUpdated           time.Time
}

func convertDtoToEntity(serviceID string, serviceName string, serviceDescription string, mrReport maturityReport) serviceReportEntity {
	serviceReport := serviceReportEntity{ID: serviceID, Name: serviceName, Description: serviceDescription,
		UnitTestScore: mrReport.LatestCheckResults[0].Message, overallLevel: mrReport.OverallLevel.Name}

	for _, categoryBreakdown1 := range mrReport.CategoryBreakdown {
		if categoryBreakdown1.Category.Name == "Security" {
			serviceReport.SecurityLevel = categoryBreakdown1.Level.Name
		}
		if categoryBreakdown1.Category.Name == "Resiliency" {
			serviceReport.ResiliencyLevel = categoryBreakdown1.Level.Name
		}
		if categoryBreakdown1.Category.Name == "Infrastructure" {
			serviceReport.InfrastructureLevel = categoryBreakdown1.Level.Name
		}
		if categoryBreakdown1.Category.Name == "Quality" {
			serviceReport.QualityLevel = categoryBreakdown1.Level.Name
		}
		if categoryBreakdown1.Category.Name == "Application Architecture" {
			serviceReport.AppArchitectureLevel = categoryBreakdown1.Level.Name
		}
	}
	serviceReport.LastUpdated = time.Now()

	return serviceReport
}

func setupDBConfig() (mongo.Client, context.Context) {

	DB_URL := "mongodb://localhost/?retryWrites=true&w=majority"

	client, err := mongo.NewClient(options.Client().ApplyURI(DB_URL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return *client, ctx
}

func main() {
	fmt.Println("welcome to go")

	dbClient, ctx := setupDBConfig()
	jsonData := map[string]string{
		"query": `
		{
			account {
			services (ownerAlias: "knp_-_content_ingestion") {
				nodes {
				name
				id
				description  
				}
			}
			}
		}
	`}

	token := "X6YVXLwIRjceklWoQOrXD5RHa2VmDrJKMnnf"

	//call get Services endpoint
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://api.opslevel.com/graphql", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	client := &http.Client{Timeout: time.Second * 10}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Accept", "application/json")

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close() // makes sure that response body is closed.
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(data))

	var servicesResponse ServicesResponse
	//var test interface{}
	json.Unmarshal(data, &servicesResponse)
	fmt.Println("response received")

	fmt.Println(servicesResponse)
	// fetch maturity report for each service.

	for _, node := range servicesResponse.Data.Account.Services.Nodes {
		if node.ID == "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNjIw" {

			jsonMR := map[string]string{
				"query": `
				{
					account {
						  service(id: "` + node.ID + `") {
							name
							aliases
							description
							maturityReport {
								categoryBreakdown {
									category {
									name
									id
									description
									}
									level {
									name
									id
									description
									alias
									}
								}
								latestCheckResults(ids : "Z2lkOi8vb3BzbGV2ZWwvQ2hlY2tzOjpQYXlsb2FkLzE5MDY") {
									message
									status
								}
								overallLevel {
									name
									id
									description
									alias
								}
							}
						}
					}
				}
			`}
			bytesMRRequest, _ := json.Marshal(jsonMR)
			//fmt.Println("json Mr is", jsonMR)

			requestMaturityReport, errMR := http.NewRequest("POST", "https://api.opslevel.com/graphql", bytes.NewBuffer(bytesMRRequest))
			if errMR != nil {
				panic(errMR)
			}
			client := &http.Client{Timeout: time.Second * 10}
			requestMaturityReport.Header.Add("Content-Type", "application/json")
			requestMaturityReport.Header.Add("Authorization", "Bearer "+token)
			requestMaturityReport.Header.Add("Accept", "application/json")

			response, err := client.Do(requestMaturityReport)
			if err != nil {
				panic(err)
			}

			defer response.Body.Close() // makes sure that response body is closed.
			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				panic(err)
			}
			//fmt.Println(string(data))

			var mrResponse ServicesResponse
			//var test interface{}
			json.Unmarshal(data, &mrResponse)
			fmt.Println("response received for fetching Maturity Report")
			fmt.Println(mrResponse)

			//save information to database.
			serviceReport := convertDtoToEntity(node.ID, node.Name, node.Description, mrResponse.Data.Account.Service.MaturityReport)
			collection := dbClient.Database("opslevel").Collection("services_report")
			res, err := collection.InsertOne(context.Background(), serviceReport)

			if err != nil {
				panic(err)
			}
			defer dbClient.Disconnect(ctx)

			fmt.Println(res)

		}
	}
}

/* Steps
1. Fetch list of services. - done
2. For each service, fetch maturity report. - done
3. Store maturity report results in table- services_report, and update change_log table with any deltas.
4. Update confluence page with the services report as well as change, with following fields:
	* Service name
	* Levels for 5 different areas- Security, Resiliency, Infrasturcture, Quality, Application Architecture.
	* Unit test coverage - possible?
	* changes since last run.

*/
