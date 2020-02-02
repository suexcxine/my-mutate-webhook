all: app

app:
	CGO_ENABLED=0 GOOS=linux go build -o my-mutate cmd/main.go
