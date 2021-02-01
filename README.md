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

# Check Out Process Details (Algorithm)
1. Start transaction
2. Find cart details calling cart microservice
3. Acquire lock 
4. Check order quantity is available or not using below snipet
    if (available quanity - reserve quantity < order quantity) then product "out of stock"
5. If order quantity available then proceed
    5.1 update reserve quantity += order quantity
    5.2 process payment
    5.3 if payment success then decrease available quantity as below
        reserve quantity -= order quantity
        available quantity -= order quantity
    5.4 if payment fails then revert reserve quanity as below
        reserve quantity -= order quantity
6. Release lock
7. Commit transaction


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
nginx in root.