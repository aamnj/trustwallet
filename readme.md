## OverView
This documents describe how to setup this project in local environment. Along with API documentation
### Features
- APIs to fetch current block, subcribe to an address and fetching all transactions for an address
- Job to fetch transactions starting from current block and then fetching all transactions for each block in specified time interval
- Once transactions are fetched for each block, stores the transaction for subscribed address. Transactions for other addresses are not stored

## Get Started
### Install Go
Install using brew for mac https://formulae.brew.sh/formula/go

### Project Setup
- Clone the repository
- Once you have go installed, run project using `make run`

## API Documentation

### Get Current Block
#### Request
```
curl --location 'http://localhost:8080/current_block'
```
#### Response
```
{
  "status": 200,
  "data": {
      "block": 20898910
}
```

### Subscribe
#### Request
```
curl --location 'http://localhost:8080/subscribe' \
--header 'Content-Type: application/json' \
--data '{
    "address": "0x560075ba1990dded3b8cb49f1955e607999c942d"
}'
```
#### Success Response
```
{
    "status": 200,
    "message": "Address 0xc7bcbb3f18b5478e275b38067f0ed633895b2b5b subscribed"
}
```
#### Error Response
```
{
    "error": "Already subscribed to address 0x560075ba1990dded3b8cb49f1955e607999c942d",
    "status": 400
}

```

### Get Transactions
#### Request
```
curl --location --request GET 'http://localhost:8080/transactions?address=0xc7bcbb3f18b5478e275b38067f0ed633895b2b5b' \
--header 'Content-Type: application/json' \
--data '{
    "address": "0xc7bcbb3f18b5478e275b38067f0ed633895b2b5b"
}'
```
#### Success Response
```
{
    "status": 200
    "data":[{
      "from": "",
      "to": "",
      "value": "",
      "hash": ""
    }]
}
```
#### Error Response
```
{
    "error": "No transactions found for this address",
    "status": 400
}
```

