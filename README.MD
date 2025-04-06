# About Operation-Borderless
This is a FinTech App that mimics some of the key features of how a real FinTech works. It uses your registered email address to send money.

## Tech Stack Used
- Golang (Backend)
- React (Frontend)

## Features Walk Through
No authentication was done. Anybody with a valid email can create a wallet or use the app. 

Features Implemented are
- Create/Retrieve Wallet: Home page gives a box for you to put your email address (for new or returning users). If you don't have an existing wallet, it creates one for you and redirects you to your dashboard. If you do, it just gets your wallet and redirects you to your dashboard.
- Deposit: Allows you to deposit money in our currencies (cUSD, cNGN, cXAF, cEUR).
- View Wallet (Balance): Allows you to access to your wallet to view your wallet balances and also total balance in cUSD.
- Swap Currencies: Allows you to swap your currencies. You can deposit in cUSD and swap or convert a portion to cNGN.
- Transfer: Using a user's email address, it allows you to send money to the user in any currency. Eg. You can have money in cUSD and send it to another user in cEUR. It does the conversion automatically using real-live rate tracker.
- View Transactions: Allows you to view all your transactions including deposits, swaps, and transfers

- N/B: Audit feature is implemented internally to save the logs of users (userID, IP Address, Device used, Country, Browser and Timestamp) that used the app in the database and log in. You can only access this if you have access to the database.


## Project Parts
The Project has 2 parts:
- Backend Implementation: This (borderless-money) repository
- Frontend Implementation: https://github.com/Kjasuquo/borderless-money-frontend

## Backend
### Major API Integration
- FOREX_API_KEY: used to get real live rates. (My own key provided in the .env.sample). It is advisable you create your own key to test locally with

Testing deployed backend can be done through `https://backend-borderless-money.xyz/`

### Deployment
The was done on GCP kubernetes cluster. And `https://backend-borderless-money.xyz/` was acquired from namecheap inorder to secure it with https

### Installation
Before installing, download and install go version 1.23.0 and above from **[golang's official website](https://go.dev/doc/install)**
Ensure you have docker if you are running it on docker.

### Install make 
for MacOS run

```brew install make```

### Make available the correct set of dependencies 
```go mod tidy``` 

### Starting Server Locally
- create a `.env` file in the base directory with sample fields from `.env.sample`
- run `make run` to run the app in docker.

## Frontend
This was implemented using React.js

### Deployment
The was done on GCP kubernetes cluster. And `https://borderless-money.xyz/create-wallet` was acquired from namecheap inorder to secure it with https and connect to the backend https

### Starting Frontend Locally
- npm install
- npm start

But note that if you run the React app locally, you won't be able to make calls to the https urls. You would need to point the axios url to your local backend endpoint/localhost

# WARNING
The GCP Kubernetes Cluster and provisioned database (for the backend) will be shut down after this codebase achieves its purpose. So for clarification contact `okoasuquo@gmail.com` 
