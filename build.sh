export CGO_ENABLED=0
export GOOS=linux
export GOARCH=s390x
go build -o goIM .
docker build -t=wright1992/goim:latest .