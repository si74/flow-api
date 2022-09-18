# flow-api

# Setup and Use

# Testing 

Simply run `make test` to run unit tests. 

# Design & Implementation 

flow-api is a flow aggregation surface with two key components: 

1. A flow RESTful API. 

This is a fairly standard structure which utilized the golang http package - the 
server and the mux along with several handlers. Note the server handles context cancellations gracefully via the use of error groups. 

I also added logging that was configurable and prometheus go metrics - alongside some custom http server and flow datastore metrics. 

2. A flow datastore and aggregation package

The main flow datastore structure utilized a thread-safe map of flow tuple identifier to a flowList structure. The flow tuple identifier consists of three values - the src app, the dst app, and the vpc ID. 

I debated the structure of the flowList quite a bit and elected to create a generic interface that would enable me to easily swap out implementations if I wanted to pursue more efficient structures. 

- The first implementation - flowlistV1 - was a simple unordered doubly linked list of flow data points. While this made insertion quite fast, it made data aggregation significantly more time-intensive as it required full iteration through the entire list of flows (O(N)). 

# Limitations and Next-Steps 

1. Datastore: 
- The first thing I would look into is further optimizing the flow datastore. 

2. HTTP Server: 

Further steps: 

- Additional testing: 
- Load testing 

# Next Steps 
1. Higher performance HTTP server for more simultaneous requests and something that is more RESTful and modular 
2. use some kind of 3rd party queue like kafka in between whatever is sampling network flows and our internal memory store to buffer and handle flow requests. 
3. some form of authentication btw flow client and server 
4. add middleware for logging and metrics instead of custom metrics - less familiar with this and didn't have time to implement it 
5. Add limit to size of incoming payloads
6. rate limit server - token bucket rate limiter
7. create a distributed service with data backend? 
8. add normal tests and load tests 

# TODO
- tests for flow
- tests for server 
- integration tests 
- makefile
- metrics for flowdb
- metrics for http server 
- load tests
