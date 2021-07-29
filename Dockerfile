FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod ./
#COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /shortner github.com/ahd99/urlshortner/cmd/server
EXPOSE 8081
CMD [ "/shortner" ]

# build command:
# docker build -t shortner:latest .  

# run command:
# docker run -p 8081:8081 shortner:latest