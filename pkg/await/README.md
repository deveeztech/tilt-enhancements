# Await process for debbuger

> [!NOTE]  
> Currently only OS based on unix are supported

## How to use

To use this functionality in your code you only have to import the await package:

```go

package main

import(
    _ ""github.com/deveeztech/tilt-enhancements/pkg/await""
)


func main() {

    // your business logic here
}

```


## How it works 

This package is designed to monitor TCP connections for a new connection on a specified port. It operates as follows:

- Infinite Loop: Enters an infinite loop to continuously check for new TCP connections.
Fetch TCP Connections: Attempts to fetch current TCP connections using netstat.Tcp(). If an error occurs, it logs the error and returns it, exiting the function.
- Check Connections: Iterates over the fetched TCP connections. If a connection is found with the specified port and is in the ESTABLISHED state, it logs the detection of a new connection and exits the function successfully.
- Retry Logic: If no matching connection is found, it logs that no connections were detected and retries after a specified polling interval, using time.Sleep() to wait.

To avoid use dependencies of the system, this project implements a similar `netstat` logic