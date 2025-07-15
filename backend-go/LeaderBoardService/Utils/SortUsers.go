package Utils

import (
	"sort"
	"xxx/shared"
)

func SortUserScoresByScoreDesc(scores []shared.UserScore) []shared.UserScore {
	// Копируем, чтобы не менять оригинальный слайс
	sorted := make([]shared.UserScore, len(scores))
	copy(sorted, scores)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].TotalScore > sorted[j].TotalScore
	})

	return sorted
}
