### TODO:

- [X] Client who send the messages
- [X] Microservice 1
    - [X] Listen messages from websocket
    - [X] Send messages to queue
- [X] Microservice 2
    - [X] Read messages from queue
    - [X] Send messages to subscriptions
- [X] Subscribers
    - [X] Subscribes to messages
    - [X] Listen messages
    - [X] Print messages
- [ ] Queue
    - [X] Basic connect, read and send
    - [ ] Configure security and more parameters
- [X] Organize project structure following clean architecture
- [X] Cmd
- [X] Graph to initiate the app
- [X] Tests for all the app
    - [X] Microservice 1
    - [X] Microservice 2
    - [X] Subscriber
- [X] Improve logs
- [X] Makefile
- [X] Improve naming
- [X] Serve test html in a port
- [ ] Channels to listen the go routines errors 
- [ ] Review naming

### Start rabbitmq image

docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management

Dashboard can be seen in http://localhost:15672/