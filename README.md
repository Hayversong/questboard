# QuestBoard

Kanban gamificado para acompanhar o desenvolvimento de jogos, construido do zero com Go.

QuestBoard e um sistema web estilo Kanban onde cada tarefa e uma Quest, o progresso vira XP e projetos ganham Level e Rank conforme as tarefas sao concluidas.

O projeto existe principalmente para aprender engenharia de software construindo algo real: Go, net/http, html/template, CSS puro e JavaScript vanilla.

## Funcionalidades

### Projetos
- Criar projeto
- Renomear projeto
- Excluir projeto
- Dashboard com metricas globais de projetos, quests e XP total

### Cards / Quests
- Criar card com titulo, descricao, raridade e prazo
- Editar card inline
- Excluir card
- Mover entre colunas via drag and drop
- Reordenar cards dentro da coluna
- Persistir status e ordem ao recarregar a pagina

### Gamificacao
- XP por raridade de card
- Level calculado a partir do XP acumulado
- Rank baseado no level
- Conquistas desbloqueadas por progressao
- Log de atividades recentes do projeto

### Operacao
- Persistencia em JSON para uso local simples
- Persistencia em SQLite para deploy
- Migracao de JSON para SQLite
- Healthcheck em `GET /healthz`
- Porta configuravel por `PORT`
- Imagem Docker com volume persistente em `/data`

## Stack

**Backend**
- Go
- `net/http`
- `html/template`

**Frontend**
- HTML semantico
- CSS puro
- JavaScript vanilla

**Persistencia**
- JSON por padrao em desenvolvimento local
- SQLite recomendado para producao

## Arquitetura

```text
HTTP Request
     |
  Handler     -> le request, chama service, renderiza template
     |
  Service     -> regras de negocio, validacoes, erros sentinela
     |
  Storage     -> le e escreve JSON ou SQLite
     |
  Model       -> structs e metodos de dominio
```

Cada camada tem responsabilidade unica. Handlers nao acessam storage diretamente.

## Estrutura principal

```text
questboard/
  cmd/
    server/      # servidor web
    migrate/     # migracao JSON -> SQLite
  internal/
    handler/     # handlers HTTP
    model/       # entidades e regras de dominio
    service/     # regras de negocio
    storage/     # persistencia JSON e SQLite
  web/
    templates/   # HTML server-side
    static/      # CSS e JS
  data/
    projects.example.json
  docs/
```

## Rotas

| Metodo | Rota               | Descricao |
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
| POST   | `/cards/move`      | Avancar status do card |
| POST   | `/cards/status`    | Atualizar status diretamente |
| POST   | `/cards/reorder`   | Salvar ordem e status |

## Variaveis de ambiente

| Variavel | Padrao | Descricao |
|----------|--------|-----------|
| `PORT` | `8080` | Porta HTTP usada pelo servidor |
| `QUESTBOARD_STORAGE` | vazio | Use `sqlite` para ativar SQLite explicitamente |
| `QUESTBOARD_DATA_FILE` | `data/projects.json` | Caminho do arquivo JSON |
| `QUESTBOARD_DATA_DIR` | `data` | Diretorio base para dados quando o arquivo nao e informado |
| `QUESTBOARD_DB_FILE` | `data/questboard.db` | Caminho do banco SQLite |

Observacao: se `QUESTBOARD_DB_FILE` estiver definido, o projeto usa SQLite mesmo sem `QUESTBOARD_STORAGE=sqlite`.

## Como executar localmente com JSON

**Requisitos:** Go 1.25+ ou uma versao compativel com o `go.mod`.

```bash
go run ./cmd/server
```

Acesse:

```text
http://localhost:8080
```

Nesse modo, os dados sao salvos em `data/projects.json`. O arquivo e criado automaticamente na primeira gravacao.

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

SQLite e o caminho recomendado para deploy porque usa um unico arquivo de banco e funciona bem com volume persistente.

## Migracao de JSON para SQLite

Se voce ja tem dados em JSON, gere um banco SQLite com:

```bash
go run ./cmd/migrate -from data/projects.json -to data/questboard.db
```

Se o banco de destino ja existir e voce quiser sobrescrever:

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

A imagem Docker define por padrao:

```bash
PORT=8080
QUESTBOARD_STORAGE=sqlite
QUESTBOARD_DB_FILE=/data/questboard.db
```

Importante: em producao, use sempre um volume persistente para `/data`. Sem volume persistente, os dados podem ser perdidos em rebuild, redeploy ou recriacao do container.

## Deploy

Para o primeiro deploy, use uma plataforma que suporte container Docker e volume persistente.

Configuracao recomendada:

```bash
PORT=8080
QUESTBOARD_STORAGE=sqlite
QUESTBOARD_DB_FILE=/data/questboard.db
```

Se a plataforma injeta `PORT` automaticamente, nao fixe a porta manualmente; deixe a variavel da plataforma controlar isso.

## Verificacao

Antes de publicar uma versao, rode:

```bash
go test ./...
go build ./cmd/server
go build ./cmd/migrate
docker build -t questboard .
```

Em Windows acessando o projeto por `\\wsl.localhost`, o Go pode falhar com erro de lock no `go.mod`. Nesse caso, rode os comandos dentro do WSL com Go instalado ou copie o projeto para uma pasta temporaria local do Windows para validar.

## Seguranca

Esta primeira versao nao possui autenticacao e os formularios POST nao possuem protecao CSRF.

Para uso publico, coloque o app atras de uma camada externa de protecao, como acesso privado da plataforma, VPN, basic auth no proxy ou solucao equivalente. Qualquer pessoa com acesso ao app pode criar, editar e excluir projetos e cards.

## Roadmap

- Filtro e busca de cards
- Tags nos cards
- Responsaveis por card
- API REST
- Autenticacao
- Melhorias de operacao e observabilidade

## Agentes de IA

Para continuar o preparo de deploy ou revisar os criterios de aceite, leia `docs/DEPLOY_AGENT_BRIEF.md`.
