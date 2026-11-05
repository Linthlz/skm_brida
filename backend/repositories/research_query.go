package repositories

import (
	"errors"

	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/models"
	"gorm.io/gorm"
)

type ResearchQuery struct {
	Statuses []string
	Filter   string
	Category string
	Location string
	Sort     string
	Offset   int
	Limit    int
}

func applyResearchQuery(db *gorm.DB, q ResearchQuery) *gorm.DB {
	if len(q.Statuses) > 0 {
		db = db.Where("status IN ?", q.Statuses)
	}
	if q.Filter != "" {
		db = db.Where("title ILIKE ?", "%"+q.Filter+"%")
	}
	if q.Category != "" {
		db = db.Where("category = ?", q.Category)
	}
	if q.Location != "" {
		db = db.Where("location = ?", q.Location)
	}
	return db
}

// SearchResearches lists researches with optional status/category/location filters.
func SearchResearches(q ResearchQuery) ([]models.Research, int64, error) {
	var list []models.Research
	var total int64
	base := applyResearchQuery(config.DB.Model(&models.Research{}).Where("deleted_at IS NULL"), q)
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	order := "created_at DESC"
	switch q.Sort {
	case "oldest":
		order = "created_at ASC"
	case "title":
		order = "title ASC"
	}
	err := applyResearchQuery(config.DB.Model(&models.Research{}).Where("deleted_at IS NULL"), q).
		Preload("Uploader").Preload("Reviewer").Preload("Group").
		Order(order).Offset(q.Offset).Limit(q.Limit).Find(&list).Error
	return list, total, err
}

// ResearchFacet counts researches per value of a column (category or location).
type ResearchFacet struct {
	Value string `json:"value"`
	Total int64  `json:"total"`
}

func ResearchFacets(column string, statuses []string) ([]ResearchFacet, error) {
	if column != "category" && column != "location" {
		return nil, errors.New("unsupported facet")
	}
	rows := make([]ResearchFacet, 0)
	db := config.DB.Model(&models.Research{}).Where("deleted_at IS NULL")
	if len(statuses) > 0 {
		db = db.Where("status IN ?", statuses)
	}
	err := db.Select(column + " AS value, COUNT(*) AS total").Group(column).Order("total DESC").Scan(&rows).Error
	return rows, err
}

func FetchResearchBySlug(slug string) (*models.Research, error) {
	var r models.Research
	err := config.DB.Preload("Uploader").Preload("Reviewer").
		Where("slug = ? AND deleted_at IS NULL", slug).First(&r).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func FetchResearchComments(researchID int64) ([]models.SubmissionComment, error) {
	var list []models.SubmissionComment
	err := config.DB.Preload("Author").Where("submission_internal_id = ?", researchID).
		Order("created_at ASC").Find(&list).Error
	return list, err
}

func FetchGroupWithMembers(groupID int64) (*models.ResearchGroup, error) {
	var g models.ResearchGroup
	err := config.DB.Preload("Leader").Preload("Members", "left_at IS NULL").Preload("Members.User").
		Where("group_id = ?", groupID).First(&g).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &g, err
}

type DocumentationQuery struct {
	ResearchSlug string
	Category     string
	Offset       int
	Limit        int
}

// PublicDocumentation is a documentation entry joined with its research.
type PublicDocumentation struct {
	models.GroupDocumentation
	ResearchSlug     string `json:"research_slug"`
	ResearchTitle    string `json:"research_title"`
	ResearchCode     string `json:"research_code"`
	ResearchCategory string `json:"research_category"`
	GroupPublicID    string `json:"group_public_id"`
}

// PublicDocumentations lists documentation of on-going/approved researches.
func PublicDocumentations(q DocumentationQuery) ([]PublicDocumentation, int64, error) {
	build := func() *gorm.DB {
		db := config.DB.Table("group_documentations gd").
			Joins("JOIN researches r ON r.group_internal_id = gd.group_internal_id").
			Where("gd.deleted_at IS NULL AND r.deleted_at IS NULL AND r.status IN ?", []string{"on-going", "approved"})
		if q.ResearchSlug != "" {
			db = db.Where("r.slug = ?", q.ResearchSlug)
		}
		if q.Category != "" {
			db = db.Where("r.category = ?", q.Category)
		}
		return db
	}
	var total int64
	if err := build().Count(&total).Error; err != nil {
		return nil, 0, err
	}
	type row struct {
		DocumentationID int64
		ResearchSlug    string
		ResearchTitle   string
		ResearchCode    string
		Category        string
		GroupPublicID   string
	}
	var rows []row
	err := build().Select("gd.documentation_id, r.slug AS research_slug, r.title AS research_title, COALESCE(r.research_code, '') AS research_code, COALESCE(r.category, '') AS category, (SELECT g.public_id::text FROM research_groups g WHERE g.group_id = gd.group_internal_id) AS group_public_id").
		Order("gd.event_date DESC NULLS LAST, gd.created_at DESC").Offset(q.Offset).Limit(q.Limit).Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]PublicDocumentation, 0, len(rows))
	for _, r := range rows {
		var doc models.GroupDocumentation
		if err := config.DB.Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("display_order ASC") }).
			Preload("Author").First(&doc, "documentation_id = ?", r.DocumentationID).Error; err != nil {
			return nil, 0, err
		}
		out = append(out, PublicDocumentation{GroupDocumentation: doc, ResearchSlug: r.ResearchSlug, ResearchTitle: r.ResearchTitle, ResearchCode: r.ResearchCode, ResearchCategory: r.Category, GroupPublicID: r.GroupPublicID})
	}
	return out, total, nil
}
