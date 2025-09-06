package main

import "time"
import "math/rand"

func coinManager(jogo *Jogo) {
	var existingCoin bool = false
	var posicaoMoedaX, posicaoMoedaY int

	spawnTicker := time.NewTicker(10 * time.Second)
	defer spawnTicker.Stop()

	for {
		select {
		case <-spawnTicker.C:
			if existingCoin {
				clearCoin(jogo, posicaoMoedaX, posicaoMoedaY)
				existingCoin = false
			}
			spawned, x, y := spawnCoin(jogo)
			if spawned {
				existingCoin = true
				posicaoMoedaX = x
				posicaoMoedaY = y
			}
		case <-gameOverChannel:
			return
		}
	}
}

func spawnCoin(jogo *Jogo) (bool, int, int) {

	maxY := len(jogo.Mapa)
	maxX := len(jogo.Mapa[0])
	
	// Tenta spawnar a moeda indefinidamente até achar uma posição válida
	for {
		// Gera uma posição aleatória para x e y
		x := rand.Intn(maxX)
		y := rand.Intn(maxY)

		// Verifica se a posição é válida (vazio e não tangível)
		if !jogo.Mapa[y][x].tangivel && jogo.Mapa[y][x].simbolo != Moeda.simbolo {

			// Escreve o comando atualizando o mapa
			cmd := func(jogo *Jogo) {
				jogo.Mapa[y][x] = Moeda
			}
			
			// Envia o comando para o mapManager
			mapChannel <- cmd
			// Retorna a posição onde a moeda foi spawnada
			return true, x, y
		}
	}
}

func clearCoin(jogo *Jogo, x int, y int) {

	// Verifica se ainda há uma moeda na posição
	if jogo.Mapa[y][x].simbolo == Moeda.simbolo {

		// Escreve o comando atualizando o mapa
		cmd := func(jogo *Jogo) {
			jogo.Mapa[y][x] = Vazio
		}
		// Envia o comando para o mapManager
		mapChannel <- cmd
	}
}