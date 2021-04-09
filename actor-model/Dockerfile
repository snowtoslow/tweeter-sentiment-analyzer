FROM golang:1.16

LABEL maintainer="snowtoslow <snowtoslow@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/tweeter-sentiment-analyzer/actor-model

COPY go.mod go.sum ./

RUN go mod download

# Copy everything from the current directory to the PWD (Present Working Directory) inside
COPY . .

# Build the Go app
RUN go build -o main .

# Command to run the executable
# CMD ["./main"]
ENTRYPOINT ["./main"]
