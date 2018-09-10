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

## documentation
For API documentation could be found [here](https://documenter.getpostman.com/view/4946571/RWaGVVmS)