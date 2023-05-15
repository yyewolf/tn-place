FROM node:19-alpine AS build-front
LABEL maintainer="yyewolf"

WORKDIR /usr/src/app

COPY front/package*.json ./
RUN npm install
COPY front .
RUN npm run build

FROM golang:1.20-alpine AS build-back
LABEL maintainer="yyewolf"

WORKDIR /app
COPY back .

WORKDIR /app/cmd/place
RUN go mod download
RUN go build -o /app/main

FROM alpine:3.18
LABEL maintainer="yyewolf"

WORKDIR /app
COPY --from=build-back /app/main /app/main
COPY --from=build-front /usr/src/app/dist /app/dist

RUN chmod +x /app/main

EXPOSE 80
CMD ["/app/main"]