# email-search-engine-backend

Este es un proyecto para busqueda rápida de correos electrónicos.

Utiliza Vue.js para el frontend y Go y ZincSearch para el backend.

Para correr el proyecto:
```
go main.go
```

Para correr Zincsearch:
```
set ZINC_FIRST_ADMIN_USER=admin
set ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123
set ZINC_MAX_DOCUMENT_SIZE=2097152
mkdir data
zinc.exe
```

El frontend se encuentra en el repositorio: https://github.com/Miguel219/email-search-engine
