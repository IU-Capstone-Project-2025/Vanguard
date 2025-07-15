package Utils

import (
	"sort"
	"xxx/LeaderBoardService/models"
)

func SortUserScoresByScoreDesc(scores []models.UserScore) []models.UserScore {
	// Копируем, чтобы не менять оригинальный слайс
	sorted := make([]models.UserScore, len(scores))
	copy(sorted, scores)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].TotalScore > sorted[j].TotalScore
	})

	return sorted
}
