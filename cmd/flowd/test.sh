#!/bin/bash

# Curl command to insert a number of test flows 

curl -X POST "http://localhost:8080/flows" -vv \
-H 'Content-Type: application/json' \
-d '[{"src_app":"foo","dest_app":"bar","vpc_id":"vpc-0","bytes_tx":100,"bytes_rx":300,"hour":1},{"src_app":"foo","dest_app":"bar","vpc_id":"vpc-0","bytes_tx":200,"bytes_rx":600,"hour":1},{"src_app":"baz","dest_app":"qux","vpc_id":"vpc-0","bytes_tx":100,"bytes_rx":500,"hour":1},{"src_app":"baz","dest_app":"qux","vpc_id":"vpc-0","bytes_tx":100,"bytes_rx":500,"hour":2},{"src_app":"baz","dest_app":"qux","vpc_id":"vpc-1","bytes_tx":100,"bytes_rx":500,"hour":2}]'

