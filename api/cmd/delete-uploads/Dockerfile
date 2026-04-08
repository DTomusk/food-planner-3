# Stage 1: Build binaries 
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN ls -la

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o /bin/food-planner-delete-uploads ./cmd/delete-uploads

# Stage 2: Create the final image
FROM alpine:3.22.2

ENV ENV=docker

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /bin/food-planner-delete-uploads /usr/local/bin/

COPY --from=builder /app/reference ./reference

#RUN chmod +x /usr/local/bin/food-planner-delete-uploads

CMD ["/usr/local/bin/food-planner-delete-uploads"]