# Backend Tickitz with Golang

<p align="center">
        <img src="https://res.cloudinary.com/doncmmfaa/image/upload/v1705476761/samples/Tickitz_1_qjg2bh.png" width="200px" alt="logo"></img>
</p>

A web api project for ordering ticket movie online. There are 4 operations that can be performed, Get (fetching data), Post (insert data), Update (update partial data), delete (delete data).

## Built With

- [Golang](https://go.dev/)
- [Postgre SQL](https://www.postgresql.org/)
- [GinGonic](https://gin-gonic.com/)
- [GoValidator](https://github.com/asaskevich/govalidator)
- [Cloudinary](https://github.com/cloudinary/cloudinary-go)

## Configure app

Create file `.env` then edit it with your settings
according to your needs. You will need:

DB_HOST = Your Database Host
DB_NAME = Your Database Host
DB_USER = Your Database User
DB_PASS = Your Database Password
JWT_KEY = Your JWT Key
ISSUER = Your Issuer
CLOUDINERY_NAME = Your Cloudinary Name
CLOUDINERY_KEY = Your Cloudinary Key
CLOUDINERY_SECRET = Your Cloudinary Secret
MIDTRANS_ID_MERCHANT = Your Midtrans ID Merchant
MIDTRANS_CLIENT_KEY = Your Midtrans Client Key
MIDTRANS_SERVER_KEY =  Your Midtrans Server Key

## Install And Run Locally

1.  Clone project from github repository

        $ git clone https://github.com/ridwanbahtiar15/Tickitz-Golang

2.  go to folder coffee-shop

        $ cd coffee-shop-golang

3.  install dependencies

        $ go get .

4.  Start the server

        $ go run ./cmd/main.go

## API Reference

Auth
| Route | Method | Description |
| -------------- | ----------------------- | ------ |
| /auth/login | POST | Login user |
| /auth/register | POST | Register user |
| /auth/logout | DELETE | Logout user |

Users
| Route | Method | Description |
| -------------- | ----------------------- | ------ |
| /user/login | GET | Login user |
| /user/authorization | POST | Checking user registration |
| /user/ | PATCH | Update user |

Movies
| Route | Method | Description |
| -------------- | ----------------------- | ------ |
| /movie/ | GET | Get all movie |
| /movie/:movie_id | GET | Get movie detail by id |
| /movie/:schedule_id | GET | Get movie schedule |

Orders
| Route | Method | Description |
| -------------- | ----------------------- | ------ |
| /orders/ | GET | Get all order |
| /orders/ | POST | Create transaction |
| /orders/notification | GET | Push notification for sucessfull payment |

## Documentation

[Postman Documentation](https://documenter.getpostman.com/view/28541505/2s9YyqhMwh)

## Related Project

[Front End with React JS](https://github.com/ridwanbahtiar15/Tickitz-Frontend)

## Collaborator

- [Gilang Rizaltin](https://github.com/GilangRizaltin)
