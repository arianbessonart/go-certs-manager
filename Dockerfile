FROM golang

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN go build -mod vendor -o cert-manager .

EXPOSE 8080

CMD [ "./cert-manager" ]
