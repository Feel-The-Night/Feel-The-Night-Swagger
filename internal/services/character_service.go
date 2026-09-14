package services

import (
	"errors"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CharacterService struct {
	db *gorm.DB
}

func NewCharacterService(db *gorm.DB) *CharacterService {
	return &CharacterService{db: db}
}

func (s *CharacterService) CreateCharacter(character *models.Character) error {
	var existingCharacter models.Character

	if err := s.db.Where("name = ?", character.Name).First(&existingCharacter).Error; err == nil {
		return errors.New("Personagem ja existente")
	}

	return s.db.Create(character).Error
}

func (s *CharacterService) GetCharacterById(id uint) (*models.Character, error) {
	var character models.Character
	if err := s.db.First(&character, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &character, nil
}

func (s *CharacterService) DeleteCharacterById(id uint) error {
	result := s.db.Delete(&models.Character{CharacterID: id})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *CharacterService) UpdateCharacterById(id uint, character *models.Character) (*models.Character, error) {
	existingCharacter, err := s.GetCharacterById(id)
	if err != nil {
		return nil, err
	}

	if character.Name != "" && character.Name != existingCharacter.Name {
		var checkName models.Character
		err := s.db.Where("name = ? AND character_id != ?", character.Name, id).First(&checkName).Error
		if err == nil {
			return nil, errors.New("Personagem ja existente")
		}
	}

	updates := make(map[string]any)

	if character.Name != "" {
		updates["name"] = character.Name
	}
	if character.Description != "" {
		updates["description"] = character.Description
	}
	if character.Type != "" {
		updates["type"] = character.Type
	}
	if character.ImageURL != "" {
		updates["image_url"] = character.ImageURL
	}

	zap.L().Info("[CharacterService] UpdateCharacter",
		zap.Any("updates", updates),
		zap.Any("updates_length", len(updates)),
	)

	if len(updates) == 0 {
		return nil, errors.New("Nenhum dado recebido")
	}

	if err := s.db.Model(existingCharacter).Updates(updates).Error; err != nil {
		return nil, err
	}

	return existingCharacter, nil

}
