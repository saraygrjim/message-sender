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
- [ ] Serve test html in a port


### Start rabbitmq image

docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management

Dashboard can be seen in http://localhost:15672/