  FROM golang:1.24.5-alpine 

  WORKDIR /app

  COPY server/go.mod ./
  COPY server/go.sum ./
  RUN go mod download

  COPY . ./

  RUN go build -o bin/server ./server/

  RUN ls -la ./bin/server

  RUN chmod +x bin/server

  EXPOSE 8080

  ENTRYPOINT ["./bin/server"]
