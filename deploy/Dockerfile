FROM node:26-alpine AS web
RUN npm install --global pnpm@12.4.2
WORKDIR /src/web
COPY web/package.json web/pnpm-lock.yaml web/.npmrc ./
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm build

FROM golang:1.27.2-alpine AS server
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
COPY --from=web /src/web/build/ internal/web/dist/
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/poker ./cmd/poker

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=server /out/poker /poker
USER nonroot:nonroot
EXPOSE 8080 2222
ENTRYPOINT ["/poker"]
