# go-loadbalancer

A load balancer supporting multiple load balancing strategies

## Load Balancing Algorithms 
- Round Robin
- Weighted Round Robin

### [Round Robin](https://en.wikipedia.org/wiki/Round-robin_scheduling)
Round Robin is a simple load balancing algorithm that involves 
distributes network traffic to different app servers (replicas) 
in a sequential order with no preference (one after the other).

### [Weighted Round Robin](https://en.wikipedia.org/wiki/Weighted_round_robin)
Weighted Round Robin is a variant of Round Robin that takes into 
consideration the possibility of the different app servers having 
different computing capacity. The capacity of each app server is 
specified using __weights__. A weight is simply the relative ratio 
of the computing capacity of each app server.

## Features
- [X] Easy configuration of services and its replicas using `yaml`
- [X] Health checks for each app replica
- [X] Path forwarding for each service
- [X] Forwarding to replicas through reverse proxies
- [ ] Hot Reloading for config files

## Setup

- Clone repository using 
```shell
git clone https://github.com/Emmrys-Jay/go-loadbalancer.git
```

- Create a config yaml file in the root directory. __Sample:__
```yaml
# Sample config yaml
strategy: "RoundRobin"
services: 
  - name: "Service01"
    matcher: "/api/v1"
    strategy: "RoundRobin"
    replicas: 
    - url: "http://localhost:8081"
    - url: "http://localhost:8082"
```

| Yaml Field | Description                                                                                                  |
|------------|--------------------------------------------------------------------------------------------------------------|
| `strategy` | The default load balancing algorithm for all <br/>services                                                   |
| `services` | An array of services made up of the name <br/> `matcher`, `strategy` and `replicas`                          |
| `matcher`  | A string in the format of a url path that matches a network request to the <br/> load balancer  to a service |
| `strategy` | The load balancing algorithm for a particular service.                                                       |
| `replicas` | An array which contains the address of each instance of the server in <br> the form `<host>:<port>`.         ||
| `url`      | The url of each server replica in the form `<host>:<port>`.                                                  |
| `weight`   | Used in `WeightedRoundRobin` strategy to signify the weights of each <br> server replica.                    |

- Run the project using the command below. The default port is `8080`.
```shell
go run cmd/server/main.go --port /port/to/start/the/lb/server --config /path/to/your/config/file
```

- To test the program without going through the process of creating your own config
files, run the program using the default config file for a `RoundRobin` strategy in 
`example/config.yaml` using:
```shell
go run cmd/server/main.go --port /port/to/start/the/lb/server
```

- Forward requests to the load balancer using curl or via the browser. Ensure the 
path in your load balancer url indicates the matcher for the service you want to
visit. For example, sending requests to the load balancer for the service in the 
example `RoundRobin` config file above will be:
```shell
curl localhost:8080/api/v1
```

- You can view a sample of a `WeightedRoundRobin` config file in 
`example/config-weighted.yaml`

## Building a demo 

If you want to try load balancer, you can run the demo server in the `cmd/demo/main.go` directory, which 
starts up a hello world server listening to a specified port.

Launching a demo server is as easy as:
```shell
go run cmd/demo/main.go --port /port/to/start/the/lb/server
```

__NB__: You are free to fork and use this project.