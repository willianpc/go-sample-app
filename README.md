# Go Sample App

A simple Go application which contains the following features:

1. Acts as an HTTP client
1. Acts as an HTTP server
1. Makes use of Redis for caching

## What it does

1. The application exposes a web server that receives a querystring `q` as its only argument.
1. The string provided for `q` will be used to query Google.
1. The results from the first page are then filtered and all titles of results are returned in a JSON format.
1. The result is then cached in Redis for 15 seconds.
1. If the same query is made within 15 seconds, the data is retrieved from Redis instead of querying Google again.

## How to start the application

### Start Redis

    $ docker compose up -d

### Start the application

    $ go run .

## Running the application

1. Call the URL http://localhost:9090?q=myquery

