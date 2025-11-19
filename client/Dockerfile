FROM node:24-alpine AS builder
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.29.2-alpine
COPY --from=builder /dist /usr/share/nginx/html
