package services

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"gorm.io/gorm"
)

type NotificationService interface {
	Create(notification *models.Notification) error
	NotifPiles(RecipientID int64, filter, sort string, limit, offset int) ([]models.Notification, int64, error)
	Dispatch(event models.NotificationEvent) error
	NotifyGroupMembers(groupID int64, exceptUserID int64, event models.NotificationEvent) error

	ReviewSubmissionNotification(research *models.Research, author *models.User) error
	SubmissionRevisionNotificaation(research *models.Research, author *models.User) error
	SubmissionReviewNotfication(research *models.Research, reviewer *models.User, message *models.SubmissionComment) error

	NotifyDecision(research *models.Research, decision models.SubmissionDecision, feedback string) error
	NotifyResearcherReply(research *models.Research, researcher *models.User, comment *models.SubmissionComment) error
	NotifyGroupCreated(group *models.ResearchGroup) error
	NotifyMemberJoined(group *models.ResearchGroup, newMember *models.User) error
	NotifyMemberRemoved(group *models.ResearchGroup, removedUser *models.User) error
	NotifyMemberLeft(group *models.ResearchGroup, leavingUser *models.User) error
	NotifyDocumentationCreated(doc *models.GroupDocumentation, groupID int64, authorName string) error
	NotifyDocumentationDeleted(doc *models.GroupDocumentation, groupID int64, actorID int64) error
	NotifyDeliverableCreated(groupID int64, actorID int64, items []models.ResearchDeliverable) error
	NotifyDeliverableDone(group *models.ResearchGroup, item *models.ResearchDeliverable, actorName string) error
	NotifyDeliverableReopened(group *models.ResearchGroup, item *models.ResearchDeliverable, actorName string) error
	NotifyDeliverableDeleted(groupID int64, actorID int64, title string) error
}

type NotificationServices struct {
	NotifRepo repositories.NotificationRepository
	db        *gorm.DB
}

func NewNotificationService(NotifRep repositories.NotificationRepository, db *gorm.DB) NotificationService {
	return &NotificationServices{NotifRepo: NotifRep, db: db}
}

