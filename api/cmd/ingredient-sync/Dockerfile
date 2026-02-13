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
RUN go build -o /bin/food-planner-ingredient-sync ./cmd/ingredient-sync

# Stage 2: Create the final image
FROM alpine:3.22.2

ENV ENV=docker

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /bin/food-planner-ingredient-sync /usr/local/bin/

COPY --from=builder /app/reference ./reference

#RUN chmod +x /usr/local/bin/food-planner-ingredient-sync

CMD ["/usr/local/bin/food-planner-ingredient-sync"]