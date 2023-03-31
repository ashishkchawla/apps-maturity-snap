package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

/** Confluence Models **/
type template struct {
	Name        string      `json:"name"`
	ContentBody contentBody `json:"body"`
}

type contentBody struct {
	Storage storage `json:"storage"`
}

type storage struct {
	Content         string   `json:"value"`
	Representation  string   `json:"representation"`
	EmbeddedContent []string `json:"embeddedContent"`
}

type createPageRequest struct {
	Type        string      `json:"type"`
	Title       string      `json:"title"`
	Space       space       `json:"space"`
	ContentBody contentBody `json:"body"`
	Ancestors   []ancestor  `json:"ancestors"`
}

type space struct {
	Key string `json:"key"`
}

type ancestor struct {
	ID string `json:"id"`
}

//<parentPage_id> <title> <sprint_name> <team_name>
/**Confluence model ends **/

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

	fmt.Println("infrastructure level old***" + oldReportEntity.InfrastructureLevel)
	fmt.Println("infrastructure level new***" + newReportEntity.InfrastructureLevel)

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

func updatePageCreateRequestwithAMM(content string, serviceReport serviceReportEntity, changeLog changeLogEntity) string {

	regexTableRow, err := regexp.Compile("<tr><td><p><at:var at:name=\"service_name\" />.*<at:var at:name=\"change_integration_test_coverage\" /> </p></td></tr>")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(regexTableRow.MatchString(content))
	//fmt.Println(regexTableRow.FindString(content))
	var row string
	if regexTableRow.MatchString(content) {
		row = regexTableRow.FindString(content)
	}
	//fmt.Println(row)

	tempRow := row

	// Build change description
	var changeAMM string
	if changeLog.ChangeAppArchitecture != "" {
		changeAMM = "<br/>Architecture : " + changeLog.ChangeAppArchitecture
	}
	if changeLog.ChangeInfrastructure != "" {
		changeAMM = "<br/>Infrastructure : " + changeLog.ChangeInfrastructure
	}
	if changeLog.ChangeQuality != "" {
		changeAMM = "<br/>Quality : " + changeLog.ChangeQuality
	}
	if changeLog.ChangeResiliency != "" {
		changeAMM = "<br/>Resiliency : " + changeLog.ChangeResiliency
	}
	if changeLog.ChangeSecurity != "" {
		changeAMM = "<br/>Security : " + changeLog.ChangeSecurity
	}

	tempRow = strings.Replace(tempRow, "<at:var at:name=\"service_name\" />", serviceReport.Name, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"security_level\" />", serviceReport.SecurityLevel, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"resiliency_level\" />", serviceReport.ResiliencyLevel, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"infra_level\" />", serviceReport.InfrastructureLevel, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"quality_level\" />", serviceReport.QualityLevel, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"arch_level\" />", serviceReport.AppArchitectureLevel, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"overall_level\" />", serviceReport.OverallLevel, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"change_in_AMM\" />", changeAMM, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"unit_test_coverage\" />", serviceReport.UnitTestScore, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"change_unit_test_coverage\" />", changeLog.ChangeUnitTest, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"integration_test_coverage\" />", serviceReport.IntegrationTestScore, 1)
	tempRow = strings.Replace(tempRow, "<at:var at:name=\"change_integration_test_coverage\" />", changeLog.ChangeIntegrationTest, 1)

	tempRow = tempRow + row

	//fmt.Println("tempRow**" + tempRow)

	content = regexTableRow.ReplaceAllString(content, tempRow)

	// Have AMM score names instaead of levels
	content = strings.ReplaceAll(content, "<td><p>Level 1</p></td>", "<td class=\"highlight-green\" data-highlight-colour=\"green\"><p>Awesome</p></td>")
	content = strings.ReplaceAll(content, "<td><p>Level 2</p></td>", "<td class=\"highlight-yellow\" data-highlight-colour=\"yellow\"><p>Best</p></td>")
	content = strings.ReplaceAll(content, "<td><p>Level 3</p></td>", "<td class=\"highlight-blue\" data-highlight-colour=\"blue\"><p>Better</p></td>")
	content = strings.ReplaceAll(content, "<td><p>Level 4</p></td>", "<td class=\"highlight-red\" data-highlight-colour=\"red\"><p>Good</p></td>")
	content = strings.ReplaceAll(content, "<td><p>Level 5</p></td>", "<td class=\"highlight-grey\" data-highlight-colour=\"grey\"><p>Basic</p></td>")

	// append other variables
	content = strings.ReplaceAll(content, "<at:var at:name=\"team_name\" />", os.Args[5])
	content = strings.ReplaceAll(content, "<at:var at:name=\"sprint_name\" />", os.Args[4])

	return content
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

	if len(os.Args) <= 5 {
		fmt.Println("Please pass following arguments: <opslevel_owner_alias> <parentPage_id> <title> <sprint_name> <team_name>")
		os.Exit(1)
		// do something with command
	}

	dbClient, ctx := setupDBConfig()

	serviceRequestBody := map[string]string{
		"query": `
		{
			account {
			services (ownerAlias: "` + os.Args[1] + `") {
				nodes {
				name
				id
				description  
				}
			}
			}
		}
	`}

	opsLevelToken := "X6YVXLwIRjceklWoQOrXD5RHa2VmDrJKMnnf"

	confluenceToken := "ATATT3xFfGF0KLpcwfCD_tAc5I8bodmC6AQaDiwgRfCxNguzYRxCfXqX-UjbQpZ9lmwFJ17GGOfWojy6r5c_GZedaiGgZkfvpllV-oj1Ypxq_tfyA3G39GEQpz6LugfePxhn9EOzTwH0WMRcMSxwJbUQi8KhMAJPIqv0FdVBTVWwo1Rp1DA3IBw=909708D4"

	//call get Services endpoint
	jsonValue, _ := json.Marshal(serviceRequestBody)
	request, err := http.NewRequest("POST", "https://api.opslevel.com/graphql", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	client := &http.Client{Timeout: time.Second * 10}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+opsLevelToken)
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
	//fmt.Println("response received")

	/** COnfluence template request begins*/

	requestTemplate, err := http.NewRequest("GET", "https://chegg.atlassian.net/wiki/rest/api/template/2961262126", nil)
	if err != nil {
		panic(err)
	}
	client = &http.Client{Timeout: time.Second * 10}
	requestTemplate.Header.Add("Content-Type", "application/json")
	requestTemplate.SetBasicAuth("aschawla@chegg.com", confluenceToken)
	requestTemplate.Header.Add("Accept", "application/json")

	responseTemplate, err := client.Do(requestTemplate)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close() // makes sure that response body is closed.
	dataTemplate, err := ioutil.ReadAll(responseTemplate.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dataTemplate))

	var templateResponse template
	json.Unmarshal(dataTemplate, &templateResponse)
	fmt.Println("response received for template**")
	//fmt.Println(templateResponse)

	createNewRetroPage := createPageRequest{Type: "page", Ancestors: []ancestor{{ID: os.Args[2]}}, Title: os.Args[3], Space: space{Key: "EPE"}, ContentBody: templateResponse.ContentBody}

	/*Confluence template request ends*/

	fmt.Println(servicesResponse)
	// fetch maturity report for each service.

	for _, node := range servicesResponse.Data.Account.Services.Nodes {
		//if node.ID == "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNjIw" || node.ID == "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNjI1" {

		maturityReportRequestBody := map[string]string{
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
		bytesMRRequest, _ := json.Marshal(maturityReportRequestBody)
		//fmt.Println("json Mr is", jsonMR)

		requestMaturityReport, errMR := http.NewRequest("POST", "https://api.opslevel.com/graphql", bytes.NewBuffer(bytesMRRequest))
		if errMR != nil {
			panic(errMR)
		}
		client := &http.Client{Timeout: time.Second * 10}
		requestMaturityReport.Header.Add("Content-Type", "application/json")
		requestMaturityReport.Header.Add("Authorization", "Bearer "+opsLevelToken)
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
		//fmt.Println(mrResponse)

		// change_log
		// read information from service_report, compare with current response, generate change_log, persist in the two collections
		//save information to database.
		currentServiceReport := convertDtoToEntity(node.ID, node.Name, node.Description, mrResponse.Data.Account.Service.MaturityReport)
		serviceReportCollection := dbClient.Database("opslevel").Collection("services_report")
		changeLogCollection := dbClient.Database("opslevel").Collection("change_log")
		filter := bson.D{{"_id", node.ID}}

		var oldServiceReport serviceReportEntity

		err1 := serviceReportCollection.FindOne(context.Background(), filter).Decode(&oldServiceReport)
		var isChanged bool
		var changeLog changeLogEntity
		if err1 != nil {
			//doesn't exist already
			fmt.Println("it doesn't exist, so insert new row.")
			fmt.Println(err1)

			res, err := serviceReportCollection.InsertOne(context.Background(), currentServiceReport)
			fmt.Println(res)

			if err != nil {
				panic(err)
			}
		} else {
			// exist already.
			fmt.Println("it exists, so update the row.")

			changeLog, isChanged = captureChangeLog(node.ID, currentServiceReport, oldServiceReport)
			fmt.Println("isChanged**")
			fmt.Println(isChanged)
			fmt.Println(changeLog)

			res, err := serviceReportCollection.ReplaceOne(context.Background(), filter, currentServiceReport)
			if err != nil {
				panic(err)
			}
			fmt.Println(res)

			if isChanged {
				// save to database, and update confluence document
				res, err := changeLogCollection.InsertOne(context.Background(), changeLog)
				fmt.Println(res)

				if err != nil {
					panic(err)
				}

			}

		}
		createNewRetroPage.ContentBody.Storage.Content = updatePageCreateRequestwithAMM(createNewRetroPage.ContentBody.Storage.Content, currentServiceReport, changeLog)

		fmt.Println("oldServiceReport is***")
		fmt.Println(oldServiceReport)

		defer dbClient.Disconnect(ctx)

		//}
	}
	/** Confluence create new page begins **/

	/*createPageRequestJson := []byte(`{
					{
						"type": "page",
						"title": "Retro test-4",
						"space": {
							"key": "EPE"
						},
						"body": {
							"storage": {
								"value": "testing",
								"representation": "storage",
	            				"embeddedContent": []
							}
				}`)*/
	jsonCreatePageRequest, _ := json.Marshal(createNewRetroPage)
	fmt.Println("create page request***")
	fmt.Println(createNewRetroPage)

	//fmt.Println("confluenceToken***")
	//fmt.Println(confluenceToken)

	request, err = http.NewRequest("POST", "https://chegg.atlassian.net/wiki/rest/api/content/", bytes.NewBuffer(jsonCreatePageRequest))
	if err != nil {
		panic(err)
	}
	client2 := &http.Client{Timeout: time.Second * 30}
	request.Header.Add("Content-Type", "application/json")
	request.SetBasicAuth("aschawla@chegg.com", confluenceToken)
	request.Header.Add("Accept", "application/json")

	response, err = client2.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close() // makes sure that response body is closed.
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Confluence page created")
	fmt.Println(string(data))

	/** Confluence create new page ends **/
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
