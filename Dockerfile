FROM golang:alpine AS builder


#RUN apk update \
# && apk add --no-cache \
#    git \
#    openssh \
#    'su-exec>=0.2'

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on
#ENV GOPRIVATE=github.com/offerm/lnagent
#    CGO_ENABLED=0 \
#    GOOS=linux \
#    GOARCH=amd64

#RUN git config --global url.git@github.com:.insteadOf https://github.com/
#RUN mkdir -p /root/.ssh && \
#    chmod 0700 /root/.ssh && \
#    ssh-keyscan -t rsa github.com > /root/.ssh/known_hosts && \
#    echo "${SSH_KEY}" > /root/.ssh/id_rsa && \
#    chmod 600 /root/.ssh/id_rsa

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o lnagent cmd/lnagent/main.go

## Move to /dist directory as the place for resulting binary folder
#WORKDIR /dist

FROM golang:alpine 

RUN apk update \
 && apk add --no-cache \
    tini 
#    openssh \
#    'su-exec>=0.2'

COPY --from=builder /build/lnagent /


# Command to run when starting the container
#CMD ["lnagent", "run", "--lnch", "34.73.171.247"]

COPY docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["/sbin/tini", "--", "/bin/sh", "/usr/local/bin/docker-entrypoint.sh"]
