FROM golang:latest
# ENV GOPATH /home/luoguibin/CompanyCode/go
# ENV APP_ROOT /var/lib/jenkins/workspace/go-game
# WORKDIR ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
ENV APP_ROOT ./
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app/
RUN mkdir conf && touch ./conf/app.conf
RUN mkdir -p ./data/
COPY --from=0 ${APP_ROOT}/main .
COPY --from=0 ${APP_ROOT}/conf/app.conf ./conf/app.conf
COPY --from=0 ${APP_ROOT}/data/sys-peotry-set.json ./data/sys-peotry-set.json
COPY --from=0 ${APP_ROOT}/data/sys-peotry.json ./data/sys-peotry.json
EXPOSE 8089
ENV SGHENENV prod
ENTRYPOINT ["/app/main"]