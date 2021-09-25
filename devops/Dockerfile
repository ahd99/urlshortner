FROM golang:1.16-alpine
WORKDIR /app
# suppose context is project root
COPY go.mod ./   
#COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /shortner github.com/ahd99/urlshortner/cmd/server
CMD [ "/shortner" ]
