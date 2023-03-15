FROM node:19-alpine AS build

WORKDIR /usr/src/app

COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM busybox:1.35

RUN adduser -D static
USER static
WORKDIR /home/static
COPY --from=build /usr/src/app/dist .

EXPOSE 80
CMD ["busybox", "httpd", "-f", "-v", "-p", "80"]