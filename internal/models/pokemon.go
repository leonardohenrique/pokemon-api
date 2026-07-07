package models

type Pokemon struct {
	ID    int      `json:"id"`
	Name  string   `json:"name" validate:"required,min=2,max=50"`
	Type  []string `json:"type" validate:"required,min=1,max=2,dive,oneof=normal fire water electric grass ice fighting poison ground flying psychic bug rock ghost dragon dark steel fairy"`
	Level int      `json:"level" validate:"required,min=1,max=100"`
	HP    int      `json:"hp" validate:"required,min=1,max=999"`
}
