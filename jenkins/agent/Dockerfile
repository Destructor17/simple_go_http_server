FROM jenkins/ssh-agent:debian-jdk11

# WORKDIR /usr/src/app

# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

# COPY . .
RUN apt update
RUN apt install -y golang
