# --- Estágio 1: Build ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Argumentos para dizer qual serviço construir
ARG SERVICE_PATH
ARG MAIN_FILE_PATH

# --- SEÇÃO ALTERADA ---
# Copia o go.work e TODOS os arquivos de módulo primeiro.
# Isso garante que o 'go work sync' tenha a visão completa do projeto.
COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/pedidos/go.mod services/pedidos/go.sum ./services/pedidos/
COPY services/clientes/go.mod services/clientes/go.sum ./services/clientes/

# Agora o 'go work sync' funcionará, pois ele encontra todos os módulos.
RUN go work sync
RUN go mod download

# Copia todo o resto do código fonte para o contêiner
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./${MAIN_FILE_PATH}


# --- Estágio 2: Final ---
FROM alpine:latest

WORKDIR /

# Copia o binário compilado do estágio 'builder'
COPY --from=builder /app/server /server

# Expõe a porta que a nossa aplicação usa
EXPOSE 8080

# Comando para iniciar a aplicação quando o contêiner rodar
CMD ["/server"]