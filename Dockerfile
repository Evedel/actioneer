FROM golang:1.21.1
COPY . /app
WORKDIR /app
RUN go build -o on-call-automaton . 
ENTRYPOINT [ "/app/on-call-automaton" ]
