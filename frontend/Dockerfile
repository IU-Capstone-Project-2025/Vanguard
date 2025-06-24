FROM node:18-alpine AS builder
WORKDIR /app
COPY ./app .
RUN npm install && npm run build && npm cache clean --force

FROM nginx:1.28-alpine-slim
COPY --from=builder /app/build /usr/share/nginx/html
EXPOSE 80