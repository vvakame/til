FROM golang:1.12

LABEL "name"="md-to-blogpost"
LABEL "maintainer"="vvakame <vvakame@gmail.com>"
LABEL "version"="1.0.0"

LABEL "com.github.actions.name"="Markdown to Blog post"
LABEL "com.github.actions.description"="Convert Markdown to Blog post"
LABEL "com.github.actions.icon"="send"
LABEL "com.github.actions.color"="yellow"

WORKDIR /project/
ENV GO111MODULE=on

COPY ./ ./

RUN go build -o /md2blog main.go

ENTRYPOINT ["/md2blog"]
