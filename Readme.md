
### Clone the 3 services listed below
https://github.com/toluwase1/userservice
https://github.com/toluwase1/companyservice
https://github.com/toluwase1/notification

### CD into each of the 3 clones
### Run docker compose up for each

### User service
#### example payload for the user service signup: http://localhost:8081/api/v1/auth/signup
{
"name":"toluwase",
"email":"tolu@gmail.com.com",
"phone_number":"+2348900989156",
"password":"toluwase"
}

#### example payload for the user service login: http://localhost:8081/api/v1/auth/login
{
"email":"tolu@gmail.com.com",
"password":"toluwase"
}

### Company service
#### example payload for create company service: http://localhost:8080/api/v1/company/
{
"id":"1",
"support_email" : "any@gmail.com",
"name" : "new company",
"description" : "my new company 1",
"amount_of_employees" : 6,
"is_registered" : true,
"type" : "NonProfit"
}

#### example payload for update company service: http://localhost:8080/api/v1/company
{
"support_email" : "any@gmail.com",
"name" : "new company",
"description" : "my new coy",
"amount_of_employees" : 6,
"is_registered" : false,
"type" : "NonProfit"
}

#### example request for delete company service:
http://localhost:8080/api/v1/company/{uuid}
