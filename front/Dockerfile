FROM node:19-alpine

WORKDIR /usr/src/app

COPY package*.json ./
RUN npm install

COPY . .
RUN rm ./src/.env.development && rm ./src/.env.example && mv ./src/.env.production ./src/.env

EXPOSE 5173
CMD [ "npm", "run", "dev", "--", "--host" ]