FROM oven/bun:1 AS builder
WORKDIR /usr/src/app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile
COPY . .
ARG node_env
ENV NODE_ENV=$node_env
RUN bun build.js

FROM nginx:1.29.2-alpine
COPY --from=builder /usr/src/app/dist /usr/share/nginx/html
