package main

import (
	"math/rand"
	"time"
)

var portalConsultaChannel = make(chan chan bool)

func portalManager(jogo *Jogo) {
	var portalAtivo bool = false
	var posicaoPortalX, posicaoPortalY int

	for {
		select {
			case ativar := <-portalChannel:
				if ativar && !portalAtivo {
					portalAtivo, posicaoPortalX, posicaoPortalY = ativarPortal(jogo)

					go func() {
						time.Sleep(15 * time.Second)
						portalChannel <- false
					}()	
				}else {
					if portalAtivo {
						clearPortal(jogo, posicaoPortalX, posicaoPortalY)
						portalAtivo = false
					}
				}
			case resposta := <-portalConsultaChannel:
				resposta <- portalAtivo

			case <-gameOverChannel:
				return
		}
	}
}

func ativarPortal(jogo *Jogo) (bool, int, int) {

	maxY := len(jogo.Mapa)
	maxX := len(jogo.Mapa[0])

	for {
		// Gera uma posição aleatória para x e y
		x := rand.Intn(maxX)
		y := rand.Intn(maxY)

		// Verifica se a posição é válida (vazio e não tangível)	
		if !jogo.Mapa[y][x].tangivel && jogo.Mapa[y][x].simbolo == Vazio.simbolo {

			// Escreve o comando atualizando o mapa
			cmd := func(jogo *Jogo) {
				jogo.Mapa[y][x] = PortalAtivo
			}
			// Envia o comando para o mapManager
			mapChannel <- cmd
			// Retorna a posição onde o portal foi spawnado
			return true, x, y
		}
	}
}

func clearPortal(jogo *Jogo, x int, y int) {
	// Verifica se ainda há um portal na posição
	if jogo.Mapa[y][x].simbolo == PortalAtivo.simbolo {
		// Escreve o comando atualizando o mapa
		cmd := func(jogo *Jogo) {
			jogo.Mapa[y][x] = PortalInativo
		}
		// Envia o comando para o mapManager
		mapChannel <- cmd
	}
}

func teleportarJogador(jogo *Jogo) (int, int) {
    maxY := len(jogo.Mapa)
    maxX := len(jogo.Mapa[0])

    for {
        x := rand.Intn(maxX)
        y := rand.Intn(maxY)

        if !jogo.Mapa[y][x].tangivel && jogo.Mapa[y][x].simbolo == Vazio.simbolo {
            return x, y
        }
    }
}

func isPortalAtivo() bool {
    resposta := make(chan bool)
    portalConsultaChannel <- resposta
    return <-resposta
}