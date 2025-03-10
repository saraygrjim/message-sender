

# 🚀 Message Sender

This project implements a real-time messaging system using **WebSockets** and **RabbitMQ**, consisting of two microservices and a demo subscriber.

## 📌 Architecture

The system consists of three main components:

### 1️⃣ **Receiver (Microservice 1)**
- Listens for messages sent via **WebSocket** at `ws://localhost:8080/ws`.
- Publishes received messages to a **RabbitMQ queue** (`localhost:5672`).
- Serves a **web interface at `http://localhost:9090/`** for sending test messages.

### 2️⃣ **Broadcaster (Microservice 2)**
- Reads messages from the **RabbitMQ queue**.
- Forwards them to **all connected subscribers** via WebSocket at `ws://localhost:8081/echo`.

### 3️⃣ **Subscriber**
- Connects to `ws://localhost:8081/echo` and listens for messages sent by the Broadcaster.
- Multiple instances of this component can run simultaneously.

---

## 📂 Project Structure

```
/message-sender
├── /cmd            # Code to start microservices
├── /internal       # Common components used in multiple services
│ ├── rabbitmq      # RabbitMQ wrapper
│ └── websocket     # WebSocket wrapper
├── microservices   # Source code for the microservices
│ ├── broadcaster   # Source code for the broadcaster
│ ├── receiver      # Source code for the receiver
│ └── subscriber    # Source code for the subscriber
├── static          # Web interface files (HTML, JS, CSS)
├── Makefile        # Commands for running the services
└── README.md       # Project documentation
```

---

## 🛠 Technologies Used

- **Go** (Golang)
- **RabbitMQ** (Message queue)
- **WebSockets** (Real-time communication)
- **Gorilla WebSocket** (`github.com/gorilla/websocket`)
- **RabbitMQ AMQP** (`github.com/rabbitmq/amqp091-go`)

---

## 🚀 How to Run the Project

### 1️⃣ Start RabbitMQ

Run the following command to launch a RabbitMQ container:

```sh
make provision
```

📌 *The RabbitMQ management dashboard is available at* [http://localhost:15672/](http://localhost:15672/)  
*(Username: guest | Password: guest)*

### 2️⃣ Start the **Receiver** Microservice
```sh
make receiver
```
📌 *Listens for WebSocket connections at `ws://localhost:8080/ws`*  
📌 *Web interface available at* [http://localhost:9090/](http://localhost:9090/)

### 3️⃣ Start the **Broadcaster** Microservice
```sh
make broadcaster
```
📌 *Reads messages from RabbitMQ and broadcasts them via WebSocket at `ws://localhost:8081/echo`*

### 4️⃣ Start a **Subscriber**
```sh
make subscriber
```
📌 *Listens for messages from `ws://localhost:8081/echo`*  
📌 *You can run multiple instances of this command.*

---

## 🧪 Running Tests

To run the tests:

```sh
go test ./...
```

**Testify** (`github.com/stretchr/testify`) is used for unit testing.

---

## 📜 License

This project is licensed under the **MIT License**. Feel free to contribute! 🚀