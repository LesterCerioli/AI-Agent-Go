package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRequirement struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Prompt    string    `json:"prompt" gorm:"type:text"`
	Status    string    `json:"status" gorm:"type:varchar(50)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Architecture *ProjectArchitecture `json:"architecture,omitempty" gorm:"foreignKey:RequirementID"`
}

func (pr *ProjectRequirement) BeforeCreate(tx *gorm.DB) (err error) {
	pr.ID = uuid.New()
	return
}

type ProjectArchitecture struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	RequirementID uuid.UUID `json:"requirement_id" gorm:"type:uuid"`
	TechStack     string    `json:"tech_stack" gorm:"type:varchar(100)"` // e.g., "Go, Fiber, Postgres"
	Language      string    `json:"language" gorm:"type:varchar(50)"`    // e.g., "Go"
	Structure     string    `json:"structure" gorm:"type:text"`          // JSON representation of directory tree
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Files []GeneratedFile `json:"files,omitempty" gorm:"foreignKey:ArchitectureID"`
}

func (pa *ProjectArchitecture) BeforeCreate(tx *gorm.DB) (err error) {
	pa.ID = uuid.New()
	return
}

type GeneratedFile struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	ArchitectureID uuid.UUID `json:"architecture_id" gorm:"type:uuid"`
	FilePath       string    `json:"file_path" gorm:"type:varchar(255)"`
	Content        string    `json:"content" gorm:"type:text"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (gf *GeneratedFile) BeforeCreate(tx *gorm.DB) (err error) {
	gf.ID = uuid.New()
	return
}
