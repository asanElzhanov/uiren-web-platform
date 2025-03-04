package exercises

import (
	"time"
)

const (
	multipleChoiceType = "multiple_choice"
	manualTypingType   = "manual_typing"
	matchPairsType     = "match_pairs"
	orderWordsType     = "order_words"
)

type Pair struct {
	Term  string `bson:"term" json:"term"`   // match_pairs
	Match string `bson:"match" json:"match"` // match_pairs
}

type Exercise struct {
	// common fields
	Code         string     `bson:"code" json:"code"`
	ExerciseType string     `bson:"type" json:"type"`
	Question     string     `bson:"question" json:"question"`
	Hints        []string   `bson:"hints" json:"hints"`
	Explanation  string     `bson:"explanation" json:"explanation"`
	CreatedAt    time.Time  `bson:"created_at" json:"created_at"`
	DeletedAt    *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`

	// multiple_choice, order_words
	Options []string `bson:"options,omitempty" json:"options,omitempty"`

	// multiple_choice, manual_typing
	CorrectAnswer string `bson:"correct_answer,omitempty" json:"correct_answer,omitempty"`

	// order_words
	CorrectOrder []string `bson:"correct_order,omitempty" json:"correct_order,omitempty"`

	// match_pairs
	Pairs []Pair `bson:"pairs,omitempty" json:"pairs,omitempty"`
}

// repo dto
type CreateExerciseDTO struct {
	Code          string     `bson:"code" json:"code"`
	ExerciseType  string     `bson:"type" json:"type"`
	Question      string     `bson:"question" json:"question"`
	Hints         []string   `bson:"hints" json:"hints"`
	Explanation   string     `bson:"explanation" json:"explanation"`
	CreatedAt     time.Time  `bson:"created_at" json:"created_at"`
	DeletedAt     *time.Time `bson:"deleted_at" json:"deleted_at"`
	Options       []string   `bson:"options,omitempty" json:"options,omitempty"`
	CorrectAnswer *string    `bson:"correct_answer,omitempty" json:"correct_answer,omitempty"`
	CorrectOrder  []string   `bson:"correct_order,omitempty" json:"correct_order,omitempty"`
	Pairs         []Pair     `bson:"pairs,omitempty" json:"pairs,omitempty"`
}

type UpdateExerciseDTO struct {
	Question      *string  `bson:"question,omitempty" json:"question,omitempty"`
	Hints         []string `bson:"hints,omitempty" json:"hints,omitempty"`
	Explanation   *string  `bson:"explanation,omitempty" json:"explanation,omitempty"`
	Options       []string `bson:"options,omitempty" json:"options,omitempty"`
	CorrectAnswer *string  `bson:"correct_answer,omitempty" json:"correct_answer,omitempty"`
	CorrectOrder  []string `bson:"correct_order,omitempty" json:"correct_order,omitempty"`
	Pairs         []Pair   `bson:"pairs,omitempty" json:"pairs,omitempty"`
}

type ExerciseDTO interface {
	GetOptions() []string
	GetCorrectAnswer() *string
	GetCorrectOrder() []string
	GetPairs() []Pair
	SetOptions([]string)
	SetCorrectAnswer(*string)
	SetCorrectOrder([]string)
	SetPairs([]Pair)
}

func (dto *CreateExerciseDTO) GetOptions() []string           { return dto.Options }
func (dto *CreateExerciseDTO) GetCorrectAnswer() *string      { return dto.CorrectAnswer }
func (dto *CreateExerciseDTO) GetCorrectOrder() []string      { return dto.CorrectOrder }
func (dto *CreateExerciseDTO) GetPairs() []Pair               { return dto.Pairs }
func (dto *CreateExerciseDTO) SetOptions(opts []string)       { dto.Options = opts }
func (dto *CreateExerciseDTO) SetCorrectAnswer(ans *string)   { dto.CorrectAnswer = ans }
func (dto *CreateExerciseDTO) SetCorrectOrder(order []string) { dto.CorrectOrder = order }
func (dto *CreateExerciseDTO) SetPairs(pairs []Pair)          { dto.Pairs = pairs }

func (dto *UpdateExerciseDTO) GetOptions() []string           { return dto.Options }
func (dto *UpdateExerciseDTO) GetCorrectAnswer() *string      { return dto.CorrectAnswer }
func (dto *UpdateExerciseDTO) GetCorrectOrder() []string      { return dto.CorrectOrder }
func (dto *UpdateExerciseDTO) GetPairs() []Pair               { return dto.Pairs }
func (dto *UpdateExerciseDTO) SetOptions(opts []string)       { dto.Options = opts }
func (dto *UpdateExerciseDTO) SetCorrectAnswer(ans *string)   { dto.CorrectAnswer = ans }
func (dto *UpdateExerciseDTO) SetCorrectOrder(order []string) { dto.CorrectOrder = order }
func (dto *UpdateExerciseDTO) SetPairs(pairs []Pair)          { dto.Pairs = pairs }
