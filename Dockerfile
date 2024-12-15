# Etapa de build
FROM golang:1.23 AS build
WORKDIR /app
COPY . .
# Configuração para build estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o buscacep

# Etapa de runtime
FROM alpine:latest
WORKDIR /app

# Copiar binário gerado da etapa anterior
COPY --from=build /app/buscacep .

# Definir variável de ambiente (a chave será passada no momento do deploy)
ENV WEATHER_API_KEY="5ef06d7fbd5743b69ed150449241512"

# Expor porta usada pelo aplicativo
EXPOSE 8080

# Comando de inicialização
ENTRYPOINT ["./buscacep"]
