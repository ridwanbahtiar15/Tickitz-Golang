# Backend Tickitz Team Project with Golang

This project is about to show you on my performance in developing backend architecture using Golang. It has couple of features and API also several security authorization. This project contain JWT, Cloudinary uploader, Hashing password with Argon2d, Docker, Migrations, and Midtrans as payment gateway.

## API Reference

#### Authentication & Authorization

```http
  /auth
```

| Method   | Endpoint      | Description                        |
| :------- | :------------ | :--------------------------------- |
| `POST`   | `"/register"` | register user                      |
| `POST`   | `"/acivate"`  | activating user through OTP        |
| `POST`   | `"/login"`    | get access and identity of user    |
| `DELETE` | `"/logout"`   | delete access and identity of user |

#### Users

```http
  /user
```

| Method  | Endpoint           | Description                     |
| :------ | :----------------- | :------------------------------ |
| `GET`   | `"/profile"`       | Fet user's profile              |
| `POST`  | `"/authorization"` | Checking user's authorization   |
| `PATCH` | `"/"`              | Update users detail and profile |

#### Products

```http
  /movie
```

| Method   | Endpoint                | Description                                       |
| :------- | :---------------------- | :------------------------------------------------ |
| `GET`    | `"/"`                   | GET all movie                                     |
| `GET`    | `"/movie/:movie_id"`    | Get movie details **Required** movie_id           |
| `GET`    | `"/movie/:schedule_id"` | Get movie schedule **Required** schedule_id       |
| `POST`   | `"/"`                   | Add movie (admin only)                            |
| `PATCH`  | `"/:movie_id"`          | Edit a movie **Required** movie_id (admin only)   |
| `DELETE` | `"/:movie_id"`          | Deleting movie **Required** movie_id (admin only) |

#### Orders

```http
  /order
```

| Method | Endpoint     | Description                                           |
| :----- | :----------- | :---------------------------------------------------- |
| `GET`  | `"/"`        | Get orders per users **Required** authorization token |
| `GET`  | `"/stat"`    | Get order statistic (admin only)                      |
| `POST` | `"/"`        | Create transaction order                              |
| `POST` | `"/success"` | Push notification for sucessfull payment              |
| `POST` | `"/failed"`  | Push notification for failed or expired payment       |

## Authors

Authors By Me _as known as_ Gilang Muhamad Rizaltin \
**Github link**

- [@Gilang Rizaltin](https://github.com/GilangRizaltin)
