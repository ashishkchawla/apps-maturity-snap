# Automate AMM capture


## Steps for reading template from Confluence
* Create a template

* Create API token - https://id.atlassian.com/manage-profile/security/api-tokens



## opslevel graphql- https://app.opslevel.com/graphiql

* Create a token for reading amm metrices - https://app.opslevel.com/api_tokens


## How the utility works

### Logic:
1. Fetch list of services.
1. For each service, fetch maturity report. 
1. Store maturity report results in table- services_report, and update change_log table with any deltas.
1. Update confluence page with the services report as well as change, with following fields:
	* Service name
	* Levels for 5 different areas- Security, Resiliency, Infrasturcture, Quality, Application Architecture.
	* Unit test coverage and integration test coverage - possible?
	* changes since last run.

  

## Executing the utility
### Optional - First time setup only
1. Create API tokens in Opslevel and Confluent.
1. Install Go lang- https://go.dev/doc/install
1. go mod init github.com/ashishkchawla/automate-confluence
1. Setup Mongodb
    * docker-compose up -d
    * docker ps -a
    * docker exec -it 436fe4428393 bash

1. Create database and table 
    * docker exec -it db bash
    * mongosh
    * use opslevel // create database
    * db.createCollection("services_report") // create collection
    * db.createCollection("change_log")
    * // insert documents

### Run the utility
1. run go program
```
    go run main.go <opslevel_owner_alias> <parentPage_id> <title> <sprint_name> <team_name> <opslevel_token> <confluence_token>
```


## References
* https://github.com/mongodb/mongo-go-driver
* https://www.mongodb.com/docs/drivers/go/current/quick-start/#std-label-golang-quickstart
* Get list of templates: https://chegg.atlassian.net/wiki/rest/api/template/page
* Use basic auth with token as password 
* Fetch list of tempates - https://chegg.atlassian.net/wiki/rest/api/template/page?spaceKey=EPE=
* Fetch template detail- https://chegg.atlassian.net/wiki/rest/api/template/2961262126
```
{
    "templateId": "2725740718",
    "name": "Responsibilities ",
    "description": "",
    "space": {
        "id": 2725511356,
        "key": "CIE",
        "name": "Content Ingestion Engineering",
        "type": "global",
        "status": "current",
        "_expandable": {
            "settings": "/rest/api/space/CIE/settings",
            "metadata": "",
            "operations": "",
            "lookAndFeel": "/rest/api/settings/lookandfeel?spaceKey=CIE",
            "identifiers": "",
            "permissions": "",
            "icon": "",
            "description": "",
            "theme": "/rest/api/space/CIE/theme",
            "history": "",
            "homepage": "/rest/api/content/2725511837"
        },
        "_links": {
            "webui": "/spaces/CIE",
            "self": "https://chegg.atlassian.net/wiki/rest/api/space/CIE"
        }
    },
    "labels": [],
    "templateType": "page",
    "editorVersion": "v2",
    "body": {
        "storage": {
            "value":"<at:declarations><at:string at:name=\"team_name\" /><at:list at:name=\"team_members\"><at:option at:value=\"\" /></at:list><at:string at:name=\"lead_name\" /><at:string at:name=\"em_name\" /><at:string at:name=\"pm_name\" /><at:string at:name=\"date\" /><at:string at:name=\"retro_member\" /><at:string at:name=\"chaos_member\" /><at:string at:name=\"sre_champion_name\" /><at:string at:name=\"fma_member\" /><at:list at:name=\"code_review_members\"><at:option at:value=\"\" /></at:list><at:string at:name=\"tdi_member\" /></at:declarations><ac:structured-macro ac:name=\"info\" ac:schema-version=\"1\" ac:macro-id=\"bb3c4149-8cc9-4d4b-9402-87f8b6c7918e\"><ac:rich-text-body><p>This page highlights the current assignment of shared responsibilities for <at:var at:name=\"team_name\" /> team</p></ac:rich-text-body></ac:structured-macro><p /><p>Team members</p><table data-layout=\"default\" ac:local-id=\"aa683385-f043-47ba-8f95-3c694f413b2a\"><colgroup><col style=\"width: 263.0px;\" /><col style=\"width: 587.0px;\" /></colgroup><tbody><tr><td data-highlight-colour=\"#fffae6\"><p><strong>Team</strong></p></td><td><p><at:var at:name=\"team_name\" /> </p></td></tr><tr><td data-highlight-colour=\"#fffae6\"><p><strong>Team members</strong></p></td><td><p><at:var at:name=\"team_members\" /></p></td></tr><tr><td data-highlight-colour=\"#fffae6\"><p><strong>Lead</strong></p></td><td><p><at:var at:name=\"lead_name\" /></p></td></tr><tr><td data-highlight-colour=\"#fffae6\"><p><strong>EM</strong></p></td><td><p><at:var at:name=\"em_name\" /></p></td></tr><tr><td data-highlight-colour=\"#fffae6\"><p><strong>PM</strong></p></td><td><p><at:var at:name=\"pm_name\" /></p></td></tr><tr><td data-highlight-colour=\"#fffae6\"><p><strong>Date</strong></p></td><td><p><at:var at:name=\"date\" /></p></td></tr></tbody></table><p>&nbsp;</p><h2><ac:emoticon ac:name=\"blue-star\" ac:emoji-shortname=\":blue_book:\" ac:emoji-id=\"1f4d8\" ac:emoji-fallback=\"\\uD83D\\uDCD8\" /></h2><table data-layout=\"default\" ac:local-id=\"20087def-a37c-4fa7-9956-36a86b01f87a\"><tbody><tr><th><p><strong>Rsponsibility</strong></p></th><th><p><strong>Member</strong></p></th></tr><tr><td><p>AMM report, burndown Chart, retro meeting</p></td><td><p><at:var at:name=\"retro_member\" /> (Everyone)</p></td></tr><tr><td><p>Chaos Engineering</p></td><td><p><at:var at:name=\"chaos_member\" /> (All Devs and QE)</p></td></tr><tr><td><p>SRE Champion</p></td><td><p><at:var at:name=\"sre_champion_name\" />(All Devs)</p></td></tr><tr><td><p>FMA discussions, documentation upto date, demos</p></td><td><p><at:var at:name=\"fma_member\" /></p></td></tr><tr><td><p>Code review champions</p></td><td><p><at:var at:name=\"code_review_members\" /></p></td></tr><tr><td><p>TDI dashboard</p></td><td><p><at:var at:name=\"tdi_member\" /></p></td></tr></tbody></table><p /><p />",
            "representation": "storage",
            "embeddedContent": []
        }
    },
    "_links": {
        "self": "https://chegg.atlassian.net/wiki/rest/api/template/2725740718",
        "base": "https://chegg.atlassian.net/wiki",
        "context": "/wiki"
    }
}
```
* get list of services for a given team:
```
{
  account {
    services (ownerAlias: "knp_-_content_ingestion") {
      nodes {
        name
        aliases
        id
        description
        owner {
          name
          manager {
            email
          }
        }
      }
    }
  }
}
```

