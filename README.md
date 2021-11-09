# Steps to Start Survivor Server

- start postgres docker container locally:
  - docker run --name survivor-db -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres
- set environment variables:
  - `export APP_DB_HOST=localhost`
  - `export APP_DB_PORT=5432`
  - `export APP_DB_USERNAME=postgres`
  - `export APP_DB_PASSWORD=mysecretpassword`
  - `export APP_DB_NAME=postgres`
- pull the code locally and run:
  - `git clone https://github.com/cbhakar/zombies`
  - `cd zombies`
  - `go mod vendor`
  - `go run main.go`
  - server should start on port 8080


# APIs and their responses

### GET
- description : get details of single survivor with **survivor_id** 
- path `/survivor/{id}`
  - request body : nil
  - eg : localhost:8080/survivor/1
  - success resp : ```{
    "id": 4,
    "name": "billy",
    "age": 30,
    "gender": "male",
    "last_known_location": {
        "latitude": 10.1123,
        "longitude": 11.1111
    },
    "resources": [
        "gun",
        "medicine"
    ],
    "is_infected": false
    }```
  - error response : ```{
    "error": "survivor not found"
    }```

### POST
- description : add details of single survivor 
- path `/survivor`
  - eg : localhost:8080/survivor
  - request body : ```{
    "name":"billy",
    "age" : 30,
    "gender" : "male",
    "latitude" : 10.1123,
    "longitude" : 11.1111,
    "resources" : ["gun", "medicine"]
    }```
  - success resp : ```{
    "survivor_id": 5
    }```
  - error response : ```{
    "error": "Invalid request payload"
    }```

### PUT
- description : update location details of single survivor with **survivor_id** 
- path `/survivor/{id}`
  - request body : ```{
    "latitude" : 11.1111,
    "longitude" : 11.1111
    }```
  - eg : localhost:8080/survivor/1
  - success resp : ```{
    "result": "success"
    }```
  - error response : ```{
    "error": "Invalid request payload"
    }```
    
### DELETE
- description : delete details of single survivor with **survivor_id** 
- path `/survivor/{id}`
  - request body : nil
  - eg : localhost:8080/survivor/11
  - success resp : ```{
    "result": "success"
    }```
  - error response : ```{
    "error": "survivor not found"
    }```
    
### GET
- description : get details of all survivor
- path `/survivors`
  - request body : nil
  - eg : localhost:8080/survivors
  - success resp : ```[
    {
        "id": 1,
        "name": "bob",
        "age": 30,
        "gender": "male",
        "last_known_location": {
            "latitude": 10.1123,
            "longitude": 11.1111
        },
        "resources": [
            "tea",
            "coffee"
        ],
        "is_infected": false
    },
    {
        "id": 2,
        "name": "tom",
        "age": 30,
        "gender": "male",
        "last_known_location": {
            "latitude": 10.1123,
            "longitude": 11.1111
        },
        "resources": [
            "tea"
        ],
        "is_infected": true
      }
    ]```
  - error response : ```{
    "error": "survivors not found"
    }```
    
### POST
- description : report a single survivor with **survivor_id**
- path `/report/survivor/{id}`
  - eg : localhost:8080/report/survivor/1
  - request body : nil
  - success resp : ```{
      "result": "survivor flagged successfully"
    }```
  - error response : ```{
    "error": "survivor not found"
    }```
    
### GET
- description : get analytical report of all survivors, infected and robots. 
- path `/report`
  - request body : nil
  - eg : localhost:8080/report
  - success resp : ```{
    "percentage_of_infected_survivors": 20,
    "percentage_of_non_infected_survivors": 80,
    "infected_survivors": [
        {
            "id": 2,
            "name": "tom",
            "age": 30,
            "gender": "male",
            "last_known_location": {
                "latitude": 10.1123,
                "longitude": 11.1111
            },
            "resources": [
                "tea"
            ],
            "is_infected": true
        }
    ],
    "non_infected_survivors": [
        {
            "id": 1,
            "name": "bob",
            "age": 30,
            "gender": "male",
            "last_known_location": {
                "latitude": 10.1123,
                "longitude": 11.1111
            },
            "resources": [
                "tea",
                "coffee"
            ],
            "is_infected": false
        },
        {
            "id": 3,
            "name": "billy",
            "age": 30,
            "gender": "male",
            "last_known_location": {
                "latitude": 10.1123,
                "longitude": 11.1111
            },
            "resources": [
                "gun",
                "medicine"
            ],
            "is_infected": false
        },
        {
            "id": 5,
            "name": "billy",
            "age": 30,
            "gender": "male",
            "last_known_location": {
                "latitude": 10.1123,
                "longitude": 11.1111
            },
            "resources": [
                "gun",
                "medicine"
            ],
            "is_infected": false
        },
        {
            "id": 4,
            "name": "billy",
            "age": 30,
            "gender": "male",
            "last_known_location": {
                "latitude": 11.1111,
                "longitude": 11.1111
            },
            "resources": [
                "gun",
                "medicine"
            ],
            "is_infected": false
        }
    ],
    "robots": [
        {
            "model": "1G1EJ",
            "serialNumber": "3YAPLZWS1CA72Z7",
            "manufacturedDate": "2021-12-13T07:23:55.6296893Z",
            "category": "Land"
        },
        {
            "model": "W03V4",
            "serialNumber": "A9X1YKONXUNGCMU",
            "manufacturedDate": "2021-11-22T07:23:55.6296699Z",
            "category": "Land"
        },
        {
            "model": "ZOLZ2",
            "serialNumber": "0NI52PYYJ0EMYWT",
            "manufacturedDate": "2021-11-18T07:23:55.6296662Z",
            "category": "Flying"
        }
    ]
    }```
