# Stream-go rpc communication between server/client Workflow

1. Start the server in src/server/server.go in main() function - this will start the grpc server and it will start listening locally on localhost: 50005
2. Start the client in src/client/client.go in main() function - in this a client will start communicating with the grpc server.
3. Client will generate the requests (based on value of num_request) and will communicate to and fro. It will greet the server based on time (as defined in greet() function) for the interval of 2 seconds (as defined in interval_time variable)
4. After each time interval, client will send a random string generated from list of random strings (defined in generateString() function)
5. Upon receiving any request, the server will also take the input, greet the client and based on time_interval, it will also generate random string and send it to client
6. The client when getting any response, will check for 'hello' substring inside it. If it exists, then it will print that.
7. Currently I have commented out the part where messages from client and server and logged. You can uncomment them in client.go and server.go - they will be accompanied by a comment asking to uncomment the line below (currently line 73 in client.go and 99 in client.go)

