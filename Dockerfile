FROM node:22-alpine AS frontend-builder

WORKDIR /frontend

COPY RPO-Frontend/package.json RPO-Frontend/package-lock.json ./
RUN npm ci

COPY RPO-Frontend/index.html RPO-Frontend/tsconfig.json RPO-Frontend/vite.config.ts ./
COPY RPO-Frontend/src ./src
RUN npm run build

FROM golang:1.24-alpine AS backend-builder

RUN apk add --no-cache gcc libc-dev

WORKDIR /backend

COPY RPO-Backend/go.mod RPO-Backend/go.sum ./
RUN go mod download

COPY RPO-Backend ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /out ./

FROM alpine:3.20

RUN apk add --no-cache nginx openssl

WORKDIR /app

COPY RPO-Backend/data /app/data
COPY --from=frontend-builder /frontend/dist /app/web/dist
COPY --from=backend-builder /out /api
COPY RPO-Backend/deploy/nginx.conf /etc/nginx/nginx.conf
COPY RPO-Backend/deploy/certs /deploy/certs
COPY RPO-Backend/deploy/entrypoint.sh /deploy/entrypoint.sh

RUN chmod +x /deploy/entrypoint.sh

EXPOSE 8888

CMD ["/deploy/entrypoint.sh"]