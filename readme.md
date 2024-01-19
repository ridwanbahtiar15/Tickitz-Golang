# Backend Tickitz with Golang

<div align="center">
  <img src="https://res.cloudinary.com/doncmmfaa/image/upload/v1705476761/samples/Tickitz_1_qjg2bh.png" alt="Logo" />
</div>

This project is about to show you on my performance in developing backend architecture using Golang. It has couple of features and API also several security authorization. It is a website for purchasing cinema tickets with main features including a list of films and their details, ordering cinema tickets based on the desired time and place. There are 2 roles, namely Consumer and Admin. Its has authentication and authorization for several accessible pages based on role.

## Technologies used in this project

- Gin Gonic \
  Gin Gonic is a lightweight and fast web framework for Golang. \
  [Gin Gonic Documentation](https://pkg.go.dev/github.com/gin-gonic/gin#section-readme)

- JSON Web Token \
  JSON Web Tokens provide a secure and compact way to transmit information between parties. \
  [JSON Web Token](https://jwt.io/introduction)

- Cloudinary \
  Cloudinary is a cloud-based service for managing and optimizing images and videos. \
  [Cloudinary Documentation](https://cloudinary.com/documentation)

- Midtrans \
  Midtrans is a payment gateway service that simplifies online transactions. \
  [Midtrans Documentation](https://docs.midtrans.com/)

- Govalidator \
  Govalidator is a versatile validation library for Golang. \
  [Govalidator Documentation](https://github.com/asaskevich/govalidator)

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

```bash
  DB_HOST = "YOUR DB_HOST"
  DB_NAME = "YOUR DB_NAME"
  DB_USER = "YOUR DB_USER"
  DB_PASSWORD = "YOUR DB_PASSWORD"
  JWT_KEY = "YOUR JWT_KEY"
  ISSUER = "YOUR ISSUER"
  CLOUDINARY_NAME = "YOUR CLOUDINARY_NAME"
  CLOUDINARY_KEY = "YOUR CLOUDINARY_KEY"
  CLOUDINARY_SECRET = "YOUR CLOUDINARY_SECRET"
  MIDTRANS_ID_MERCHANT = "YOUR MIDTRANS_ID_MERCHANT"
  MIDTRANS_CLIENT_KEY = "YOUR MIDTRANS_CLIENT_KEY"
  MIDTRANS_SERVER_KEY = "YOUR MIDTRANS_SERVER_KEY"
```

## Run Locally

Clone the project

```bash
  $ git clone https://github.com/GilangRizaltin/Tickitz-Golang
```

Go to the project directory

```bash
  $ cd Tickitz-Golang
```

Install dependencies

```bash
  $ go get .
```

Start the server

```bash
  $ go run ./cmd/main.go
```

## Running Tests

To run tests, run the following command

```bash
  $ go test .
```

## API Reference

#### Authentication & Authorization

```http
  /auth
```

| Method   | Endpoint      | Description                        |
| :------- | :------------ | :--------------------------------- |
| `POST`   | `"/register"` | register user                      |
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

| Method | Endpoint                | Description                                 |
| :----- | :---------------------- | :------------------------------------------ |
| `GET`  | `"/"`                   | GET all movie                               |
| `GET`  | `"/movie/:movie_id"`    | Get movie details **Required** movie_id     |
| `GET`  | `"/movie/:schedule_id"` | Get movie schedule **Required** schedule_id |

#### Orders

```http
  /order
```

| Method | Endpoint          | Description                                           |
| :----- | :---------------- | :---------------------------------------------------- |
| `GET`  | `"/"`             | Get orders per users **Required** authorization token |
| `POST` | `"/"`             | Create transaction order                              |
| `POST` | `"/notification"` | Push notification for sucessfull payment              |

## Documentation

[Postman Documentation](https://documenter.getpostman.com/view/29696636/2s9Ykn8MDe)

## Related Project

[Front End (React Js)](https://github.com/GilangRizaltin/Tickitz-Frontend-Gilang)

## Collaborator

- [@Ridwan Bahtiar](https://github.com/ridwanbahtiar15)

## Support

For support, email gilangzaltin@gmail