func safeTitle(r *models.Research) string {
	if r != nil && r.Title != nil && *r.Title != "" {
		return *r.Title
	}
	return "Untitled"
}
func safeUserName(u *models.User, fallback string) string {
	if u == nil {
		return fallback
	}
	if u.Name != nil && *u.Name != "" {
		return *u.Name
	}
	return fallback
}
func (s *NotificationServices) Create(notification *models.Notification) error {
	if notification.PublicID == uuid.Nil {
		notification.PublicID = uuid.New()
	}
	return s.NotifRepo.Create(notification)
}
func (s *NotificationServices) NotifPiles(RecipientID int64, filter, sort string, limit, offset int) ([]models.Notification, int64, error) {
	return s.NotifRepo.NotificationPiles(RecipientID, filter, sort, limit, offset)
}
func (s *NotificationServices) Dispatch(event models.NotificationEvent) error {
	notif := &models.Notification{
		PublicID:            uuid.New(),
		RecipientInternalID: event.RecipientID,
		Subject:             event.Subject,
		Message:             event.Message,
		Type:                event.Type,
		ReferenceType:       &event.ReferenceType,
		ReferenceInternalID: &event.ReferenceInternal,
		ReferencePublicID:   &event.ReferencePublic,
	}
	return s.NotifRepo.Create(notif)
}
func (s *NotificationServices) NotifyGroupMembers(groupID int64, exceptUserID int64, event models.NotificationEvent) error {
	var members []models.ResearchGroupMember
	if err := s.db.Where("group_internal_id = ?", groupID).Find(&members).Error; err != nil {
		return err
	}
	for _, m := range members {
		if m.UserInternalID == exceptUserID {
			continue
		}
		e := event
		e.RecipientID = m.UserInternalID
		if err := s.Dispatch(e); err != nil {
			log.Printf("[notif] broadcast failed for %d: %v", m.UserInternalID, err)
		}
	}
	return nil
}
func (s *NotificationServices) ReviewSubmissionNotification(research *models.Research, author *models.User) error {
	var reviewers []models.User
	if err := s.db.Where("role = ?", "reviewer").Find(&reviewers).Error; err != nil {
		log.Printf("[notif] failed to fetch reviewers: %v", err)
		return err
	}
	if len(reviewers) == 0 {
		return nil
	}
	title := safeTitle(research)
	authorName := safeUserName(author, "Unknown")
	refType := "research"
	subject := "New Research Submission"
	body := fmt.Sprintf(`%s submitted a new research proposal: "%s". Click the "Pending Research" tab to review it.`, authorName, title)
	for _, reviewer := range reviewers {
		if research.UploaderInternalID != nil && reviewer.InternalID == *research.UploaderInternalID {
			continue
		}
		notif := &models.Notification{
			PublicID:            uuid.New(),
			RecipientInternalID: reviewer.InternalID,
			Subject:             subject,
			Message:             body,
			Type:                "info",
			ReferenceType:       &refType,
			ReferenceInternalID: research.InternalID,
			ReferencePublicID:   &research.PublicID,
		}
		if err := s.NotifRepo.Create(notif); err != nil {
			log.Printf("[notif] create failed for reviewer %d: %v", reviewer.InternalID, err)
		}
	}
	return nil
}
func (s *NotificationServices) SubmissionRevisionNotificaation(research *models.Research, author *models.User) error {
	var reviewers []models.User
	if err := s.db.Where("role = ?", "reviewer").Find(&reviewers).Error; err != nil {
		log.Printf("[notif] failed to fetch reviewers: %v", err)
		return err
	}
	if len(reviewers) == 0 {
		return nil
	}
	title := safeTitle(research)
	authorName := safeUserName(author, "Unknown")
	refType := "research"
	subject := "Research Revision Submitted"
	body := fmt.Sprintf(`%s submitted a revision for: "%s". Click the "Revision" tab to review it.`, authorName, title)
	for _, reviewer := range reviewers {
		if research.UploaderInternalID != nil && reviewer.InternalID == *research.UploaderInternalID {
			continue
		}
		notif := &models.Notification{
			PublicID:            uuid.New(),
			RecipientInternalID: reviewer.InternalID,
			Subject:             subject,
			Message:             body,
			Type:                "info",
			ReferenceType:       &refType,
			ReferenceInternalID: research.InternalID,
			ReferencePublicID:   &research.PublicID,
		}
		if err := s.NotifRepo.Create(notif); err != nil {
			log.Printf("[notif] create failed for reviewer %d: %v", reviewer.InternalID, err)
		}
	}
	return nil
}
func (s *NotificationServices) SubmissionReviewNotfication(research *models.Research, reviewer *models.User, comment *models.SubmissionComment) error {
	if research.UploaderInternalID == nil {
		return nil
	}
	title := safeTitle(research)
	reviewerName := safeUserName(reviewer, "Reviewer")
	messageContent := "No comment provided"
	if comment != nil && comment.Message != "" {
		messageContent = comment.Message
	}
	refType := "research"
	subject := "Feedback from Reviewer"
	body := fmt.Sprintf(`%s has reviewed your proposal "%s". Comment: "%s".`, reviewerName, title, messageContent)
	notif := &models.Notification{
		PublicID:            uuid.New(),
		RecipientInternalID: *research.UploaderInternalID,
		Subject:             subject,
		Message:             body,
		Type:                "info",
		ReferenceType:       &refType,
		ReferenceInternalID: research.InternalID,
		ReferencePublicID:   &research.PublicID,
	}
	return s.NotifRepo.Create(notif)
}
func (s *NotificationServices) NotifyDecision(research *models.Research, decision models.SubmissionDecision, feedback string) error {
	if research.UploaderInternalID == nil {
		return nil
	}
	title := safeTitle(research)
	var subject, message, notifType string
	switch decision {
	case models.DecisionApproved:
		subject = "Proposal Approved 🎉"
		message = fmt.Sprintf(`Your proposal "%s" has been approved.`, title)
		notifType = "success"
	case models.DecisionRejected:
		subject = "Proposal Rejected"
		message = fmt.Sprintf(`Your proposal "%s" was rejected.<br/><br/><strong>Feedback:</strong><br/>%s`, title, feedback)
		notifType = "error"
	case models.DecisionRevision:
		subject = "Revision Requested"
		message = fmt.Sprintf(`Your proposal "%s" requires revision.<br/><br/><strong>Feedback:</strong><br/>%s`, title, feedback)
		notifType = "warning"
	default:
		return nil
	}
	return s.Dispatch(models.NotificationEvent{
		RecipientID:       *research.UploaderInternalID,
		Subject:           subject,
		Message:           message,
		Type:              notifType,
		ReferenceType:     "research",
		ReferenceInternal: *research.InternalID,
		ReferencePublic:   research.PublicID,
	})
}
func (s *NotificationServices) NotifyResearcherReply(research *models.Research, researcher *models.User, comment *models.SubmissionComment) error {
	if research.ReviewerInternalID == nil {
		return nil
	}
	title := safeTitle(research)
	name := safeUserName(researcher, "The researcher")
	content := "No message provided"
	if comment != nil && comment.Message != "" {
		content = comment.Message
	}
	return s.Dispatch(models.NotificationEvent{
		RecipientID:       *research.ReviewerInternalID,
		Subject:           "Researcher Replied",
		Message:           fmt.Sprintf(`%s replied on "%s": "%s"`, name, title, content),
		Type:              "info",
		ReferenceType:     "research",
		ReferenceInternal: *research.InternalID,
		ReferencePublic:   research.PublicID,
	})
}
func (s *NotificationServices) NotifyGroupCreated(group *models.ResearchGroup) error {
	return s.Dispatch(models.NotificationEvent{
		RecipientID: group.LeaderInternalID,
		Subject:     "Research Group Created",
		Message: fmt.Sprintf(
			`Your research group "%s" has been created. Share this code with collaborators: <strong>%s</strong>`,
			group.Name, group.JoinCode,
		),
		Type:              "success",
		ReferenceType:     "group",
		ReferenceInternal: group.InternalID,
		ReferencePublic:   group.PublicID,
	})
}
func (s *NotificationServices) NotifyMemberJoined(group *models.ResearchGroup, newMember *models.User) error {
	name := safeUserName(newMember, "A new member")
	if err := s.Dispatch(models.NotificationEvent{
		RecipientID:       group.LeaderInternalID,
		Subject:           "New Group Member",
		Message:           fmt.Sprintf(`%s joined your research group "%s".`, name, group.Name),
		Type:              "info",
		ReferenceType:     "group",
		ReferenceInternal: group.InternalID,
		ReferencePublic:   group.PublicID,
	}); err != nil {
		log.Printf("[notif] leader notify failed: %v", err)
	}
	if newMember != nil {
		if err := s.Dispatch(models.NotificationEvent{
			RecipientID:       newMember.InternalID,
			Subject:           "Welcome to the Group",
			Message:           fmt.Sprintf(`You joined "%s". You can now view and contribute to the group.`, group.Name),
			Type:              "success",
			ReferenceType:     "group",
			ReferenceInternal: group.InternalID,
			ReferencePublic:   group.PublicID,
		}); err != nil {
			log.Printf("[notif] member notify failed: %v", err)
		}
	}
	return nil
}

