# orbital_HeBao
This is the orbital23 project for team HeBao.

## Introduction
The API gateway is mainly built upon two Golang frameworks, namely Kitex (an RPC framework) and Hertz (an HTTP framework). The most basic function of the gateway is to translate the HTTP request from client to Thrift binary protocol, conduct a remote procedure call (RPC) to the backend RPC servers, and finally translate the response back to HTTP. 

Various extensions are integrated into the API gateway to handle extra logic besides the basic function. We use nacos (an extension for Kitex) for service registration and discovery. Thrift Interface Definition Language (IDL) is also used to define the interface contracts between the client, the gateway, and backend microservices. 

## System Design
![system design](https://github.com/Linda-ui/orbital_HeBao/assets/83194176/9236f888-f401-4567-a58d-53cb6219ca62)
The API gateway is implemented as a Hertz server containing multiple Kitex clients, each corresponding to a microservice (implemented as a Kitex server). The Hertz
server receives incoming HTTP requests from the client and invokes the respective Kitex client. Then, the Kitex client will translate the HTTP request to Thrift binary request and initiate a generic RPC call to its backend server. The backend Kitex server handles the request and sends a response back to the Kitex client, which is then translated back to HTTP response and directed back to the client. 

## Components
### Parameter Binding and Validation

This component will process the HTTP request provided by the client for correct identification, extraction and validation of information from the request to be used in the application.

This is generated using the Hertz framework.

### HTTP Mapping Generic Call
This component converts the HTTP request into a generic request based on the interface mapping specifications in the Thrift IDL files. The generic call is then performed to obtain a generic response which is converted back into HTTP response.

This is implemented using Kitex’s generic call feature.

### Service Registry and Discovery
This component refers to a service registry that is responsible for registering and identifying services the API gateway interacts with. The gateway dynamically discovers the service location through the registry, establishes connection and communicates with the service needed. New services can be added, removed, or updated without requiring manual reconfiguration of the API gateway.

This is implemented using an extension for Kitex, nacos.

### Load Balancing
This component distributes incoming API requests across multiple instances of backend services to ensure even distribution of workload and optimal resource utilization. This component helps achieve scalability by adding or removing backend service instances dynamically based on the incoming traffic.

This is implemented using Kitex’s default load-balancer, which uses a weighted round-robin strategy based on weights.

### Rate Limiting
This component prevents the server from being overloaded by sudden traffic increase from a client by controlling the rate of incoming requests for each rate-limited endpoint.

This is implemented using Kitex’s default rate limiter, which specifies the maximum number of concurrent connections and the maximum number of queries per second to a specific endpoint.

### Interactions between Components
1. The Hertz Server component is responsible for receiving incoming requests from clients. It performs parameter binding and validation, extracting data from the request and mapping it to the corresponding API parameters. The Hertz Server then passes the processed request to the generic client.
2. The generic client analyzes the request metadata, such as the HTTP method, headers, URL path, and query parameters, to map the HTTP request to a generic call to the corresponding backend service.
3. The incoming request first passes through the rate limiter component in the API gateway. The rate limiter evaluates the request against predefined limits, such as the number of requests per second (QPS) or concurrent connections. If the request exceeds the configured limits, it may be delayed or rejected based on the defined behavior.
4. If the request is accepted, nacos will help us discover the available instances of the backend service required to process the request from all services registered with it.
5. Then the load balancer uses a predefined algorithm, such as weighted round-robin, to select an instance from the pool of available instances to evenly distribute the incoming requests across the backend service instances to achieve optimal resource utilization and performance.
6. The request is forwarded to the selected instance for further processing. 

## Future Development
We propose the following improvements to be made in the future:
- Enhance the security features:
Implement features such as OAuth 2.0 to have more robust authentication and authorization mechanisms to protect APIs and sensitive data.
- Implement dynamic configuration: 
Allow for on-the-fly changes to routing rules, rate limits, security policies, and other parameters. This enables more flexible and agile management of the gateway configuration without requiring a complete restart or redeployment.
- Prepare for deployment:
Make our API gateway available for deployment by users. Several steps include: Package our API gateway software and related configuration files into a deployable artifact. Provide comprehensive documentation that guides clients through the deployment process. Create automation scripts or configuration templates that streamline the deployment process, with tools such as Docker Compose.

## Conclusion
Our API gateway will serve as a powerful tool to process all incoming requests through the API management system with various functions like rate-limiting and load-balancing, decouple the client interface from backend implementation, and manage microservices in one place. It is especially helpful for hosting large-scale APIs and greatly improves flexibility and performance. 
