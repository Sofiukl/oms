# Order Management System (OMS)
This repository contains sample code of Order Management System. The code is written keeping focus on
concurrency handling.

# Modules
It has the below modules which can be run as separate microservice.
product: Implementation of product related APIs [port 3004]
Cart: Implementation of cart [port 3006]
Checkout: Implementation of checkout module [port 3005]

All the microservices uses postgress database. Database URL is mentioned in respective
service app.env file using key DBURL. app.env file also contain the application server port.

Note: Different business cases, business validations, error handling and code comments are not in place.
We are focing here on concurrency handling.

# Architrcture
To Be Updated


# SQL Script
Please refer to the file postgress.sql. Database and tables need to created before
running the application.


# API Documentation
1. Checkout API
<pre>POST
http://localhost:3005/checkout-service/api/v1/checkout/
Request Body
    <pre> {
        "cart_id": "c1",
        "amount": 100
    } </pre>

Response: 
    Status: 201 (checkout will be processed internally) </pre>

Others to be updated


# API Gateway
NGINX is used for api gateway and load balance. Configuration can be found in folder
nginx in root