func (s *NotificationServices) NotifyMemberRemoved(group *models.ResearchGroup, removedUser *models.User) error {
	if removedUser == nil {
		return nil
	}
	return s.Dispatch(models.NotificationEvent{
		RecipientID:       removedUser.InternalID,
		Subject:           "Removed from Group",
		Message:           fmt.Sprintf(`You have been removed from "%s".`, group.Name),
		Type:              "warning",
		ReferenceType:     "group",
		ReferenceInternal: group.InternalID,
		ReferencePublic:   group.PublicID,
	})
}
func (s *NotificationServices) NotifyMemberLeft(group *models.ResearchGroup, leavingUser *models.User) error {
	name := safeUserName(leavingUser, "A member")
	return s.Dispatch(models.NotificationEvent{
		RecipientID:       group.LeaderInternalID,
		Subject:           "Member Left Group",
		Message:           fmt.Sprintf(`%s left your research group "%s".`, name, group.Name),
		Type:              "info",
		ReferenceType:     "group",
		ReferenceInternal: group.InternalID,
		ReferencePublic:   group.PublicID,
	})
}
func (s *NotificationServices) NotifyDocumentationCreated(doc *models.GroupDocumentation, groupID int64, authorName string) error {
	return s.NotifyGroupMembers(groupID, doc.AuthorInternalID, models.NotificationEvent{
		Subject:           "New Documentation",
		Message:           fmt.Sprintf(`%s posted a new documentation: "%s".`, authorName, doc.Title),
		Type:              "info",
		ReferenceType:     "documentation",
		ReferenceInternal: doc.InternalID,
		ReferencePublic:   doc.PublicID,
	})
}
func (s *NotificationServices) NotifyDocumentationDeleted(doc *models.GroupDocumentation, groupID int64, actorID int64) error {
	return s.NotifyGroupMembers(groupID, actorID, models.NotificationEvent{
		Subject:           "Documentation Removed",
		Message:           fmt.Sprintf(`Documentation "%s" was removed.`, doc.Title),
		Type:              "warning",
		ReferenceType:     "documentation",
		ReferenceInternal: doc.InternalID,
		ReferencePublic:   doc.PublicID,
	})
}
func (s *NotificationServices) NotifyDeliverableCreated(groupID int64, actorID int64, items []models.ResearchDeliverable) error {
	if len(items) == 0 {
		return nil
	}
	var subject, message string
	if len(items) == 1 {
		subject = "New Deliverable"
		message = fmt.Sprintf(`A new deliverable was added: "%s".`, items[0].Title)
	} else {
		subject = "New Deliverables"
		message = fmt.Sprintf(`%d new deliverables were added to your group.`, len(items))
	}
	return s.NotifyGroupMembers(groupID, actorID, models.NotificationEvent{
		Subject:           subject,
		Message:           message,
		Type:              "info",
		ReferenceType:     "group",
		ReferenceInternal: groupID,
		ReferencePublic:   uuid.Nil,
	})
}
func (s *NotificationServices) NotifyDeliverableDone(group *models.ResearchGroup, item *models.ResearchDeliverable, actorName string) error {
	if item.CompletedByInternalID != nil && group.LeaderInternalID == *item.CompletedByInternalID {
		return nil
	}
	return s.Dispatch(models.NotificationEvent{
		RecipientID:       group.LeaderInternalID,
		Subject:           "Deliverable Completed",
		Message:           fmt.Sprintf(`%s completed "%s".`, actorName, item.Title),
		Type:              "success",
		ReferenceType:     "deliverable",
		ReferenceInternal: item.InternalID,
		ReferencePublic:   item.PublicID,
	})
}
func (s *NotificationServices) NotifyDeliverableReopened(group *models.ResearchGroup, item *models.ResearchDeliverable, actorName string) error {
	return s.Dispatch(models.NotificationEvent{
		RecipientID:       group.LeaderInternalID,
		Subject:           "Deliverable Reopened",
		Message:           fmt.Sprintf(`%s reopened "%s".`, actorName, item.Title),
		Type:              "warning",
		ReferenceType:     "deliverable",
		ReferenceInternal: item.InternalID,
		ReferencePublic:   item.PublicID,
	})
}
func (s *NotificationServices) NotifyDeliverableDeleted(groupID int64, actorID int64, title string) error {
	return s.NotifyGroupMembers(groupID, actorID, models.NotificationEvent{
		Subject:           "Deliverable Removed",
		Message:           fmt.Sprintf(`Deliverable "%s" was removed.`, title),
		Type:              "warning",
		ReferenceType:     "group",
		ReferenceInternal: groupID,
		ReferencePublic:   uuid.Nil,
	})
}
