FROM golang:latest as builder
ENV GOPATH /var/lib/jenkins/workspace-go
ENV APP_ROOT /var/lib/jenkins/workspace-go/src/go-game
WORKDIR ${APP_ROOT}
COPY ./ ${APP_ROOT}
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest as final
ENV APP_ROOT /var/lib/jenkins/workspace-go/src/go-game
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app/
RUN mkdir conf && touch ./conf/app.conf
RUN mkdir -p ./data/
COPY --from=builder ${APP_ROOT}/main .
COPY --from=builder ${APP_ROOT}/conf/app.conf ./conf/app.conf
COPY --from=builder ${APP_ROOT}/data/sys-peotry-set.json ./data/sys-peotry-set.json
COPY --from=builder ${APP_ROOT}/data/sys-peotry.json ./data/sys-peotry.json
EXPOSE 8089
ENV SGHENENV prod
ENTRYPOINT ["/app/main"]