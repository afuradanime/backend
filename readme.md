# Afuradanime Backend

 Um *backend* desenvolvido em Go que disponibiliza o *API* principal para o Afuradanime.

 Este módulo 100% *Open Source* (Código aberto) e desenvolvido com [arquitetura onion](https://blog.ploeh.dk/2013/12/03/layers-onions-ports-adapters-its-all-the-same/) em mente, fornece**rá** os *APIs* necessários para a manutenção de ulizadores, *animelists*, *threads*, *posts*, etc... assim como a conexão direta ao módulo de *Anime* e *Manga*.

*Packages* externos utilizados:
* **[Chi](github.com/go-chi/chi/v5)** - setting up routing
* **[Godotenv](github.com/joho/godotenv)** - loading env vars from .env files