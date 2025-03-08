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
- [X] Microservice 2
    - [X] Read messages from queue
    - [X] Send messages to subscriptors
- [X] Subscribers
    - [X] Subscribes to messages
    - [X] Listen messages
    - [X] Print messages
- [ ] Queue
    - [X] Basic connect, read and send
    - [ ] Configure security and more parameters
- [x] Organize project structure following clean architecture
- [x] Cmd
- [X] Graph to initiate the app
- [ ] Tests for all the app
- [X] Improve logs
- [X] Makefile
- [x] Improve naming


### Start rabbitmq image

docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management

Dashboard can be seen in http://localhost:15672/