<h1 align="center">
  ⚓ Afuradanime Backend
</h1>

Um *Backend* desenvolvido em [Go](https://go.dev/), responsável por disponibilizar a **API** principal do Afuradanime.

 Este módulo 100% *Open Source* (Código aberto) e desenvolvido com [Clean Architecture / Onion Architecture](https://blog.ploeh.dk/2013/12/03/layers-onions-ports-adapters-its-all-the-same/) em mente, fornece os *endpoints* necessários para a gestão de utilizadores, *animelists*, *threads* e *posts*, tendo também a conexão direta aos componentes de *Anime* e *Manga*.

**Bibliotecas externas utilizadas:**
1. **[Anime Facts Core](https://github.com/afuradanime/anime-facts-core)** - Acesso rápido e local a informação sobre animes
2. **[Chi](https://github.com/go-chi/chi)** - Router leve e idiomático para Go
3. **[Godotenv](https://github.com/joho/godotenv)** - Carregamento de variáveis de ambiente por ficheiros `.env`
4. **[Cors for Chi](https://github.com/go-chi/cors)** - Definição de regras de cabeçalho [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS)
5. **[JWT for Go](https://github.com/golang-jwt/jwt)** - Implementação de [JWT](https://www.jwt.io/introduction#what-is-json-web-token)'s para auth/authz
6. **[Mongo Driver](https://go.mongodb.org/mongo-driver)** - Driver para o SGBD [MongoDB](https://www.mongodb.com/)
7. **[OAuth2 for Go](https://golang.org/x/oauth2)** - Suporte para uso do protocolo de autenticação [OAuth2](https://oauth.net/2/)
8. **[Rate](https://pkg.go.dev/golang.org/x/time/rate)** - Suporte de rate limiting

## Documentação da API

<div style="display: flex; justify-content: space-around">
    <div>
        <h3> <a href="./_docs/api.http"><i>Endpoints</i> documentados</a> </h3>
        <h3> <a href="./_docs/domain.md">Modelo de domínio</a> </h3>
    </div>
<div>
    <h3> <a>Modelo da base de dados</a></h3>
    <h3> <a href="./_docs/architecture.md">Vistas arquitecturais</a> </h3>
</div>
</div>

## Dependências do backend
Para que o backend do AfuradaAnime funcione como deve, é necessário ter as suas dependências preparadas:

1. **Drivers**:
Na pasta `/drivers` devem ser inseridos executaveis/bibliotecas que alimentam o backend do Afuradanime.
    1. **anime_facts.dll**/**anime_facts.so**: A biblioteca partilhada de acesso à base de dados de animes. Instruções de compilação podem ser encontradas no repositório do [AFC](https://github.com/afuradanime/anime-facts-core); alternativamente usar os recursos na sua página de releases.

2. **A base de dados de animes**:
A base de dados de animes tem instruções de formação no projecto do [AFC](https://github.com/afuradanime/anime-facts-core)

3. **Configurações `.env`**:
O backend procura um ficheiro `.env` com certas configurações do seu sistema, um exemplo pode ser encontrado no ficheiro [example.env](./example.env)

## Como correr o backend

`go run .\cmd\api\main.go`

## Como correr os testes

TODO

## Como fazer deploy

TODO