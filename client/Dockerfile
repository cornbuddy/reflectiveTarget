FROM node:24-alpine AS builder
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.29.2-alpine
COPY --from=builder /favicon.ico /usr/share/nginx/html
COPY --from=builder /img /usr/share/nginx/html/img
COPY --from=builder /dist /usr/share/nginx/html
