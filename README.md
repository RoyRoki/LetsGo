# Random Chat App (Omegle-Like) - Scalable & Secure

## **Overview**
This project aims to build an **infinitely scalable** random chat application similar to Omegle. The app will support **text chat initially**, with plans to add **audio/video calls, user monitoring, and advanced matching (gender & tags)** in future versions. The architecture follows **SOLID principles** and **Clean Architecture** for maintainability and scalability.

## **Technologies Used**
| Component            | Technology |
|---------------------|------------|
| **Backend (Core Service)** | Go (Golang) |
| **Real-time Communication** | WebSockets + Redis Pub/Sub |
| **Matchmaking System** | Redis Sorted Sets |
| **Data Storage** | PostgreSQL / MongoDB |
| **Messaging Queue** | Redis / RabbitMQ |
| **Scaling** | Kubernetes + Load Balancing |
| **WebRTC (Future - Audio/Video Calls)** | WebRTC + Signaling Server |
| **Monitoring & Analytics** | Prometheus, Grafana, Elasticsearch |
| **Security Measures** | IP Rate Limiting, WebSocket Tokens, Encryption |

## **Project Architecture**
### **1. Layered Clean Architecture**
- **Presentation Layer** → WebSocket API (User Connection & Messaging)
- **Application Layer** → Business Logic (Pairing, Monitoring, Session Handling)
- **Domain Layer** → Core Entities & Use Cases (User, Session, Matching)
- **Infrastructure Layer** → Database, Redis, WebRTC Signaling

### **2. Scalable WebSocket Handling**
- Use **NGINX/WebSocket Load Balancer** for distributing WebSocket connections.
- Each WebSocket server instance is **stateless**, storing session data in **Redis**.

### **3. Matchmaking System**
- **Redis Sorted Sets** are used to match users efficiently.
- Users are stored with timestamps to match in a **FIFO manner**.
- Future expansion: Match users based on **gender & tags**.

### **4. Messaging System**
- **Redis Pub/Sub** ensures real-time communication across WebSocket instances.
- Future expansion: Support **RabbitMQ/Kafka** for scalable message relays.

### **5. Security Measures**
- **Rate Limiting** (Nginx + Redis) prevents spam & abuse.
- **WebSocket Token Authentication** ensures session integrity.
- **End-to-End Encryption (E2EE) for WebRTC Calls** (Future Feature).

## **Implementation Roadmap**
### **Phase 1: Core Text Chat System** ✅
1. **Set up WebSocket server** in Go.
2. **Implement matchmaking** using Redis Sorted Sets.
3. **Use Redis Pub/Sub** for message delivery.
4. **Rate limit & secure WebSockets.**

### **Phase 2: Scaling & Monitoring** 🔄
5. Deploy **NGINX Load Balancer** for WebSocket connections.
6. Use **PostgreSQL/MongoDB** to store chat metadata.
7. Implement **Prometheus & Grafana** for real-time analytics.

### **Phase 3: Advanced Matchmaking & Filtering** 🔜
8. Extend matchmaking to support **gender & tag-based pairing**.
9. Optimize Redis queries for faster matching.

### **Phase 4: WebRTC for Audio & Video Calls** 🔜
10. Implement **WebRTC signaling server** in Go.
11. Allow P2P connections for video/audio (or SFU for scalability).

## **How to Run the Project**
### **1. Install Dependencies**
```sh
sudo apt update && sudo apt install redis postgresql nginx
```

### **2. Run Redis Server**
```sh
redis-server
```

### **3. Run WebSocket Server**
```sh
go run main.go
```

### **4. Set Up Nginx as WebSocket Load Balancer**
```sh
sudo cp nginx.conf /etc/nginx/sites-enabled/
sudo systemctl restart nginx
```

## **Future Enhancements**
✅ Multi-region scaling with Kubernetes
✅ End-to-End Encryption for WebRTC Calls
✅ AI-based smart matchmaking system
✅ Real-time message moderation & filtering

---
🚀 **This README will be updated as the project evolves!**



internal/config: Use this for configuration logic, such as loading environment variables, parsing config files at runtime, and handling configuration dependencies within your application.
config/ (Root Level): Use this for static configuration files (e.g., .env, config.json, config.yaml) that your application reads at runtime.

