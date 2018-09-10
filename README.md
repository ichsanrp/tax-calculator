# Tax Calculator
A tax calculator project

---------------------------------------
 * [Requirements](#requirements)
 * [Installation](#installation)
 * [Development](#development)
 * [Documentation](#documentation)

 ## Requirements
  * Docker 18.06.1-ce
  * Go 1.8 or higher
  * MySQL (4.1+)

## Installation
Simply run 
```bash
$ sudo docker-compose up
```

## development
```
Project Structure
 /public
    index.html
 /src                     
    /taxCalculator
        /controller.go      #Controller is collection of http router handler that handler incoming http request
        /model.go           #Model is collection of business logic and algorithm
        /statement.go       #Statement is SQL Query collection that will be used to manipulate data in db
        /type.go            #type is collection of Object that will be used to serialize and deserialize
 Dockerfile                 #Multi-stage build dockerfile to optimize image for docker use
 docker-compose.yml         #docker orchestration script that will create deployment env setting 
```

to develop new feature please create http handler in controller and create new business logic in model. for handling db operation please prepare any db query ini statement.go

### Database Design
![ERD Diagram](https://github.com/ichsanrp/tax-calculator/blob/master/ERD.png "Tax Calculator ERD")

Tax calculator had two table item and session. Session table will store session that created for each calculation, one session will be have multiple item since item will be stored based on session. In every New Calculation we create new session.

Item table is used to store item information for each session. tax will be calculated in server and stored here so any changes in tax calculation will be done in business logic model and old calculation will remain same.

## documentation
For API documentation could be found [here](https://documenter.getpostman.com/view/4946571/RWaGVVmS)