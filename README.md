# QuestBoard

[![CI](https://github.com/Hayversong/questboard/actions/workflows/ci.yml/badge.svg)](https://github.com/Hayversong/questboard/actions/workflows/ci.yml)

Kanban gamificado para acompanhar o desenvolvimento de jogos, construído do zero com Go.

QuestBoard é um sistema web estilo Kanban onde cada tarefa é uma Quest, o progresso vira XP e projetos ganham Level e Rank conforme as tarefas são concluidas.

O projeto existe principalmente para aprender engenharia de software construindo algo real: Go, net/http, html/template, CSS puro e JavaScript vanilla.

## Funcionalidades

### Projetos
- Criar projeto
- Renomear projeto
- Excluir projeto
- Dashboard com métricas globais de projetos, quests e XP total

### Cards / Quests
- Criar card com título, descrição, raridade e prazo
- Editar card inline
- Excluir card
- Mover entre colunas via drag and drop
- Reordenar cards dentro da coluna
- Persistir status e ordem ao recarregar a página

### Gamificação
- XP por raridade de card
- Level calculado a partir do XP acumulado
- Rank baseado no level
- Conquistas desbloqueadas por progressão
- Log de atividades recentes do projeto

### Operação
- Persistência em JSON para uso local simples
- Persistência em SQLite para deploy
- Migração de JSON para SQLite
- Healthcheck em `GET /healthz`
- Porta configurável por `PORT`
- Imagem Docker com volume persistente em `/data`

## Stack

**Backend**
- Go
- `net/http`
- `html/template`

**Frontend**
- HTML semântico
- CSS puro
- JavaScript vanilla

**Persistência**
- JSON por padrão em desenvolvimento local
- SQLite recomendado para produção

## Arquitetura

```text
HTTP Request
     |
  Handler     -> lê request, chama service, renderiza template
     |
  Service     -> regras de negócio, validações, erros sentinela
     |
  Storage     -> lê e escreve JSON ou SQLite
     |
  Model       -> structs e métodos de domínio
```

Cada camada tem responsabilidade única. Handlers não acessam storage diretamente.

## Estrutura principal

```text
questboard/
  cmd/
    server/      # servidor web
    migrate/     # migração JSON -> SQLite
  internal/
    handler/     # handlers HTTP
    model/       # entidades e regras de domínio
    service/     # regras de negócio
    storage/     # persistencia JSON e SQLite
  web/
    templates/   # HTML server-side
    static/      # CSS e JS
  data/
    projects.example.json
  docs/
```

## Rotas

| Método | Rota               | Descrição |
|--------|--------------------|-----------|
| GET    | `/`                | Home com lista de projetos |
| GET    | `/project?id=`     | Board Kanban do projeto |
| GET    | `/healthz`         | Healthcheck |
| POST   | `/projects`        | Criar projeto |
| POST   | `/projects/rename` | Renomear projeto |
| POST   | `/projects/delete` | Excluir projeto |
| POST   | `/cards`           | Criar card |
| POST   | `/cards/update`    | Editar card |
| POST   | `/cards/delete`    | Excluir card |
| POST   | `/cards/move`      | Avançar status do card |
| POST   | `/cards/status`    | Atualizar status diretamente |
| POST   | `/cards/reorder`   | Salvar ordem e status |

## Variáveis de ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `PORT` | `8080` | Porta HTTP usada pelo servidor |
| `QUESTBOARD_STORAGE` | vazio | Use `sqlite` para ativar SQLite explicitamente |
| `QUESTBOARD_DATA_FILE` | `data/projects.json` | Caminho do arquivo JSON |
| `QUESTBOARD_DATA_DIR` | `data` | Diretório base para dados quando o arquivo não é informado |
| `QUESTBOARD_DB_FILE` | `data/questboard.db` | Caminho do banco SQLite |

Observação: se `QUESTBOARD_DB_FILE` estiver definido, o projeto usa SQLite mesmo sem `QUESTBOARD_STORAGE=sqlite`.

## Como executar localmente com JSON

**Requisitos:** Go 1.25+ ou uma versão compatível com o `go.mod`.

```bash
go run ./cmd/server
```

Acesse:

```text
http://localhost:8080
```

Nesse modo, os dados são salvos em `data/projects.json`. O arquivo é criado automaticamente na primeira gravacao.

## Como executar localmente com SQLite

```bash
QUESTBOARD_STORAGE=sqlite QUESTBOARD_DB_FILE=data/questboard.db go run ./cmd/server
```

No PowerShell:

```powershell
$env:QUESTBOARD_STORAGE = "sqlite"
$env:QUESTBOARD_DB_FILE = "data/questboard.db"
go run ./cmd/server
```

SQLite é o caminho recomendado para deploy porque usa um único arquivo de banco e funciona bem com volume persistente.

## Migração de JSON para SQLite

Se você já tem dados em JSON, gere um banco SQLite com:

```bash
go run ./cmd/migrate -from data/projects.json -to data/questboard.db
```

Se o banco de destino já existir e você quiser sobrescrever:

```bash
go run ./cmd/migrate -from data/projects.json -to data/questboard.db -force
```

## Docker

Build da imagem:

```bash
docker build -t questboard .
```

Executar com volume persistente:

```bash
docker run --rm -p 8080:8080 -v questboard-data:/data questboard
```

A imagem Docker define por padrão:

```bash
PORT=8080
QUESTBOARD_STORAGE=sqlite
QUESTBOARD_DB_FILE=/data/questboard.db
```

Importante: em produção, use sempre um volume persistente para `/data`. Sem volume persistente, os dados podem ser perdidos em rebuild, redeploy ou recriacao do container.

## Deploy

Para o primeiro deploy, use uma plataforma que suporte container Docker e volume persistente.

Configuração recomendada:

```bash
PORT=8080
QUESTBOARD_STORAGE=sqlite
QUESTBOARD_DB_FILE=/data/questboard.db
```

Se a plataforma injeta `PORT` automaticamente, não fixe a porta manualmente; deixe a variavel da plataforma controlar isso.

## Verificação

Antes de publicar uma versão, rode:

```bash
go test ./...
go build ./cmd/server
go build ./cmd/migrate
docker build -t questboard .
```

Em Windows acessando o projeto por `\\wsl.localhost`, o Go pode falhar com erro de lock no `go.mod`. Nesse caso, rode os comandos dentro do WSL com Go instalado ou copie o projeto para uma pasta temporária local do Windows para validar.

## Segurança

Esta primeira versão não possui autenticação e os formulários POST não possuem proteção CSRF.

Para uso publico, coloque o app atras de uma camada externa de proteção, como acesso privado da plataforma, VPN, basic auth no proxy ou solução equivalente. Qualquer pessoa com acesso ao app pode criar, editar e excluir projetos e cards.

## Roadmap

- Filtro e busca de cards
- Tags nos cards
- Responsáveis por card
- API REST
- Melhorias de operação e observabilidade
