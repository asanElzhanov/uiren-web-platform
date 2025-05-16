package exercises

func normalizeMultipleChoiceExerciseDTO[T ExerciseDTO](dto T, targetDTO T) error {
	if dto.GetOptions() == nil {
		return ErrOptionsRequired
	}
	if dto.GetCorrectAnswer() == nil {
		return ErrCorrectAnswerRequired
	}

	targetDTO.SetOptions(dto.GetOptions())
	targetDTO.SetCorrectAnswer(dto.GetCorrectAnswer())
	return nil
}

func normalizeManualTypingExerciseDTO[T ExerciseDTO](dto T, targetDTO T) error {
	if dto.GetCorrectAnswer() == nil {
		return ErrCorrectAnswerRequired
	}

	targetDTO.SetCorrectAnswer(dto.GetCorrectAnswer())
	return nil
}

func normalizeMatchPairsExerciseDTO[T ExerciseDTO](dto T, targetDTO T) error {
	if dto.GetPairs() == nil {
		return ErrPairsRequired
	}

	targetDTO.SetPairs(dto.GetPairs())
	return nil
}

func normalizeOrderWordsExerciseDTO[T ExerciseDTO](dto T, targetDTO T) error {
	if dto.GetOptions() == nil {
		return ErrOptionsRequired
	}
	if dto.GetCorrectOrder() == nil {
		return ErrCorrectOrderRequired
	}

	targetDTO.SetOptions(dto.GetOptions())
	targetDTO.SetCorrectOrder(dto.GetCorrectOrder())
	return nil
}
