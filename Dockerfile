# -----------------------------------------------------------------------------
# Stage 1: 构建前端 views (Vue)
# -----------------------------------------------------------------------------
FROM node:20-slim AS frontend
# 临时跳过 SSL 校验（仅用于排查证书问题，生产请移除）
ENV NODE_TLS_REJECT_UNAUTHORIZED=0
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
RUN corepack enable pnpm
WORKDIR /views
COPY views/package.json views/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY views/ .
RUN pnpm build

# -----------------------------------------------------------------------------
# Stage 2: 构建后端 coredns-admin (Go)
# -----------------------------------------------------------------------------
FROM golang:1.22-alpine AS backend
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /views/dist ./dist
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o coredns-admin .

# -----------------------------------------------------------------------------
# Stage 3: 运行镜像（etcd + coredns-admin + 静态前端）
# -----------------------------------------------------------------------------
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend /app/coredns-admin .
COPY --from=backend /app/dist ./dist
COPY docker/config.yaml ./config.yaml
EXPOSE 8088
ENTRYPOINT ["./coredns-admin"]
CMD ["-C", "/app/config.yaml"]
