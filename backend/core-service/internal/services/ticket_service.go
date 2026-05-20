package services

import (
	"context"
	"strings"

	"core-service/internal/dao"
	"core-service/internal/domain/entities"
	"core-service/internal/domain/errors"
	"core-service/internal/domain/repositories"
)

type TicketService struct {
	ticketRepo       *repositories.TicketRepository
	specialistRepo   *repositories.SpecialistRepository
	analysisAccessor *dao.AnalysisAccessor
}

func NewTicketService(
	ticketRepo *repositories.TicketRepository,
	specialistRepo *repositories.SpecialistRepository,
	analysisAccessor *dao.AnalysisAccessor,
) *TicketService {
	return &TicketService{
		ticketRepo:       ticketRepo,
		specialistRepo:   specialistRepo,
		analysisAccessor: analysisAccessor,
	}
}

func (s *TicketService) CreateTicket(
	ctx context.Context,
	userID int64,
	message string,
) (*entities.Ticket, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil, errors.ErrEmptyMessage
	}

	category := "SUPPORT"

	if s.analysisAccessor != nil {
		res, err := s.analysisAccessor.ClassifyMessage(ctx, message)
		if err == nil && res != nil {
			res.Category = strings.TrimSpace(res.Category)
			if res.Category != "" {
				category = res.Category
			}
		}
	}

	ticket := &entities.Ticket{
		UserId:   userID,
		Message:  message,
		Category: category,
		Status:   entities.StatusNew,
	}

	if err := s.ticketRepo.Create(ctx, ticket); err != nil {
		return nil, err
	}

	assignedTicket, err := s.assignFreeSpec(ctx, ticket)
	if err != nil {
		return ticket, nil
	}

	return assignedTicket, nil
}

func (s *TicketService) GetTicket(ctx context.Context, id int64) (*entities.Ticket, error) {
	return s.ticketRepo.GetById(ctx, id)
}

func (s *TicketService) ListTickets(ctx context.Context, filters *entities.Filter) ([]*entities.Ticket, error) {
	return s.ticketRepo.List(ctx, filters)
}

func (s *TicketService) assignFreeSpec(
	ctx context.Context,
	ticket *entities.Ticket,
) (*entities.Ticket, error) {
	specialists, err := s.specialistRepo.FindByCategory(ctx, ticket.Category)
	if err != nil {
		return nil, err
	}

	if len(specialists) == 0 {
		return nil, errors.ErrNotFound
	}

	selected := specialists[0]
	for _, sp := range specialists[1:] {
		if sp.CurrentLoad < selected.CurrentLoad {
			selected = sp
		}
	}

	ticket.SpecialistId = &selected.UserId
	ticket.Status = entities.StatusAssigned

	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	if err := s.specialistRepo.IncrementLoad(ctx, selected.UserId, ticket.Category); err != nil {
		return nil, err
	}

	selected.CurrentLoad++

	return ticket, nil
}

func (s *TicketService) CloseTicket(ctx context.Context, ticketID int64) (*entities.Ticket, error) {
	ticket, err := s.ticketRepo.GetById(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	if ticket.SpecialistId != nil {
		if err := s.specialistRepo.DecrementLoad(ctx, *ticket.SpecialistId, ticket.Category); err != nil {
			return nil, err
		}
	}

	ticket.Status = entities.StatusClosed

	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *TicketService) AssignTicket(
	ctx context.Context,
	ticketID int64,
	specialistID int64,
) (*entities.Ticket, error) {
	ticket, err := s.ticketRepo.GetById(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	specialist, err := s.specialistRepo.GetById(ctx, specialistID, ticket.Category)
	if err != nil {
		return nil, err
	}
	if specialist == nil {
		return nil, errors.ErrNotFound
	}

	if ticket.SpecialistId != nil && *ticket.SpecialistId == specialist.UserId {
		return ticket, nil
	}

	if ticket.SpecialistId != nil {
		if err := s.specialistRepo.DecrementLoad(ctx, *ticket.SpecialistId, ticket.Category); err != nil {
			return nil, err
		}
	}

	ticket.SpecialistId = &specialist.UserId
	ticket.Status = entities.StatusAssigned

	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	if err := s.specialistRepo.IncrementLoad(ctx, specialist.UserId, ticket.Category); err != nil {
		return nil, err
	}

	return ticket, nil
}
