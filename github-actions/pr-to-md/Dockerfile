FROM golang:1.12

LABEL "name"="pr-to-md"
LABEL "maintainer"="vvakame <vvakame@gmail.com>"
LABEL "version"="1.0.0"

LABEL "com.github.actions.name"="Pull Request to Markdown"
LABEL "com.github.actions.description"="Convert PR to Markdown"
LABEL "com.github.actions.icon"="user-plus"
LABEL "com.github.actions.color"="green"

WORKDIR /project/
ENV GO111MODULE=on

COPY go.* ./
COPY *.go ./

RUN go build -o /pr2md main.go

ENTRYPOINT ["/pr2md"]
