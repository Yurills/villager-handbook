package engine

import "github.com/Yurills/villager-handbook/internal/model"

func (e *Engine) LookaheadBestCandidate() []model.LookaheadResult { //return player ID
	playersEntropy := []model.LookaheadResult{}

	for playerID, playerState := range e.Players {

		if playerState == -1 {
			continue //skip eliminated players
		}

		originalStats := e.GetPlayerRoleProbabilities(playerID)
		evilProb := originalStats.RoleProbabilities[model.Werewolf]
		goodProb := 1.0 - evilProb

		if evilProb > 0.98 || goodProb < 0.02 {
			continue //skip players with extreme probabilities
		}

		//simulate
		entropyIfEvil := 0.0
		if evilProb > 0 {
			simEngine := e.fork()
			simEngine.ProcessMove(
				model.Interaction{
					Actor:  playerID,
					Target: playerID,
					Type:   "fact",
					Result: "evil",
				})
			entropyIfEvil = simEngine.calculateEntropy(simEngine.GetStats())
		}
		entropyIfGood := 0.0
		if goodProb > 0 {
			simEngine := e.fork()
			simEngine.ProcessMove(
				model.Interaction{
					Actor:  playerID,
					Target: playerID,
					Type:   "fact",
					Result: "good",
				})
			entropyIfGood = simEngine.calculateEntropy(simEngine.GetStats())

		}

		expectedEntropy := evilProb*entropyIfEvil + goodProb*entropyIfGood

		playersEntropy = append(playersEntropy, model.LookaheadResult{
			ID:      playerID,
			Entropy: expectedEntropy,
		})

	}
	return playersEntropy

}
