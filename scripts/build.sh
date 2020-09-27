mkdir -p ./dist
go build -ldflags="-s -w" -o ./dist/liege ./liege.go
