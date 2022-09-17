# flow-api

# Description

# Performnce Limitations/Testing

# Next Steps 
1. Higher performance HTTP server for more simultaneous requests and something that is more RESTful and modular 
2. use some kind of 3rd party queue like kafka in between whatever is sampling network flows and our internal memory store to buffer and handle flow requests. 
3. some form of authentication btw flow client and server 
4. add middleware for logging and metrics instead of custom metrics - less familiar with this and didn't have time to implement it 
5. Add limit to size of incoming payloads
6. rate limit server - token bucket rate limiter
7. create a distributed service with data backend? 
8. add normal tests and load tests 

