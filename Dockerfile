# --- Estágio 1: Build ---
# Usamos uma imagem Go completa para compilar a aplicação
FROM golang:1.25-alpine AS builder

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Um argumento que vamos passar no momento do build para dizer qual serviço construir.
# Ex: 'services/pedidos'
ARG SERVICE_PATH
ARG MAIN_FILE_PATH

# Copia os arquivos de módulo primeiro para aproveitar o cache do Docker
COPY go.work go.work.sum ./
COPY ${SERVICE_PATH}/go.mod ${SERVICE_PATH}/go.mod
COPY ${SERVICE_PATH}/go.sum ${SERVICE_PATH}/go.sum
COPY pkg/go.mod pkg/go.mod
COPY pkg/go.sum pkg/go.sum

# Baixa as dependências
RUN go work sync
RUN go mod download

# Copia todo o resto do código fonte para o contêiner
COPY . .

# Compila a aplicação. O binário será estático, sem dependências externas.
# O -o define o nome do arquivo de saída como 'server'
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./${MAIN_FILE_PATH}


# --- Estágio 2: Final ---
# Usamos uma imagem "scratch" (vazia) ou "alpine" para a imagem final. Alpine é bom para debug.
FROM alpine:latest

WORKDIR /

# Copia o binário compilado do estágio 'builder'
COPY --from=builder /app/server /server

# Expõe a porta que a nossa aplicação usa (ex: 8080)
EXPOSE 8080

# Comando para iniciar a aplicação quando o contêiner rodar
CMD ["/server"]