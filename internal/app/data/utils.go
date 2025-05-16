package data

import "strconv"

func generateXpLeaderboardKey(limit int) string {
	return "xp_leaderboard_" + strconv.FormatInt(int64(limit), 10)
}

func generateLessonKey(code string) string {
	return "lesson::" + code
}

func generateExerciseKey(code string) string {
	return "exercise::" + code
}
