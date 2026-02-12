# ⚓ Afuradanime Backend

Um *Backend* desenvolvido em [Go](https://go.dev/), responsável por disponibilizar a **API** principal do Afuradanime.

 Este módulo 100% *Open Source* (Código aberto) e desenvolvido com [Clean Architecture / Onion Architecture](https://blog.ploeh.dk/2013/12/03/layers-onions-ports-adapters-its-all-the-same/) em mente, fornece**rá** os *endpoints* necessários para a gestão de utilizadores, *animelists*, *threads* e *posts*, tendo também a conexão direta aos componentes de *Anime* e *Manga*.

**Bibliotecas externas utilizadas:**
* **[Chi](github.com/go-chi/chi/v5)** - Router leve e idiomático para Go
* **[Godotenv](github.com/joho/godotenv)** - Carregamento de variáveis de ambiente por ficheiros `.env`

## Como correr o backend
`go run .\cmd\api\main.go`