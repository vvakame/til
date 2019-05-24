FROM golang:1.12

LABEL "name"="pr-to-md"
LABEL "maintainer"="vvakame <vvakame@gmail.com>"
LABEL "version"="1.0.0"

LABEL "com.github.actions.name"="Pull Request to Markdown"
LABEL "com.github.actions.description"="Convert PR to Markdown"
LABEL "com.github.actions.icon"="edit-2"
LABEL "com.github.actions.color"="orange"

WORKDIR /project/
ENV GO111MODULE=on

COPY ./ ./

RUN go install github.com/rakyll/statik && \
    go generate . && \
    go build -o /pr2md main.go

ENTRYPOINT ["/pr2md"]
