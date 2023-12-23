# System Design

## Overview
This document provides an overview of the system design for the csv_query_server application. This application is a web-based platform that allows users to operate a csv file, include query and modify.

### Architecture
The server follows a monolithic architecture, consisting of the following reasons:

1. There are few APIs, so there is no need to adopt a microservices' architecture.
2. The usage requirements are not demanding, and a monolithic architecture is sufficient to meet the needs.

### Directory Structure
    |__ csv_query_server
        |__ doc
            |__ design.md
        |__ internal
            |__ handler
                |__ modify_handler.go
                |__ query_handler.go
                |__ request_validator.go
                |__ request_validator_test.go
            |__ repository
                |__ csv_repository.go
                |__ csv_repository_test.go
            |__ server
                |__ server.go
        |__ pkg
            |__ csv
                |__ data.csv
        |__ go.mod
        |__ main.go
        |__ README.md

## Details

### WorkFlow

- When I received this task, I initially spent some time deliberating on the choice of programming language. Using C++ could earn extra points, but my proficiency in C++ network development was insufficient, making it highly likely that I would not be able to complete the task. Therefore, after reviewing some open-source projects in C++, I decided to use Go instead of C++.
- After deciding the programming language, the next step was to choose a networking framework. After reviewing performance comparisons and relevant documentation of various networking frameworks, I decided to use the standard library of Go. Since the requirements of this task are relatively simple, the performance differences among different networking frameworks in this scenario are not significant.
- Then, I use Go build a simple service, start 2 servers listening on different ports. One server handle query, the other one handle modify.
- After test the ports, I start to coding the business logic. Handle a data.csv by an API is not easy as I think. 
  - I read the file to memory, and save it to a global variable.
  - Using a RWLock to ensure its safe when in concurrency write.
  - When modify the data, I save the data to a temp file, it will not write to data.csv util finish this operating 
  - Using regex to check the validation of requests
  - Using a lot of string operations to implement data querying and modification.
- Then I test the service, but there is not enough time for me to test all the cases. So I test all simple cases and a part of complex cases.
- Finally, write some doc for the service.

### Problems Encountered and Solutions

1. Bug of "go.uber.org/fx", I firstly use this package to build my service. But I find it will make code more complex and hard to debug, so I abandon it.
2. When I test query API, I find a bug that trim function will delete all the "1" of "[C1==test1]", and get a value of "test". After read the source code of trim, I write a new function and solve this.
3. I think it's a difficult task to judge request is valid or not. I test many regex, and finally get a good result.
4. It's also difficult to operate data by column name, I search many doc and test some ways I checked. Although the code accomplish most of the operates, but there are still many problems.
5. After I make a docker images, I run it with relative path of local directory, that case a local file map to a directory in docker.
6. Docker has no permission to read local file that map to docker. I search some docs and find the reason is docker has no permission to this directory. I solve it by sharing the directory to docker.

### Question about task
1. I find it's uncertain that when and how to create the column names. I directly set the column name as [C1 C2 C3]. And I leave a way to set column, delete data.csv and build a new data.csv which include new column names. Server will read new file and use new column.
2. It's uncertain that which column is unique, or none column is unique. When I write UPDATE, I find it's possible to update a row to make it as same as an exist row if C1 is not unique. But the description of INSERT said we can't insert 2 same rows.

### Further Optimization

- There are many optimizations we can do for this service, due to the limited time, I write them below:
1. The service should use a log framework to save logs, but it only has some debugging output now.
2. Each time we modify the data, server will write a temp file, it will block main goroutine especially when data.csv is big. 
   - We can use a cronjob to save data to data.csv, so that we don't need to write data.csv each time we modify it.
   - Also, we can use a background task to write data and do not block main goroutine. Use a signal or message queue to notify main goroutine that write is finished.
   - Thirdly, we can incrementally save the modify operations like AOF, and operate them all to data.csv when we need update data.csv.
3. Unit test coverage is incomplete. 