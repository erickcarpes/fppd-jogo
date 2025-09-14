// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"fmt"
	"os"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa           [][]Elemento // grade 2D representando o mapa
	PosX, PosY     int          // posição atual do personagem
	UltimoVisitado Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg      string       // mensagem para a barra de status
}

// Elementos visuais do jogo
var (
	Personagem    = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo       = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede        = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao     = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio         = Elemento{' ', CorPadrao, CorPadrao, false}
	Moeda         = Elemento{'ၜ', CorAmarelo, CorPadrao, false}
	PortalAtivo   = Elemento{'○', CorMagenta, CorPadrao, false}
	PortalInativo = Vazio
	Cachorro      = Elemento{'ࠎ', CorMarrom, CorPadrao, false}
)

var coinChannel = make(chan struct{})
var portalChannel = make(chan bool)
var mapChannel = make(chan func(*Jogo))
var gameOverChannel = make(chan struct{})

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{UltimoVisitado: Vazio}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Inimigo
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
			case Cachorro.simbolo:
				e = Cachorro
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição e retorna se houve teletransporte
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) bool {
	nx, ny := x+dx, y+dy
	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x]
	elementoNaNovaPosicao := jogo.Mapa[ny][nx]

	switch elementoNaNovaPosicao.simbolo {
	case Moeda.simbolo:
		jogo.Mapa[y][x] = jogo.UltimoVisitado // restaura o conteúdo anterior
		jogo.UltimoVisitado = Vazio
		jogo.Mapa[ny][nx] = elemento // move o elemento
		select {
		case portalChannel <- true:
		default:
		}
		return false

	case PortalAtivo.simbolo:
		// Encontra uma nova posição aleatória para o teletransporte
		newX, newY := teleportarJogador(jogo)
		jogo.StatusMsg = fmt.Sprintf("Teletransportado para (%d, %d)!", newX, newY)

		// Restaura a posição anterior do jogador
		jogo.Mapa[y][x] = jogo.UltimoVisitado // restaura o conteúdo anterior

		// Remove o portal (ele é consumido)
		jogo.Mapa[ny][nx] = Vazio

		// Guarda o elemento da nova posição e coloca o jogador lá
		jogo.UltimoVisitado = jogo.Mapa[newY][newX] // guarda o conteúdo atual da nova posição
		jogo.Mapa[newY][newX] = elemento            // move o elemento

		// Atualiza a posição do jogador
		jogo.PosX, jogo.PosY = newX, newY

		return true // indica que houve teletransporte

	default:
		jogo.Mapa[y][x] = jogo.UltimoVisitado   // restaura o conteúdo anterior
		jogo.UltimoVisitado = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
		jogo.Mapa[ny][nx] = elemento            // move o elemento
		return false                            // movimento normal
	}
}

func mapManager(Jogo *Jogo) {
	for {
		select {
		case cmd := <-mapChannel:
			cmd(Jogo)
		case <-gameOverChannel:
			return
		}
	}
}
