
receiver:
	go run main.go receiver

broadcaster:
	go run main.go broadcaster

subscriber:
	go run main.go subscriber

provision:
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management

remove:
	docker stop rabbitmq &&	docker rm rabbitmq