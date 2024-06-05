balancer:
	docker build -t user -f balancer.Dockerfile .
	docker run --rm -e ADDRESS=0.0.0.0:8081 --network testnet -p 8081:8081 user   

database:
	docker run --rm --name postgres  --network testnet -p 5432:5432  -e POSTGRES_PASSWORD=1234 -d postgres 
.PHONY: balancer database