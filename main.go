// main.go - Loop principal do jogo
package main

import (
	"os"
	"time"
)

var renderChannel = make(chan struct{}, 1)

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	go mapManager(&jogo)
	go coinManager(&jogo)
	go portalManager(&jogo)
	go patoManager(&jogo)
	go renderManager(&jogo)

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}

		select {
		case renderChannel <- struct{}{}:
		default:
		}
	}
}

func renderManager(jogo *Jogo) {
	// Timer para render automático a cada 100ms
	renderTicker := time.NewTicker(100 * time.Millisecond)
	defer renderTicker.Stop()

	for {
		select {
		case <-renderTicker.C:
			interfaceDesenharJogo(jogo)

		case <-renderChannel:
			interfaceDesenharJogo(jogo)

		case <-gameOverChannel:
			return
		}
	}
}
