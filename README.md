# xm-home-test

- This is a simple API that calculates the minimum number of packs required to fulfill an order.

## Instructions to run the container
### Make sure you have docker and docker-compose installed on your machine.

```sh
  make build
  make up
```

## API
- This is the base url
```sh
   http://localhost:3001
```
- As requested there is an endpoint for getting an authenticated token. So, make sure to pass a **Bearer Token** in the requests.
```sh
    curl --location --request GET 'localhost:3001/auth'
```
- create company
```sh
    curl --location --request POST 'http://localhost:3001/companies' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {YOUR_TOKEN}' \
--data '{
	"name":{COMPANY_NAME},
	"description":{DESCRIPTION},
	"amount_of_employees": {AMOUNT_OF_EMPLOYEES},
	"registered": {true | false},
	"type":{corporations | nonprofit | cooperative | sole proprietorship}
}'
```

- update company
```sh
    curl --location --request PATCH 'http://localhost:3001/companies/{COMPANY_ID}' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {YOUR_TOKEN}' \
--data '{
	"name":{COMPANY_NAME},
	"description":{DESCRIPTION},
	"amount_of_employees": {AMOUNT_OF_EMPLOYEES},
	"registered": {true | false},
	"type":{corporations | nonprofit | cooperative | sole proprietorship}
}'
```
- get one by name
```sh
    curl --location --request GET 'localhost:3001/companies/{COMPANY_NAME}'
```

- delete by id
```sh
    curl --location --request DELETE 'localhost:3001/companies/89b68763-4a5c-4f80-8bda-d1ebfb321eb7' \
--header 'Authorization: Bearer {YOUR_TOKEN}'
```

## Instructions to run the tests
#### It requires the Golang 1.23.x to be installed on the machine.
``` sh
    make unit-test
```

## To run linter
```sh
    make run-linter
```
