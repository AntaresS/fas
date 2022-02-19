Flow Aggregation Services (FAS)
-------------------------------

# Overview
FAS is a network flow aggregation service that persists network stats in MongoDB and provides aggregation queries by hours.

# Architecture

## KeyComponents
 - FAS Backend
 - MongoDB

  ```ditaa {cmd=true args=["-E"]}
  +-----------+     +-----------+   
  | Container | ->  | Container |
  |           |     |           |
  |FAS Backend|     | MongoDB   |
  |           |     |           |
  +---+-------+     +-----------+
      : 
      : - POST /flow
      : - GET  /flow?hour=1 
      :                       
      |       
     User
  ```

## Code Structure
```
/fas
../config    # server config and loading logic
../data      # data schema and validation
../../query  # query template
../server    # app server & api layer
../storage   # storage client provider 
```

# Usage

## Pre-req
Docker must be installed first.

## Run
```shell
git clone git@github.com:AntaresS/fas.git

cd fas
docker compose up --build -d
# now the system should have started running

# test with data ingestion
curl -X POST localhost:9527/flow -H 'Content-Type: application/json' -d '[{"src_app": "foo", "dest_app": "bar", "vpc_id": "vpc-0", "bytes_tx": 100, "bytes_rx": 500, "hour": 1}]'
curl -X POST localhost:9527/flow -H 'Content-Type: application/json' -d '[{"src_app": "foo", "dest_app": "bar", "vpc_id": "vpc-0", "bytes_tx": 200, "bytes_rx": 1000, "hour": 1}]'

# test with data aggregation
curl 'localhost:9527/flow?hour=1'

# to check server logs
docker logs fas-fas-1
```

# Considerations
### Server Choice
There are a good number of robust web frameworks in GO. Such as Gin, Martini, Gorilla, Echo, etc. Choosing Echo because of its high performance,
simplicity and production readiness. We can also write up a server from scratch but no need to re-invent the wheels.

### DB Choice
To achieve data persistence and aggreation, there are a few types of DB suitable for this task
- SQL (e.g. PostgreSQL)
- No-SQL (e.g. MongoDB)
- Time-series DB (e.g. InfluxDB, or TimescaleDB which is PostgreSQL based)

Since MongoDB provides a solid aggregation pipeline feature (matching, grouping) out of the box, we choose it to build the system .
Similarly, InfluxDB is also a good candidate as it is designed for storing and querying real time data, and is strong for alerting 
and monitoring purposes.  

### Scalability
The whole system is containerized and decoupled in two images. The fas-backend is a stateless service, so it can be easily scaled
up horizontally via docker swarm or k8s. MongoDB can also run in a [sharded cluster](https://hub.docker.com/r/bitnami/mongodb-sharded/) with replications. 

### Testing and More
This application is a functional MVP. There are a lot more interesting topics to explore and things to enhance.
Here is a potential list of things can be explored in the future iterations
- Test coverage: unit tests, integration-tests, e2e tests, etc.
- Stress test: 
  - stress the system with large number of concurrent requests.
  - aggregation on large data set and see how it performs.
- Performance Benchmarking: compare the performance with other DB choices, mongo vs. InfluxDB vs. TimescaleDB