package repositories

import (
	"strconv"

	"github.com/itsXanadium/singa-riset-backend/config"
)

type YearCount struct {
	Year  int   `json:"year"`
	Total int64 `json:"total"`
}

type PortalStats struct {
	TotalResearch      int64           `json:"total_research"`
	ActiveResearch     int64           `json:"active_research"`
	PendingSubmissions int64           `json:"pending_submissions"`
	Researchers        int64           `json:"researchers"`
	Institutions       int64           `json:"institutions"`
	Users              int64           `json:"users"`
	Documentations     int64           `json:"documentations"`
	PublishedNews      int64           `json:"published_news"`
	OpenFunding        int64           `json:"open_funding"`
	ActiveBudget       int64           `json:"active_budget"`
	ByCategory         []ResearchFacet `json:"by_category"`
	ByLocation         []ResearchFacet `json:"by_location"`
	ByYear             []YearCount     `json:"by_year"`
	Manual             map[string]int  `json:"manual"`
}

var activeStatuses = []string{"on-going", "approved"}

func PortalStatistics() (*PortalStats, error) {
	s := &PortalStats{Manual: map[string]int{}, ByYear: []YearCount{}, ByCategory: []ResearchFacet{}, ByLocation: []ResearchFacet{}}
	db := config.DB
	steps := []func() error{
		func() error {
			return db.Table("researches").Where("deleted_at IS NULL").Count(&s.TotalResearch).Error
		},
		func() error {
			return db.Table("researches").Where("deleted_at IS NULL AND status IN ?", activeStatuses).Count(&s.ActiveResearch).Error
		},
		func() error {
			return db.Table("researches").Where("deleted_at IS NULL AND status = ?", "pending").Count(&s.PendingSubmissions).Error
		},
		func() error {
			return db.Table("users").Where("role = ? AND verified = true", "researcher").Count(&s.Researchers).Error
		},
		func() error {
			return db.Table("users").Count(&s.Users).Error
		},
		func() error {
			return db.Table("users").Where("role = ? AND COALESCE(TRIM(institution), '') <> ''", "researcher").
				Select("COUNT(DISTINCT LOWER(TRIM(institution)))").Scan(&s.Institutions).Error
		},
		func() error {
			return db.Table("group_documentations").Where("deleted_at IS NULL").Count(&s.Documentations).Error
		},
		func() error {
			return db.Table("news").Where("deleted_at IS NULL AND status = ?", "published").Count(&s.PublishedNews).Error
		},
		func() error {
			return db.Table("funding_schemes").Where("deleted_at IS NULL AND status IN ?", []string{"open", "closing"}).Count(&s.OpenFunding).Error
		},
		func() error {
			return db.Table("researches").Where("deleted_at IS NULL AND status IN ?", activeStatuses).
				Select("COALESCE(SUM(CASE WHEN fund_amount ~ '^[0-9]+$' THEN fund_amount::bigint ELSE 0 END), 0)").
				Scan(&s.ActiveBudget).Error
		},
		func() error {
			return db.Table("researches").Where("deleted_at IS NULL").
				Select("EXTRACT(YEAR FROM created_at)::int AS year, COUNT(*) AS total").
				Group("year").Order("year").Scan(&s.ByYear).Error
		},
	}
	for _, step := range steps {
		if err := step(); err != nil {
			return nil, err
		}
	}
	var err error
	if s.ByCategory, err = ResearchFacets("category", activeStatuses); err != nil {
		return nil, err
	}
	if s.ByLocation, err = ResearchFacets("location", activeStatuses); err != nil {
		return nil, err
	}
	settings, err := GetSettings()
	if err != nil {
		return nil, err
	}
	for k := range SettingDefaults {
		if IntSetting(k) {
			n, _ := strconv.Atoi(settings[k])
			s.Manual[k] = n
		}
	}
	return s, nil
}
