FROM golang:1.16

LABEL maintainer="snowtoslow <snowtoslow@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/tweeter-sentiment-analyzer/message-broker

# Copy everything from the current directory to the PWD (Present Working Directory) inside
COPY . .

# Build the Go app
RUN go build -o mainbroker .

EXPOSE 8088

# Command to run the executable
# CMD ["./mainbroker"]
ENTRYPOINT ["./mainbroker"]
