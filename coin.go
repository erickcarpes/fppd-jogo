package main

import "time"
import "math/rand"

func coinSpawner(jogo *Jogo) {

	spawnTicker := time.NewTicker(10 * time.Second)
	defer spawnTicker.Stop()

	for {
		select {
		case <-spawnTicker.C:
			spawnCoin(jogo)
		case <-gameOverChannel:
			return
		}
	}
}

func spawnCoin(jogo *Jogo) {

	maxY := len(jogo.Mapa)
	maxX := len(jogo.Mapa[0])

	x := rand.Intn(maxX)
	y := rand.Intn(maxY)

	cmd := func(jogo *Jogo) {
		if !jogo.Mapa[y][x].tangivel && jogo.Mapa[y][x].simbolo != Moeda.simbolo {
            jogo.Mapa[y][x] = Moeda
        }
	}

	mapChannel <- cmd
}