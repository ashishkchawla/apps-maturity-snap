package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	ID                   string `bson:"_id"`
	Name                 string
	Description          string
	SecurityLevel        string
	ResiliencyLevel      string
	InfrastructureLevel  string
	QualityLevel         string
	AppArchitectureLevel string
	UnitTestScore        string
	IntegrationTestScore string
	OverallLevel         string
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
	ChangeIntegrationTest string
	LastUpdated           time.Time
}

func convertDtoToEntity(serviceID string, serviceName string, serviceDescription string, mrReport maturityReport) serviceReportEntity {
	serviceReport := serviceReportEntity{ID: serviceID, Name: serviceName, Description: serviceDescription,
		OverallLevel: mrReport.OverallLevel.Name}

	// Evaluate Unit Test Coverage
	serviceReport.UnitTestScore = "99"
	if mrReport.LatestCheckResults[0].Status == "failed" {
		message := mrReport.LatestCheckResults[0].Message
		regexNumber, _ := regexp.Compile("([0-9.]+)")

		if regexNumber.MatchString(message) {
			serviceReport.UnitTestScore = regexNumber.FindString(message)
		}
	}
	// Evaluate Integration Test Coverage
	serviceReport.IntegrationTestScore = "99"
	if mrReport.LatestCheckResults[0].Status == "failed" {
		message := mrReport.LatestCheckResults[0].Message
		regexNumber, _ := regexp.Compile("([0-9.]+)")

		if regexNumber.MatchString(message) {
			serviceReport.IntegrationTestScore = regexNumber.FindString(message)
		}
	}

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

func captureChangeLog(serviceID string, newReportEntity serviceReportEntity, oldReportEntity serviceReportEntity) (changeLogEntity, bool) {

	var isChange bool = false

	formatChange := func(old string, new string) string {
		if new != old {
			isChange = true
			return (old + "->" + new)
		}
		return ""
	}

	changeLogEntity1 := changeLogEntity{ID: serviceID,
		ChangeSecurity:        formatChange(oldReportEntity.SecurityLevel, newReportEntity.SecurityLevel),
		ChangeResiliency:      formatChange(oldReportEntity.ResiliencyLevel, newReportEntity.ResiliencyLevel),
		ChangeInfrastructure:  formatChange(oldReportEntity.InfrastructureLevel, newReportEntity.InfrastructureLevel),
		ChangeQuality:         formatChange(oldReportEntity.QualityLevel, newReportEntity.QualityLevel),
		ChangeAppArchitecture: formatChange(oldReportEntity.AppArchitectureLevel, newReportEntity.AppArchitectureLevel),
		ChangeUnitTest:        formatChange(oldReportEntity.UnitTestScore, newReportEntity.UnitTestScore),
		ChangeIntegrationTest: formatChange(oldReportEntity.IntegrationTestScore, newReportEntity.IntegrationTestScore),
		LastUpdated:           time.Now()}

	return changeLogEntity1, isChange
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
		//if node.ID == "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNjIw" {

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
								latestCheckResults(ids : ["Z2lkOi8vb3BzbGV2ZWwvQ2hlY2tzOjpQYXlsb2FkLzE5MDk","Z2lkOi8vb3BzbGV2ZWwvQ2hlY2tzOjpQYXlsb2FkLzE5OTA"]) {
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

		// change_log
		// read information from service_report, compare with current response, generate change_log, persist in the two collections
		//save information to database.
		currentServiceReport := convertDtoToEntity(node.ID, node.Name, node.Description, mrResponse.Data.Account.Service.MaturityReport)
		serviceResportCollection := dbClient.Database("opslevel").Collection("services_report")
		changeLogCollection := dbClient.Database("opslevel").Collection("change_log")
		filter := bson.D{{"_id", node.ID}}

		var oldServiceReport serviceReportEntity

		err1 := serviceResportCollection.FindOne(context.Background(), filter).Decode(&oldServiceReport)
		if err1 != nil {
			//doesn't exist already
			fmt.Println("it doesn't exist, so insert new row.")
			fmt.Println(err1)

			res, err := serviceResportCollection.InsertOne(context.Background(), currentServiceReport)
			fmt.Println(res)

			if err != nil {
				panic(err)
			}
		} else {
			// exist already.
			fmt.Println("it exists, so update the row.")

			res, err := serviceResportCollection.ReplaceOne(context.Background(), filter, currentServiceReport)
			if err != nil {
				panic(err)
			}
			fmt.Println(res)

			changeLog, isChanged := captureChangeLog(node.ID, currentServiceReport, oldServiceReport)
			fmt.Println("isChanged**")
			fmt.Println(isChanged)
			fmt.Println(changeLog)

			if isChanged {
				// save to database, and update confluence document
				res, err := changeLogCollection.InsertOne(context.Background(), changeLog)
				fmt.Println(res)

				if err != nil {
					panic(err)
				}

			}

		}

		fmt.Println("oldServiceReport is***")
		fmt.Println(oldServiceReport)

		defer dbClient.Disconnect(ctx)

		//}
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