* query to get service by id
```
{
  account {
    service(id: "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNjIw") {
      name
      aliases
      description
    }
  }
}
```
* query maturity report for a service

```
{
  account {
    service(id: "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNjIw") {
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

// output:
{
  "data": {
    "account": {
      "service": {
        "name": "media-service",
        "aliases": [
          "media-service",
          "Media Service"
        ],
        "description": "Centralised platform for handling media upload, storage &amp; delivery. It offers a standard solution to manage Media assets for both Web &amp; Mobile Platforms",
        "maturityReport": {
          "categoryBreakdown": [
            {
              "category": {
                "name": "Security",
                "id": "Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvMzIz",
                "description": null
              },
              "level": {
                "name": "Level 5",
                "id": "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvNjE",
                "description": "(Basic) Requires CTO Approval",
                "alias": "level_5"
              }
            },
            {
              "category": {
                "name": "Resiliency",
                "id": "Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvMzI0",
                "description": null
              },
              "level": {
                "name": "Level 5",
                "id": "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvNjE",
                "description": "(Basic) Requires CTO Approval",
                "alias": "level_5"
              }
            },
            {
              "category": {
                "name": "Infrastructure",
                "id": "Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvMzI4",
                "description": null
              },
              "level": {
                "name": "Level 3",
                "id": "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvMTYx",
                "description": "(Better) Requires Director Approval",
                "alias": "level_3"
              }
            },
            {
              "category": {
                "name": "Quality",
                "id": "Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvMzI5",
                "description": null
              },
              "level": {
                "name": "Level 5",
                "id": "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvNjE",
                "description": "(Basic) Requires CTO Approval",
                "alias": "level_5"
              }
            },
            {
              "category": {
                "name": "Application Architecture",
                "id": "Z2lkOi8vb3BzbGV2ZWwvQ2F0ZWdvcnkvNjM4",
                "description": null
              },
              "level": {
                "name": "Level 5",
                "id": "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvNjE",
                "description": "(Basic) Requires CTO Approval",
                "alias": "level_5"
              }
            }
          ],
          "overallLevel": {
            "name": "Level 5",
            "id": "Z2lkOi8vb3BzbGV2ZWwvTGV2ZWwvNjE",
            "description": "(Basic) Requires CTO Approval",
            "alias": "level_5"
          }
        }
      }
    }
  }
}
```
* integration test coverage- Z2lkOi8vb3BzbGV2ZWwvQ2hlY2tzOjpQYXlsb2FkLzE5MDk, Z2lkOi8vb3BzbGV2ZWwvQ2hlY2tzOjpQYXlsb2FkLzE5OTA



## Reference Commands

* docker system prune- to remove stopped containers, in case of issues.

* Mongo
  * show collections
  * db.services_report.find({}) // find all in colleection- services_report
  * To empty a collection
    * db.services_report.remove({})
  * To update a document
    * db.services_report.updateOne ({_id:'Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNjIw'}, { $set: {infrastructurelevel: 'Level 5'}})

* Connect to mongodb from go
    * go get go.mongodb.org/mongo-driver/mongo 
    * go get github.com/joho/godotenv // to read .env file.