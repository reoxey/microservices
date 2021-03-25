# E-commerce Demo - Golang Hexagonal Microservices

### Technologies and patterns:

- Hexagonal Principle
- Docker
- Kubernetes
- Mysql as main db
- Redis as Cache
- Kafka / RabbitMQ as Queue
- Kong API Gateway
- gRPC for sync communication
- JWT Authorization
- GIN for REST API

### Microservices in the demo are:

- **User** - Handle Authentications through JWT
- **Product** - Store products catalog and inventory
- **Shipping** - Store user addresses and order shipping details
- **Cart** - Shopping cart
- **Order** - Stores items from cart, update payment details and shipping status
- **UI** - Web frontend exposed through API gateway which interacts with all microservices.

### Design

![alt text](https://raw.githubusercontent.com/reoxey/microservices/master/img/ms.png)

#### ---

Kubernetes helm scripts are present in the `scripts` directory which are used in the project.
`.http` can be used to test functionality of individual services. `ecom.http` will test the combined functionality of
all services.
