# glox

Intérprete de Lox hecho en Go.

Belén Cerrotti [109566]

## Instalación

Requiere tener [Go](https://go.dev/dl/) instalado.

```bash
git clone <url-del-repo>
cd glox
go install
```

Esto compila el binario y lo instala en `$GOPATH/bin` (por defecto `~/go/bin`).  
Si `glox` no se reconoce como comando, agregá esa carpeta al PATH:

```bash
# en ~/.bashrc, ~/.zshrc, etc.
export PATH="$PATH:$HOME/go/bin"
```

Luego recargá la configuración (`source ~/.bashrc`) o abrí una terminal nueva.

## Uso

```bash
# REPL interactivo
glox

# Ejecutar un archivo
glox programa.lox

# Modos de debug
glox programa.lox --scanner   # imprime los tokens
glox programa.lox --parser    # imprime el AST
```
