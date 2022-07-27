# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:alpine as base

# Add Maintainer Info
LABEL maintainer="Domingo Sanz Marti <domingosanzmarti@gmail.com>"

RUN GOCACHE=OFF

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

FROM base as test

CMD ["sh scripts/test.sh"]

FROM base as build

# Build the Go app
RUN go build -o cmd/server/main cmd/server/main.go

FROM golang:alpine as production

COPY --from=build /app/cmd/server/main /app/main 
COPY --from=build /app/config/.. /app/config/.. 

# Command to run the executable
CMD ["/app/main"]