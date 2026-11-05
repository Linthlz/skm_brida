package repositories

import (
	"errors"
	"strconv"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ---------- Funding schemes ----------

func ListFunding(status, scheme string, offset, limit int) ([]models.FundingScheme, int64, error) {
	var list []models.FundingScheme
	var total int64
	db := config.DB.Model(&models.FundingScheme{})
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if scheme != "" {
		db = db.Where("scheme_code = ?", scheme)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("CASE status WHEN 'open' THEN 0 WHEN 'closing' THEN 1 ELSE 2 END, deadline ASC").
		Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func FetchFunding(publicID string) (*models.FundingScheme, error) {
	var f models.FundingScheme
	err := config.DB.Where("public_id = ?", publicID).First(&f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &f, err
}

func CreateFunding(f *models.FundingScheme) error { return config.DB.Create(f).Error }

func UpdateFunding(f *models.FundingScheme) error {
	return config.DB.Model(&models.FundingScheme{}).Where("funding_id = ?", f.InternalID).Updates(map[string]interface{}{
		"name": f.Name, "organizer": f.Organizer, "scheme_code": f.SchemeCode, "ceiling": f.Ceiling,
		"quota": f.Quota, "deadline": f.Deadline, "status": f.Status, "target_fields": f.TargetFields,
		"requirements": f.Requirements, "description": f.Description, "website": f.Website, "website_name": f.WebsiteName,
	}).Error
}

func DeleteFunding(publicID string) error {
	return config.DB.Where("public_id = ?", publicID).Delete(&models.FundingScheme{}).Error
}

// ---------- Site settings ----------

// SettingDefaults lists every supported key; anything else is rejected on save.
var SettingDefaults = map[string]string{
	"site_name":             "SINGA RISET BULELENG",
	"tagline":               "Sinergi Gerakan Akademisi dan Riset Buleleng",
	"contact_email":         "brida@bulelengkab.go.id",
	"contact_phone":         "(0362) 21985",
	"address":               "Jl. Ngurah Rai No. 2, Singaraja, Kabupaten Buleleng, Bali 81113",
	"batch_open":            "true",
	"registration_open":     "true",
	"maintenance_mode":      "false",
	"stat_publikasi":        "0",
	"stat_hki":              "0",
	"stat_desa_terdampak":   "0",
	"stat_adopsi_kebijakan": "0",
	"stat_policy_brief":     "0",
}

var BoolSettings = map[string]bool{"batch_open": true, "registration_open": true, "maintenance_mode": true}

func IntSetting(key string) bool {
	_, ok := map[string]bool{"stat_publikasi": true, "stat_hki": true, "stat_desa_terdampak": true, "stat_adopsi_kebijakan": true, "stat_policy_brief": true}[key]
	return ok
}

func GetSettings() (map[string]string, error) {
	out := make(map[string]string, len(SettingDefaults))
	for k, v := range SettingDefaults {
		out[k] = v
	}
	var rows []models.SiteSetting
	if err := config.DB.Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		if _, ok := SettingDefaults[r.Key]; ok {
			out[r.Key] = r.Value
		}
	}
	return out, nil
}

func SaveSettings(values map[string]string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		for k, v := range values {
			row := models.SiteSetting{Key: k, Value: v}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "setting_key"}},
				DoUpdates: clause.AssignmentColumns([]string{"setting_value"}),
			}).Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SettingBool reads a boolean setting; on any error it returns the default so
// a settings problem never blocks the whole API.
func SettingBool(key string) bool {
	values, err := GetSettings()
	if err != nil {
		v, _ := strconv.ParseBool(SettingDefaults[key])
		return v
	}
	v, _ := strconv.ParseBool(values[key])
	return v
}
