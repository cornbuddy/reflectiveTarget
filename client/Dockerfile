FROM node:24-alpine AS builder
COPY package*.json ./
RUN npm install
COPY . .
ARG node_env
ENV NODE_ENV=$node_env
RUN npm run build

FROM nginx:1.29.2-alpine
COPY --from=builder /dist /usr/share/nginx/html
