# Build the go application
FROM golang:1.23.1-alpine AS build

# Install build dependencies
RUN apk add --no-cache gcc=12.2.1_git20220924-r10 musl-dev=1.2.4-r2

# Set the working directory for the app
WORKDIR /app

# Copy the backend directory
COPY backend/ ./

# Download dependencies
RUN go mod download

# Copy the frontend directory
COPY frontend/ ./frontend/

# Build the app with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o whoknows .

# Create an image for the app - using specific Alpine version
FROM alpine:3.19.1

# Install runtime dependencies
RUN apk add --no-cache ca-certificates=20230506-r0

# Set the working directory for the runtime
WORKDIR /root/

# Copy the go binary from the build stage
COPY --from=build /app/whoknows .

# Copy the frontend directory
COPY --from=build /app/frontend ./frontend

# Copy the database directory
COPY --from=build /app/internal/db ./internal/db

# Expose the port for the application
EXPOSE 8080

# Start the go application
CMD ["./whoknows"]