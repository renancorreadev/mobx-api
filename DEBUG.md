# Debug do Container Docker

Este projeto inclui configurações para debug do aplicativo Go dentro de containers Docker usando o Delve debugger.

## Configurações Disponíveis

### 1. Debug Local
- **Nome**: `Debug Go App (Local)`
- **Descrição**: Executa o aplicativo localmente com debug
- **Uso**: Para desenvolvimento local sem Docker

### 2. Debug Container Docker
- **Nome**: `Debug Go App (Docker Container)`
- **Descrição**: Conecta ao debugger em um container Docker já em execução
- **Uso**: Quando você já tem o container rodando com debug habilitado

### 3. Debug Docker Compose
- **Nome**: `Debug Go App (Docker Compose)`
- **Descrição**: Inicia automaticamente o ambiente completo com debug
- **Uso**: Para debug completo com banco de dados e Redis

## Como Usar

### Opção 1: Debug com Docker Compose (Recomendado)

1. **Iniciar o debug**:
   - Pressione `F5` ou vá em `Run and Debug`
   - Selecione `Debug Go App (Docker Compose)`
   - O VS Code irá automaticamente:
     - Iniciar os containers (PostgreSQL, Redis, API)
     - Conectar ao debugger
     - Parar na primeira linha do `main()`

2. **Configurar breakpoints**:
   - Clique na margem esquerda do editor para adicionar breakpoints
   - Os breakpoints funcionarão normalmente

3. **Parar o debug**:
   - Pressione `Shift+F5` para parar o debug
   - Os containers serão automaticamente parados

### Opção 2: Debug Manual

1. **Construir a imagem de debug**:
   ```bash
   docker build -f Dockerfile.debug -t vfinance-api-debug .
   ```

2. **Executar o container com debug**:
   ```bash
   docker run --rm -d \
     --name vfinance-api-debug \
     -p 3001:3000 \
     -p 2345:2345 \
     --security-opt seccomp:unconfined \
     --cap-add SYS_PTRACE \
     vfinance-api-debug
   ```

3. **Conectar o debugger**:
   - Pressione `F5`
   - Selecione `Debug Go App (Docker Container)`

### Opção 3: Debug Local

1. **Certifique-se de que PostgreSQL e Redis estão rodando**:
   ```bash
   docker-compose up -d postgres redis
   ```

2. **Iniciar debug local**:
   - Pressione `F5`
   - Selecione `Debug Go App (Local)`

## Arquivos de Configuração

- `.vscode/launch.json`: Configurações do debugger
- `.vscode/tasks.json`: Tarefas automatizadas
- `Dockerfile.debug`: Dockerfile específico para debug
- `docker-compose.debug.yml`: Docker Compose para debug

## Variáveis de Ambiente

O debugger usa as mesmas variáveis de ambiente do ambiente de produção:

- `API_PORT`: Porta da API (3000)
- `DATABASE_URL`: URL do PostgreSQL
- `REDIS_URL`: URL do Redis
- `ETHEREUM_RPC`: Endpoint RPC do Ethereum
- `CONTRACT_ADDRESS`: Endereço do contrato
- `PRIVATE_KEY`: Chave privada
- `JWT_SECRET`: Chave secreta JWT
- `RATE_LIMIT_WINDOW`: Janela de rate limiting
- `RATE_LIMIT_MAX`: Máximo de requisições

## Troubleshooting

### Erro de Conexão com Debugger
- Verifique se a porta 2345 está exposta
- Certifique-se de que o container está rodando
- Verifique os logs do container: `docker logs vfinance-api-debug`

### Erro de Permissão
- O container precisa das permissões `SYS_PTRACE` e `seccomp:unconfined`
- Verifique se o Dockerfile.debug está sendo usado

### Breakpoints Não Funcionam
- Certifique-se de que o código foi compilado com símbolos de debug (`-gcflags="all=-N -l"`)
- Verifique se o `remotePath` está correto no `launch.json`

## Comandos Úteis

```bash
# Ver logs do container de debug
docker logs vfinance-api-debug

# Entrar no container
docker exec -it vfinance-api-debug sh

# Parar todos os containers
docker-compose -f docker-compose.debug.yml down

# Reconstruir imagem de debug
docker build -f Dockerfile.debug -t vfinance-api-debug .
``` 