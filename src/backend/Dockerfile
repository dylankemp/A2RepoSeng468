FROM golang:latest AS build
WORKDIR /go/src/github.com/dylankemp/A2RepoSeng468/
COPY . .

ENV CGO_ENABLED=0
RUN go build -o backend-serv .

FROM build AS development
RUN apt-get update \
    && apt-get install -y git
CMD ["go", "run", "main.go"]
#was 8000
FROM alpine:latest
EXPOSE 8089
COPY --from=build . /backend-serv
RUN chmod 777 -R /backend-serv
CMD ["/backend-serv"]

#build stage
# FROM golang:alpine
# RUN apk add --no-cache git
# WORKDIR /go/src/app
# COPY . .
# EXPOSE 8089
# RUN go get -d -v ./...
# RUN go build -o /go/bin/app -v ./... 

# CMD ["go","run","/go/src/app/main.go"]