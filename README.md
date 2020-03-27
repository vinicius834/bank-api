# Bank API

##  API Structure

```bash
bank-api
|-- modules(account, trasaction) // You can create N modules you wish
|-- config  // Set the basic configuration like router
|-- helper  // Common generic constants, common functions and other stuff used in the whole api
|-- storage // Database configs and start
|-- init.go // Api start
```

## Run

Supose you are in /bank-api directory:
```bash
docker-compose up --build --force-recreate
```

## Test
### OBS: Containers must to be up
Supose you are in /bank-api directory:
```bash
    go test ./...
``` 

## Usage

### Create account
```bash
POST /accounts
{
	"documentNumber": "<document-number."
}
```
### Get account by id
```bash
GET /accounts/{id}
```
### Create OperationType
```bash
POST /transactions
{
	"description": "<name or description>",
	"isCredit": <true or false>
}
```
### Create Transaction
```bash
POST /transactions
{
	"accountID": "<account-id>",
	"operationTypeID": "<operation-type>",
	"amount": <amount>
}
```