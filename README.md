# message-sender

- Create two microservices for real time messages using websockets, and at least one demo subscriber.
### Functionality
1. The first service should listen for incoming messages through the websocket protocol and when a new one arrives, the message should be published into message queue
2. The second service should listen for incoming messages through the message queue and when a new message arrives, the message should be published to all the subscribers through the websocket protocol
### Other Requirements
- At least one of the services should have tests
- Make sure your code is well-structured and maintainable (including tests)
- You can use frameworks and technologies by your choice, but the language for the microservices should be Javascript (node) or Go (or both).
- The source code should be hosted online using github (or similar service)

### TODO:

- [X] Client who send the messages
- [X] Microservice 1
    - [X] Listen messages from websocket
    - [X] Send messages to queue
- [ ] Microservice 2
    - [X] Read messages from queue
    - [X] Send messages to subscriptors
- [ ] Subscribers
    - [ ] Subscribes to messages
    - [ ] Listen messages
    - [ ] Print messages
- [ ] Queue
    - [X] Basic connect, read and send
    - [ ] Consumers and producer structure
    - [ ] Configure security and more parameters
- [ ] Organize project structure following clean architecture
- [ ] Cmd
- [ ] Graph to initiate the app
- [ ] Tests for all the app
- [ ] Improve logs
- [ ] Makefile
- [ ] Improve naming


### Starta rabbitmq image

docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management