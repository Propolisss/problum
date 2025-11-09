FROM node:25.1.0-alpine3.21 AS builder

WORKDIR /app

COPY package*.json .

RUN npm i --frozen-lockfile

COPY . .

RUN npm run build

FROM nginx:1.29.3-alpine3.22

COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
