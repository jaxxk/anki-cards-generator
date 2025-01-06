package create

var DefaultVersion int = 6
var FlashcardBatchSize int = 50

type AnkiRequestBody struct {
	Action  string     `json:"action"`
	Version int        `json:"version"`
	Params  AnkiParams `json:"params"`
}

type AnkiParams interface{}

type Notes struct {
	ListOfNotes []Note `json:"notes"`
}

type Note struct {
	DeckName  string            `json:"deckName"`
	ModelName string            `json:"modelName"`
	Fields    map[string]string `json:"fields"`
}

func NewNote(front, back, deckName string) Note {
	return Note{
		DeckName:  deckName,
		ModelName: "Basic",
		Fields: map[string]string{
			"Front": front,
			"Back":  back,
		},
	}
}

func NewNotes() Notes {
	return Notes{
		ListOfNotes: make([]Note, 10),
	}
}

func NewAnkiRequestBody(action string, params AnkiParams) AnkiRequestBody {
	return AnkiRequestBody{
		Action:  action,
		Version: 6,
		Params:  params,
	}
}
