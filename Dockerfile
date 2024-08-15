FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Add verbose output to see what's going wrong
RUN go build -v -o main .

EXPOSE 8080

CMD ["./main"]