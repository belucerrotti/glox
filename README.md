# glox

Intérprete del lenguaje **Lox** escrito en Go.

Belén Cerrotti [109566]

---

## Diferencias respecto a la implementación de la materia

### ¿Qué agrega glox?


### ¿Qué no implementa (respecto al Lox completo)?

| Feature | Estado |
|---|---|
| Clases y objetos (`class`, `this`, `super`) | ❌ No implementado |
| Funciones nativas (`clock()`) | ❌ No implementado |
| Strings con escape sequences (`\n`, `\t`) | ❌ No implementado |

---

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

---

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

---

## Demo

El archivo `demo.lox` muestra las principales features del lenguaje:
variables, tipos, aritmética (incluido `%`), strings, booleanos,
`if/else`, `while`, `for`, funciones, recursión, closures y scoping léxico.

```bash
go build . && ./glox demo.lox
```

Output esperado:
```
13
7
30
3.33333
1
Hola, mundo!
false
true
true
x es mayor a 5
0
1
2
0
1
2
11
1
2
3
55
global
global
```

---

## Tests

Los tests están en `tests/` y se corren con:

```bash
go build . && python3 tests/script.py
```

El script ejecuta todos los archivos `.lox` de las carpetas `basic/` y `advanced/` en orden y verifica que ninguno falle.

**basic/**
- `0-simple` — aritmética, strings, booleanos y comparaciones básicas.
- `1-flow` — estructuras de control: `if/else`, `while`, `for`.
- `2-functions` — declaración y llamada de funciones, recursión y valores de retorno.
- `3-minsky` — simulación de máquina de Minsky (programa no trivial como integración).
- `4-fizzbuzz` — FizzBuzz clásico como test de control de flujo combinado.

**advanced/**
- `0-scopes` — scoping léxico, variables locales vs globales, y closures con estado compartido.
- `1-classes` — clases, `init`, métodos, acceso y asignación de campos, múltiples instancias independientes.
- `2-class-scopes` — interacción entre `this`, variables locales dentro de métodos, referencias a métodos y llamadas entre métodos.

---

## Benchmark

Loop de **1.000.000 iteraciones** sumando enteros, comparado con lenguajes productivos:

```bash
# glox
go build . && time ./glox benchmark.lox

# Python 3
time python3 -c "i=0;s=0
while i<1000000: s+=i;i+=1
print(s)"

# C (gcc -O0, sin optimizaciones)
cat > /tmp/bench.c << 'EOF'
#include <stdio.h>
int main() {
    long s = 0;
    for (int i = 0; i < 1000000; i++) s += i;
    printf("%ld\n", s);
}
EOF
gcc -O0 /tmp/bench.c -o /tmp/bench && time /tmp/bench
```

| Lenguaje | Tiempo (real) | Notas |
|---|---|---|
| C (gcc -O0) | ~0.003s | compilado a código máquina |
| Python 3.12 | ~0.07s | compilado a bytecode + VM |
| **glox** | **~0.92s** | tree-walk interpreter |

glox es ~300× más lento que C y ~13× más lento que Python. Esto es esperable: glox es un *tree-walk interpreter*, lo que significa que cada operación requiere recorrer el AST nodo por nodo, hacer un `switch` sobre el tipo, y alocar structs en Go.
