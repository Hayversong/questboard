# QuestBoard

QuestBoard é um sistema web estilo **Kanban** para acompanhar o desenvolvimento de projetos de jogos.

O objetivo principal deste projeto é estudar **desenvolvimento web com Go**, construindo uma aplicação real sem frameworks.

---

## Motivação

Depois de construir uma To-Do App simples, surgiu a ideia de criar algo mais próximo de um produto real.

QuestBoard nasceu para resolver um problema pessoal:

> acompanhar desenvolvimento de jogos organizando features em um board simples.

Exemplo:

```text
Dungeon Survivor

Backlog
- Loja
- Boss

Doing
- Sistema XP

Done
- Movimento
- Inimigos
```

---

# Tecnologias

Backend:
- Go
- net/http
- html/template

Frontend:
- HTML
- CSS

Persistência:
- JSON

Ferramentas:
- Git
- GitHub
- WSL

---

# Estrutura do Projeto

```text
questboard/

├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── handler/
│   ├── model/
│   ├── service/
│   └── storage/
│
├── web/
│   ├── templates/
│   └── static/
│
├── data/
│   └── projects.json
│
├── go.mod
└── README.md
```

---

# Arquitetura

O projeto segue separação simples de responsabilidades.

```text
HTTP
↓
Handler
↓
Service
↓
Storage
↓
JSON
```

### Handler
Responsável por:
- receber requisições
- carregar dados
- renderizar templates

---

### Service
Responsável por:
- regras de negócio
- validações
- operações da aplicação

---

### Storage
Responsável por:
- leitura de JSON
- escrita de JSON

---

### Model
Responsável por:
- representar entidades do sistema

---

# Modelo de Dados

## Project

```go
type Project struct {
	ID string
	Name string
	Cards []Card
}
```

## Card

```go
type Card struct {
	ID string
	Title string
	Description string
	Status string
}
```

Status disponíveis:

```text
backlog
doing
done
```

---

# Funcionalidades

## Implementadas

- [x] Estrutura inicial
- [x] Servidor HTTP
- [x] Templates HTML
- [x] Persistência em JSON
- [x] Visualização de projetos
- [x] Visualização do board
- [x] Organização por colunas

---

## Em desenvolvimento

- [ ] Criar cards
- [ ] Excluir cards
- [ ] Mover entre colunas
- [ ] Melhorar UI
- [ ] Validação de dados

---

## Futuro

- [ ] SQLite
- [ ] API REST
- [ ] Drag and Drop
- [ ] Dashboard
- [ ] Tags
- [ ] Docker
- [ ] Deploy

---

# Como executar

Clone:

```bash
git clone git@github.com:SEU_USUARIO/questboard.git
```

Entrar:

```bash
cd questboard
```

Executar:

```bash
go run ./cmd/server
```

Abrir:

```text
http://localhost:8080
```

---

# Persistência

Os dados são armazenados em:

```text
data/projects.json
```

Exemplo:

```json
[
  {
    "ID": "1",
    "Name": "Dungeon Survivor",
    "Cards": [
      {
        "ID": "1",
        "Title": "Sistema XP",
        "Status": "doing"
      }
    ]
  }
]
```

---

# Objetivo de aprendizado

Este projeto existe para praticar:

- Go Web
- Arquitetura em camadas
- Persistência
- Templates
- Organização de código
- Git
- Desenvolvimento incremental

---

Feito por Hayverson