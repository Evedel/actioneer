FROM golang:1.21.1
COPY ./app /app
WORKDIR /app
RUN go build -o automaton cmd/main.go
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    rm kubectl
ENTRYPOINT [ "/app/automaton" ]
