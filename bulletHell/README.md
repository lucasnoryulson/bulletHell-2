# bulletHell

## Visão Geral

Esta é a primeira versão (v0.1) do **bulletHell**, um esqueleto de jogo *bullet hell* competitivo 2D em Go, rodando no terminal. Ele apresenta:

* Arena retangular desenhada em ASCII com bordas (`#`).
* Dois jogadores estáticos (`1` e `2`) e um "bullet" (`*`) no centro.
* Atualização em tempo real a uma taxa configurável de frames por segundo.
* Exibição do número do tick (frame) abaixo do mapa.

Esta versão inicial foca em:

1. **Render ASCII** no terminal usando códigos ANSI para limpar a tela.
2. **Loop de jogo** com `time.Ticker`, garantindo estabilidade de FPS.
3. **Entidades** básicas com coordenadas em grid e caractere associado.

---

## Pré-requisitos

* [Go](https://go.dev/dl/) (versão 1.18+ recomendada)
* Terminal compatível com ANSI escape codes (Linux, macOS, Windows PowerShell/WSL)

---

## Como Rodar

1. **Clone o repositório**

   ```bash
   git clone https://github.com/brunobaa/bulletHell-.git
   cd bulletHell
   ```

2. **Compile e execute**

   * **Sem gerar executável** (modo rápido):

     ```bash
     go run bulletHell.go
     ```

   * **Gerando binário** (modo release):

     ```bash
     # Windows
     go build -o bulletHell.exe bulletHell.go
     ./bulletHell.exe

     # macOS/Linux
     go build -o bulletHell bulletHell.go
     ./bulletHell
     ```

3. **Observe**

   * A cada frame, o terminal será limpo e redesenhado.
   * A borda da arena, os jogadores `1` e `2`, e o bullet `*` aparecerão.
   * Abaixo, o tick atual será mostrado.

---

## Configurações

Você pode ajustar as constantes no início de `bulletHell.go`:

```go
const (
    WorldWidth    = 30  // largura do mapa (colunas)
    WorldHeight   = 15  // altura do mapa (linhas)
    UpdatesPerSec = 30  // frames por segundo
)
```

* **WorldWidth/WorldHeight**: dimensionam sua arena.
* **UpdatesPerSec**: controla o ritmo de atualização.

---

## Próximos Passos

* Adicionar movimentação de players (captura de teclado em tempo real).
* Gerar e atualizar múltiplos bullets com padrões de disparo.
* Lógica de colisão e "morte" de players ao tocar balas.
* Estatísticas de partida: vidas, pontuação, tempo.

---

\_Agora você tem uma base limpa e modular para evoluir seu bulletHell \$1

---

## Contato

Para dúvidas ou contribuições, entre em contato: [brunoandradeprof7@gmail.com](mailto:brunoandradeprof7@gmail.com)
