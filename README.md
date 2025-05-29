<h1 align="center">
  <img src="https://img.shields.io/badge/LetsGo-Golang%20Tag%20Chatting-blue?style=for-the-badge&logo=go" height="40" />
  <br />
  🧠 LetsGo – Tag-Based Chatting App
</h1>

<p align="center">
  <img src="https://img.shields.io/github/stars/RoyRoki/LetsGo?style=flat-square" />
  <img src="https://img.shields.io/github/issues/RoyRoki/LetsGo?style=flat-square" />
  <img src="https://img.shields.io/badge/built%20with-Go%201.22-blue?style=flat-square&logo=go" />
  <img src="https://img.shields.io/badge/microservices-yes-success?style=flat-square" />
</p>

<p align="center">
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg" width="40" />
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg" width="40" />
  <img src="https://grpc.io/img/logos/grpc-icon-color.png" width="40" />
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg" width="40" />
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg" width="40" />
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/typescript/typescript-original.svg" width="40" />
  <img src="https://upload.wikimedia.org/wikipedia/commons/c/cf/Lua-Logo.svg" width="40" />
</p>

---

<div align="center">

🎯 <strong>Real-time chat app that matches users based on shared tags using Go, Redis (Lua), WebSockets, and gRPC.</strong>

🎥 <a href="https://youtu.be/l7S9uPnNVvI?si=0xlp8WhQQrIGwv4V">Watch Demo on YouTube</a>

🔗 <a href="https://www.linkedin.com/in/rokiroy/">Connect on LinkedIn</a>

</div>

---

## ✨ Features

* ✅ Tag-based user matching with Redis Bitmaps & Lua scripting
* 🔌 Real-time WebSocket chat between matched users
* 🎯 gRPC-based microservice communication
* 🚀 Clean Architecture + Docker-ready
* 🔄 FIFO fallback matching logic
* 🧩 Ready for scaling with audio/video service modules

---

## 🧠 Architecture

```text
                   +------------------+
                   |  Client (Web)    |
                   +--------+---------+
                            |
                            ▼
                   +--------+---------+
                   |  Chatting Service |
                   | - WebSocket       |
                   | - Usecase Logic   |
                   +--------+---------+
                            |
                            ▼ gRPC
                   +--------+---------+
                   | Matching Service |
                   | - Redis + Lua    |
                   | - Tag Matching   |
                   +------------------+
```

---

## 🛠️ Tech Stack

<table>
<tr>
  <td><img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg" width="30"/> Go 1.22</td>
  <td><img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg" width="30"/> Docker</td>
  <td><img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg" width="30"/> Redis + Lua</td>
  <td><img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg" width="30"/> Node.js</td>
  <td><img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/typescript/typescript-original.svg" width="30"/> TypeScript</td>
  <td><img src="https://grpc.io/img/logos/grpc-icon-color.png" width="30"/> gRPC</td>
  <td><img src="https://upload.wikimedia.org/wikipedia/commons/c/cf/Lua-Logo.svg" width="30"/> Lua</td>
</tr>
</table>

---

## 🚧 Upcoming

* 📹 Video/Audio calling modules via dedicated microservice
* 🌐 Admin dashboard & analytics
* 🔒 Secure match flow with JWT/Auth

---

## 📦 Setup

```bash
git clone https://github.com/RoyRoki/LetsGo
cd LetsGo
docker-compose up --build
```

Open WebSocket tester at [`test/chat_test.html`](./test/chat_test.html) and enjoy!

---

## 📬 Connect

Follow me on [LinkedIn](https://www.linkedin.com/in/rokiroy) for updates & deep dives!

🎥 [Watch Project Demo on YouTube](https://youtu.be/l7S9uPnNVvI?si=0xlp8WhQQrIGwv4V)
