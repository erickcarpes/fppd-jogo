package main

import (
	"fmt"
	"math"
	"time"
)

func patoManager(jogo *Jogo) {
	moveTicker := time.NewTicker(1 * time.Second)
	defer moveTicker.Stop()

	for {
		select {
		case <-moveTicker.C:
			tentarMoverPato(jogo)
		case <-gameOverChannel:
			return
		}
	}
}

func tentarMoverPato(jogo *Jogo) {
	cmd := func(jogo *Jogo) {

		portalAtivo := jogo.PortalAtivo
		patoInteragiu := jogo.PatoInteragiu

		// Se portal não está ativo ou pato foi interagido, não move
		if !portalAtivo {
			return
		}

		if patoInteragiu {
			return
		}

		moverPato(jogo)
	}

	mapChannel <- cmd
}

func moverPato(jogo *Jogo) {

	patoPosX, patoPosY := jogo.PatoPosX, jogo.PatoPosY
	jogo.StatusMsg = fmt.Sprintf("Pato está em (%d, %d)", patoPosX, patoPosY)
	novoY := patoPosY - 1

	// Verifica se a nova posição é válida
	if novoY >= 0 && novoY < len(jogo.Mapa) && novoY >= 0 && novoY < len(jogo.Mapa[novoY]) {
		if !jogo.Mapa[novoY][patoPosX].tangivel && jogo.Mapa[novoY][patoPosX].simbolo != Personagem.simbolo {
			// Remove o pato da posição atual
			jogo.Mapa[jogo.PatoPosY][jogo.PatoPosX] = Vazio
			// Coloca o pato na nova posição
			jogo.Mapa[novoY][patoPosX] = Pato
			// Atualiza a posição do pato no estado
			jogo.PatoPosY = novoY
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
