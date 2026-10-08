# [Backend] - EventHub API v1.0.0

<p>
  <a href="https://golang.org/"><img style="border-radius:4px" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go"></a>
  <a href="https://github.com/gin-gonic/gin"><img style="border-radius:4px" src="https://img.shields.io/badge/Gin-008B8B?style=for-the-badge&logo=go&logoColor=white" alt="Gin"></a>
  <a href="https://github.com/joho/godotenv"><img style="border-radius:4px" src="https://img.shields.io/badge/Godotenv-ECD53F?style=for-the-badge&logo=dotenv&logoColor=black" alt="Godotenv"></a>
  <a href="https://www.postgresql.org/"><img style="border-radius:4px" src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL"></a>
  <a href="https://github.com/jackc/pgx"><img style="border-radius:4px" src="https://img.shields.io/badge/pgx-339933?style=for-the-badge&logo=postgresql&logoColor=white" alt="pgx"></a>
  <a href="https://redis.io/"><img style="border-radius:4px" src="https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white" alt="Redis"></a>
</p>

A robust backend RESTful API built with Go and the Gin Framework, developed as part of the Koda Batch 9 curriculum. This project serves as a comprehensive backend service for managing community-driven events, user memberships, and authentication.

## Tech Stack

- [![Go](https://img.shields.io/badge/Go-v1.27.1-blue?logo=go&logoColor=white)](https://go.dev/learn/) <br>Digunakan sebagai bahasa pemrograman utama untuk membangun seluruh logika backend aplikasi _Gin Exercise Koda Batch 9_ ini.

- [![Gin-Gonic](https://img.shields.io/badge/Gin_Gonic-v1.12.0-green?logo=Gin&logoColor=white)](https://gin-gonic.com/en/)<br>Dipakai sebagai _HTTP router_ dan web framework untuk mengatur _endpoint_ RESTful API fitur autentikasi, manajemen _user_, komunitas, dan _events_.
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v18.6-blue?logo=PostgreSQL&logoColor=white)](https://www.postgresql.org/) <br>Menjadi basis data relasional utama untuk menyimpan data aplikasi secara permanen seperti akun pengguna, profil, daftar komunitas, detail _events_, serta relasi _join/saved events_.
- [![pgx](https://img.shields.io/badge/Pgx-v5.11.0-blue?)](https://github.com/jackc/pgx) <br>Digunakan sebagai _driver_ dan _toolkit_ penghubung antara kode Go dan PostgreSQL untuk mengeksekusi _query_ database secara cepat dan efisien pada fitur-fitur aplikasi.
- [![Redis](https://img.shields.io/badge/Redis-v8.10.2-blue?logo=go&logoColor=white)](https://redis.io/) <br>Diterapkan untuk mengelola sesi pengguna, seperti penyimpanan token autentikasi atau _blacklist_ token saat _logout_ dari aplikasi.
- [![GodotEnv](https://img.shields.io/badge/GodotEnv-v1.5.1-blue?)](https://github.com/jackc/pgx) <br>Berfungsi untuk membaca file konfigurasi `.env` guna mengelola variabel rahasia dan dinamis (seperti koneksi database PostgreSQL, port server, dan kunci rahasia token) di aplikasi.

## Features

- Authentication (Register, Login, Logout)
- User (Join Communities, Join Events, List Events Joined, List Events Saved, etc.)
- Events (List Events, Filter Events, Detail Events)
- Communities (List Communities, Filter Communities, Detail Communities)

## Usage Instruction

### Clone this repository

```bash
$ git clone https://github.com/almarufh/koda-b8-backend-eventhub project-eventhub
```

### Environment Setup

Create your environment on the root directory named **.env** example use [.env.example](https://github.com/almarufh/koda-b9-backend-eventhub/blob/main/.env.example)

### Redis Config Setup

Setup your configuration for your use cache on the root directory named **redis.conf** example [redis.conf](https://github.com/almarufh/koda-b9-backend-eventhub/blob/main/redis.conf.example)

### Docker Pull Image from [**DokerHub**](https://hub.docker.com/)

1. PosgresSQL [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v18.6-blue?logo=PostgreSQL&logoColor=white)](https://www.postgresql.org/)

```bash
$ docker pull postgres:18.6-alpine3.24
```

2. Redis [![Redis](https://img.shields.io/badge/Redis-v8.10.2-blue?logo=go&logoColor=white)](https://redis.io/)

```bash
$ docker pull redis:8.10.2-alpine3.23
```

### Running Database

1. PosgresSQL [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v18.6-blue?logo=PostgreSQL&logoColor=white)](https://www.postgresql.org/)

```bash
$ docker run --name {NAME_CONTAINER} -d --env-file .env -p {PORT_FORWARD} -v {YOUR_DOCKER_VOLUME}:/var/lib/postgresql/18/docker postgres:18.6-alpine3.24
```

2. Redis [![Redis](https://img.shields.io/badge/Redis-v8.10.2-blue?logo=go&logoColor=white)](https://redis.io/)

```bash
$ docker run -d --name {NAME_CONTAINER} -p {PORT-FORWARDING} -v {PATH_CONFIG}:/usr/local/etc/redis/redis.conf" -v "{DOCKER_VOLUME}:/data" --restart unless-stopped redis:8.8.3-alpine3.23 redis-server /usr/local/etc/redis/redis.conf
```

3. Create Schema Table and used data dummy

```zsh
$ make migrate-up
$ make db-seed
```

4. Install dependency

```bash
$ go mod download
```

3. etc.

## Routes

| Endpoint                    | Method   | Description            |
| :-------------------------- | :------- | :--------------------- |
| `/auth/register`            | `POST`   | Create New Account     |
| `/auth/login`               | `POST`   | Log In                 |
| `/auth/logout`              | `DELETE` | Log Out                |
| `/auth/forgot/password`     | `PATCH`  | Forgot Password        |
| `/communities/{id}`         | `GET`    | Get Detail Community   |
| `/communities/{id}/join`    | `POST`   | Join Community         |
| `/communities/{id}/leave`   | `DELETE` | Leave Community        |
| `/communities/{id}/members` | `GET`    | Get Community Members  |
| `/events/{id}`              | `GET`    | Get Detail Event       |
| `/events/{id}/join`         | `POST`   | Join Event             |
| `/events/{id}/leave`        | `DELETE` | Leave Event            |
| `/events/{id}/save`         | `POST`   | Save Event             |
| `/events/{id}/unsave`       | `DELETE` | Unsave Event           |
| `/user`                     | `GET`    | Get Profiles           |
| `/user`                     | `PATCH`  | Update Profile         |
| `/user/password`            | `PATCH`  | Change Password        |
| `/user/events/joined`       | `GET`    | Get Joined Events      |
| `/user/events/saved`        | `GET`    | Get Saved Events       |
| `/user/community/joined`    | `GET`    | Get Joined Communities |

### Documentation

For complete documentation, visit [docs swagger endpooint **`/docs/v1/index.html`** ](http://localhost/docs/v1/index.html)

## Changelog

| Version | Desctiption |
| ------- | ----------- |
| 1.0.0   | Initial app |

- 1.1.1
  - Feat cors .... by [almarufh](https://github.com/almarufh)
  - Feat hash password .... by [almarufh](https://github.com/almarufh)
  - Feat Autentication .... by [almarufh](https://github.com/almarufh)

## How to Contribute

- Fork this repository
- Create your changes
- Pull Request

## License [![License: MIT](https://img.shields.io/badge/License-MIT-red)](https://opensource.org/license/mit)

This project is licensed under the MIT License

## CONTACTS

[![Email](https://img.shields.io/badge/Email-almarufhidayat99@gmail.com-blue?logo=maildotru&logoColor=red)](mailto:almarufhidayat99@gmail.com)<br>
[![WhatsApp](https://img.shields.io/badge/WhatsApp-6281973779380-blue?logo=whatsapp&logoColor=green)](https://wa.me/6281973779380)

## Related Project

[![WhatsApp](https://img.shields.io/badge/FRONTEND--blue?logo=googlechrome&logoColor=blue)](https://github.com/almarufh)
