FROM node:20-alpine
WORKDIR /app
COPY ./app/package.json ./
RUN npm install
CMD ["npm", "start"]