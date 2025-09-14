package main

import (
	"math"
	"time"
)

func patoManager(jogo *Jogo) {
	moveTicker := time.NewTicker(1 * time.Second)
	defer moveTicker.Stop()

	for {
		select {
		case <-moveTicker.C:
			if isPortalAtivo() && !jogo.PatoInteragiu {
				moverPato(jogo)
			}
		case <-gameOverChannel:
			return
		}
	}
}

func moverPato(jogo *Jogo) {
	if jogo.PatoInteragiu {
		return
	}

	patoPosX, patoPosY := jogo.PatoPosX, jogo.PatoPosY

	dx := -1 // Move para cima
	novoY := patoPosY + dx

	// Verifica se a nova posição é válida
	if novoY >= 0 && novoY < len(jogo.Mapa) && novoY >= 0 && novoY < len(jogo.Mapa[novoY]) {
		if !jogo.Mapa[novoY][patoPosX].tangivel && jogo.Mapa[novoY][patoPosX].simbolo != Personagem.simbolo {
			// Move o pato
			cmd := func(jogo *Jogo) {
				// Remove o pato da posição atual
				jogo.Mapa[patoPosY][patoPosX] = Vazio
				// Coloca o pato na nova posição
				jogo.Mapa[novoY][patoPosX] = Pato
				// Atualiza a posição do pato no estado
				jogo.PatoPosY = novoY
			}
			mapChannel <- cmd
		}
	}
}

func interagirComPato(jogo *Jogo) {
	distanciaX := math.Abs(float64(jogo.PosX - jogo.PatoPosX))
	distanciaY := math.Abs(float64(jogo.PosY - jogo.PatoPosY))

	if (distanciaX <= 1 && distanciaY == 0) || (distanciaY <= 1 && distanciaX == 0) {
		cmd := func(jogo *Jogo) {
			jogo.PatoInteragiu = true
			jogo.StatusMsg = "Você fez carinho no pato! Ele parou de se mover."
		}
		mapChannel <- cmd
	}
}
