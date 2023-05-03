# user-service
This repo will hold data related to Test Service

### To Create new User
`curl --location --request POST 'http://localhost:8082/user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name" : "DevopsDemo"
}'`

### To get all Users
http://localhost:8082/users


### To get User by Id
http://localhost:8082/user?id=100

### To get user with order
http://localhost:8082/user/order?id=100


### Environment Variables
| Name | type | default value   | Description |
| :---: | :---:  | :---: | :---: |
| SERVICE_PORT | string | 8082 |  User Service http server port |
| ORDER_SVC_HOST | string | localhost | Order Service hostname |
| ORDER_SVC_PORT | string | 8081   | Order Service port number |