run:
	docker-compose up -d

env:
	@test -f .env || cp .env.example .env
		
vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: vet lint

test:
	go test ./...

clean:
	docker-compose down --rmi all --volumes --remove-orphans