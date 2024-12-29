FROM golang:1.23-alpine

# Definindo o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiando os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Baixando as dependências
RUN go mod download

# Copiando o código-fonte
COPY . .

# Construindo o executável
RUN go build -o etl

# Comando para rodar o ETL
CMD ["./etl"]