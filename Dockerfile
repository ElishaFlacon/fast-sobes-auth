# Initial stage: download modules
FROM golang:1.25-alpine

# Install air for hot reloading
RUN go install github.com/air-verse/air@v1.63.0

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Set environment
ENV config=docker

# Expose ports
EXPOSE 5000
EXPOSE 5555
EXPOSE 7070

ENV PATH="/go/bin:${PATH}"

# Use air for hot reload
CMD ["air", "-c", ".air.toml"]