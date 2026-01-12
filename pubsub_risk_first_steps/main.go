package main

// subscriber - written in 1.2 lession
func (u user) doBattles(subCh <-chan move) []piece {
	fights := []piece{}
	for mv := range subCh {
		for _, piece := range u.pieces {
			if piece.location == mv.piece.location {
				fights = append(fights, piece)
			}
		}
	}
	return fights
}

// don't touch below this line

type user struct {
	name   string
	pieces []piece
}

type move struct {
	userName string
	piece    piece
}

type piece struct {
	location string
	name     string
}

// publisher - written in 1.1 lession
func (u user) march(p piece, publishCh chan<- move) {
	publishCh <- move{
		userName: u.name,
		piece:    p,
	}
}

// distribue what publisher give, and send it to subscribers
func distributeBattles(publishCh <-chan move, subChans []chan move) {
	for mv := range publishCh {
		for _, subCh := range subChans {
			subCh <- mv
		}
	}
}
