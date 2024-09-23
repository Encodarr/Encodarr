package repository

import (
	"errors"
	"log"
	"transfigurr/models"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	DB *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		DB: db,
	}
}

func (repo *ProfileRepository) GetAllProfiles() ([]models.Profile, error) {
	var profiles []models.Profile
	if err := repo.DB.Preload("ProfileAudioLanguages").Preload("ProfileSubtitleLanguages").Preload("ProfileCodecs").Find(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}

func (repo *ProfileRepository) GetProfileById(profileId int) (models.Profile, error) {
	var profile models.Profile
	if err := repo.DB.Preload("ProfileAudioLanguages").Preload("ProfileSubtitleLanguages").Preload("ProfileCodecs").Where("id = ?", profileId).First(&profile).Error; err != nil {
		return models.Profile{}, err
	}
	return profile, nil
}

func (repo *ProfileRepository) UpsertProfile(profileId int, inputProfile models.Profile) (models.Profile, error) {
	var profile models.Profile

	if repo.DB == nil {
		log.Printf("Database connection is nil")
		return models.Profile{}, errors.New("database connection is nil")
	}

	result := repo.DB.Where("id = ?", profileId).First(&profile)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		profile = inputProfile
		if err := repo.DB.Create(&profile).Error; err != nil {
			log.Printf("Error creating profile: %v", err)
			return models.Profile{}, err
		}
	} else {
		repo.DB.Model(&profile).Select("*").Updates(inputProfile)
		if err := repo.DB.Save(&profile).Error; err != nil {
			log.Printf("Error updating profile: %v", err)
			return models.Profile{}, err
		}
	}

	// Handle ProfileAudioLanguages
	var existingAudioLanguages []models.ProfileAudioLanguage
	repo.DB.Where("profile_id = ?", profileId).Find(&existingAudioLanguages)

	// Find audio languages to delete
	var audioLanguagesToDelete []models.ProfileAudioLanguage
	for _, existingAudioLanguage := range existingAudioLanguages {
		found := false
		for _, newAudioLanguage := range inputProfile.ProfileAudioLanguages {
			if existingAudioLanguage.Language == newAudioLanguage.Language {
				found = true
				break
			}
		}
		if !found {
			audioLanguagesToDelete = append(audioLanguagesToDelete, existingAudioLanguage)
		}
	}

	// Delete unwanted audio languages
	for _, audioLanguageToDelete := range audioLanguagesToDelete {
		repo.DB.Delete(&audioLanguageToDelete)
	}

	// Replace with new audio languages
	if err := repo.DB.Model(&profile).Association("ProfileAudioLanguages").Replace(inputProfile.ProfileAudioLanguages); err != nil {
		return models.Profile{}, err
	}

	// Handle ProfileSubtitleLanguages
	var existingSubtitleLanguages []models.ProfileSubtitleLanguage
	repo.DB.Where("profile_id = ?", profileId).Find(&existingSubtitleLanguages)

	// Find subtitle languages to delete
	var subtitleLanguagesToDelete []models.ProfileSubtitleLanguage
	for _, existingSubtitleLanguage := range existingSubtitleLanguages {
		found := false
		for _, newSubtitleLanguage := range inputProfile.ProfileSubtitleLanguages {
			if existingSubtitleLanguage.Language == newSubtitleLanguage.Language {
				found = true
				break
			}
		}
		if !found {
			subtitleLanguagesToDelete = append(subtitleLanguagesToDelete, existingSubtitleLanguage)
		}
	}

	// Delete unwanted subtitle languages
	for _, subtitleLanguageToDelete := range subtitleLanguagesToDelete {
		repo.DB.Delete(&subtitleLanguageToDelete)
	}

	// Replace with new subtitle languages
	if err := repo.DB.Model(&profile).Association("ProfileSubtitleLanguages").Replace(inputProfile.ProfileSubtitleLanguages); err != nil {
		log.Printf("Error updating ProfileSubtitleLanguages: %v", err)
		return models.Profile{}, err
	}

	// Handle ProfileCodecs
	var existingCodecs []models.ProfileCodec
	repo.DB.Where("profile_id = ?", profileId).Find(&existingCodecs)

	// Find codecs to delete
	var codecsToDelete []models.ProfileCodec
	for _, existingCodec := range existingCodecs {
		found := false
		for _, newCodec := range inputProfile.ProfileCodecs {
			if existingCodec.CodecId == newCodec.CodecId {
				found = true
				break
			}
		}
		if !found {
			codecsToDelete = append(codecsToDelete, existingCodec)
		}
	}

	// Delete unwanted codecs
	for _, codecToDelete := range codecsToDelete {
		repo.DB.Delete(&codecToDelete)
	}

	// Replace with new codecs
	if err := repo.DB.Model(&profile).Association("ProfileCodecs").Replace(inputProfile.ProfileCodecs); err != nil {
		log.Printf("Error updating ProfileCodecs: %v", err)
		return models.Profile{}, err
	}

	return profile, nil
}

func (repo *ProfileRepository) DeleteProfileById(profileId int) error {
	var profile models.Profile
	if err := repo.DB.Where("id = ?", profileId).First(&profile).Error; err != nil {
		return err
	}

	// Delete associated codecs
	if err := repo.DB.Where("profile_id = ?", profileId).Delete(&models.ProfileCodec{}).Error; err != nil {
		return err
	}

	// Delete associated audio languages
	if err := repo.DB.Where("profile_id = ?", profileId).Delete(&models.ProfileAudioLanguage{}).Error; err != nil {
		return err
	}

	// Delete associated subtitle languages
	if err := repo.DB.Where("profile_id = ?", profileId).Delete(&models.ProfileSubtitleLanguage{}).Error; err != nil {
		return err
	}

	// Delete the profile
	if err := repo.DB.Delete(&profile).Error; err != nil {
		return err
	}

	return nil
}
