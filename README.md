# flow-api

## Setup and Use

First, download the repo and get all dependencies:

`$ git clone https://github.com/si74/flow-api`

Dependencies are managed by go mod and are vendored and committed in this repo. 

Run the service: 

`$ make test`

Insert flows into the database: 

`$ ./cmd/flowd/test.sh`

(Note: This contains the flow metrics: )

Retrieve aggregated flow data: 

```
$ curl -X GET "localhost:8080/flows?hour=3" | jq .
[]

$ curl -X GET "localhost:8080/flows?hour=1" | jq .
[{"src_app":"foo","dest_app":"bar","vpc_id":"vpc-0","bytes_tx":300,"bytes_rx":900,"hour":1},{"src_app":"baz","dest_app":"qux","vpc_id":"vpc-0","bytes_tx":100,"bytes_rx":500,"hour":1}]
```

Metrics are available as well: 

`$ curl localhost:8080/metrics`

## Testing 

Simply run `$ make test` to run unit tests. 

## Design & Implementation 

flow-api is a flow aggregation surface with two key components: 

1. <b>A flow RESTful API.</b>
  - This is a fairly standard structure which utilized the golang http package - the server and the mux along with several handlers. Note the server handles context cancellations gracefully via the use of error groups. 

  - I also added logging that was configurable and prometheus Go metrics - alongside some custom http server and flow datastore metrics. 

2. <b>A flow datastore and aggregation package</b>

  - The main flow datastore structure utilized a thread-safe mapping of flow tuple identifiers to flowList structures. The flow tuple identifier consists of three values - the src app, the dst app, and the vpc ID. This map is acceptable if there is a limited subset of src, dst, and vpc options; however, if this were to be IP addresses instead of apps, there would be a significantly larger subset of identifiers and a map would not be ideal. 

  - I debated the structure of the flowList quite a bit and elected to create a generic interface that would enable me to easily swap out implementations if I wanted to pursue more efficient structures. 

  - The first implementation - flowlistV1 - was a simple unordered doubly linked list of flow data points. While this made insertion quite fast, it made data aggregation significantly more time-intensive as it required full iteration through the entire list of flows (O(N)). 

## Limitations and Next-Steps 

1. <b>Datastore:</b>
  - The first thing I would look into is further optimizing the flow datastore. 
  - Improved flowList - ordered linked list 
  - Is there lock contention using a map? 
  - Move away from using map entirely and use some kind of in-memory time series database - best way to store multi-dimensional data. 
// TODO(sneha): Figure this out
  - Move from a single service to a series of microservices and a more distributed architecture 
    - use some kind of 3rd party queue like kafka in between whatever is sampling network flows and our internal memory store to buffer and handle flow requests. Either an in-memory buffer queue or external one to handle burstiness.

2. <b>HTTP Server:</b>
  - Limit size of incoming payloads
  - Token-bucket rate limiter 
  - Higher performance HTTP server for more simultaneous requests and something that is more RESTful and modular 
  - Authentication 
  - Add middleware for logging and metrics instead of custom metrics - less familiar with this and didn't have time to implement it 
  - More configurable and have a config package/config.yaml 

3. <b>Other:</b> 
  - Dashboards 
  - Load Testing
  - Dockerized Integration Testing 

## TODO

- tests for server 
- metrics for flowdb
- improve datastore efficiency 
