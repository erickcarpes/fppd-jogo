package main

import (
	"time"
)

var cachorroInteragiu = false

func cachorroManager(jogo *Jogo) {
	moveTicker := time.NewTicker(1 * time.Second)
	defer moveTicker.Stop()

	for {
		select {
		case <-moveTicker.C:
			if isPortalAtivo() && !cachorroInteragiu {
				// Move o cachorro uma posição para frente (direita)
				moverCachorro(jogo)
			}
		case <-gameOverChannel:
			return
		}
	}
}

func moverCachorro(jogo *Jogo) {
	// Encontra a posição atual do cachorro no mapa
	var cachorroPosX, cachorroPosY int = -1, -1

	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			if elem.simbolo == Cachorro.simbolo {
				cachorroPosX = x
				cachorroPosY = y
				break
			}
		}
		if cachorroPosX != -1 {
			break
		}
	}

	// Se não encontrou o cachorro, não faz nada
	if cachorroPosX == -1 {
		return
	}

	dx := 1 // Move para a direita
	novoX := cachorroPosX + dx
	novoY := cachorroPosY

	// Verifica limites do mapa
	if novoX >= len(jogo.Mapa[0]) {
		// Se chegou na borda direita, volta para a esquerda
		novoX = 0
	}

	// Verifica se a nova posição é válida
	if novoY >= 0 && novoY < len(jogo.Mapa) && novoX >= 0 && novoX < len(jogo.Mapa[novoY]) {
		if !jogo.Mapa[novoY][novoX].tangivel && jogo.Mapa[novoY][novoX].simbolo != Personagem.simbolo {
			// Move o cachorro
			cmd := func(jogo *Jogo) {
				// Remove o cachorro da posição atual
				jogo.Mapa[cachorroPosY][cachorroPosX] = Vazio
				// Coloca o cachorro na nova posição
				jogo.Mapa[novoY][novoX] = Cachorro
			}
			mapChannel <- cmd
		}
	}
}

func interagirComCachorro(jogo *Jogo) {
	// Verifica se há um cachorro adjacente ao jogador
	adjacentes := []struct{ dx, dy int }{
		{0, -1}, // cima
		{0, 1},  // baixo
		{-1, 0}, // esquerda
		{1, 0},  // direita
	}

	for _, adj := range adjacentes {
		x := jogo.PosX + adj.dx
		y := jogo.PosY + adj.dy

		if x >= 0 && x < len(jogo.Mapa[0]) && y >= 0 && y < len(jogo.Mapa) {
			if jogo.Mapa[y][x].simbolo == Cachorro.simbolo {
				cachorroInteragiu = true
				jogo.StatusMsg = "Você fez carinho no cachorro! Ele parou de se mover."
				return
			}
		}
	}

	// Se não encontrou cachorro adjacente, verifica se o jogador está na mesma posição
	if jogo.Mapa[jogo.PosY][jogo.PosX].simbolo == Cachorro.simbolo {
		cachorroInteragiu = true
		jogo.StatusMsg = "Você fez carinho no cachorro! Ele parou de se mover."
	}
}
