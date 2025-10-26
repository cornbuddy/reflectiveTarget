FROM node:24-alpine
COPY package*.json ./
RUN npm install
COPY . .
ENTRYPOINT ["npm"]
CMD ["start"]
