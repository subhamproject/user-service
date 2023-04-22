# user-service
This repo will hold data related to Test Service

### To get User
http://localhost:8082/user?id=100

### To get user with order
http://localhost:8082/user/order?id=101

### Environment Variables
| Name | type | default value   | Description |
| :---: | :---:  | :---: | :---: |
| SERVICE_PORT | string | 8082 |  User Service http server port |
| ORDER_SVC_HOST | string | localhost | Order Service hostname |
| ORDER_SVC_PORT | string | 8081   | Order Service port number |