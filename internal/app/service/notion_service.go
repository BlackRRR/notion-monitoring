package service

import (
	"context"
	"fmt"
	"github.com/BlackRRR/notion-monitoring/internal/app/bot"
	"github.com/BlackRRR/notion-monitoring/internal/app/repository"
	"github.com/BlackRRR/notion-monitoring/internal/cfg"
	"github.com/BlackRRR/notion-monitoring/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jomei/notionapi"
)

const (
	DB    = "DB"
	Cache = "Cache"
)

type Service struct {
	*repository.Repository
	*repository.StatusCache
	*repository.DescriptionCache
	*tgbotapi.BotAPI
}

func NewService(rep *repository.Repository, bot *tgbotapi.BotAPI) *Service {
	service := &Service{
		rep,
		repository.NewStatusCache(),
		repository.NewDescriptionCache(),
		bot,
	}

	return service
}

func (s *Service) GetNotionDatabaseQuery(ctx context.Context, client *notionapi.Client) (*notionapi.DatabaseQueryResponse, error) {
	database, err := client.Database.Query(ctx, notionapi.DatabaseID(cfg.NotionDBCOnfig), &notionapi.DatabaseQueryRequest{})
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (s *Service) AddToCacheOrDB(ID, status, description, repository string) error {
	if repository == "Cache" {
		s.StatusCache.Add(ID, status)
		s.DescriptionCache.Add(ID, description)
	} else {
		err := s.CreateOrUpdateNotionPages(ID, status, description)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) AddToCacheOrDBFromNotionResp(resp *notionapi.DatabaseQueryResponse, pages []model.Page, repo string) error {
	if resp != nil {
		for pageNum := range resp.Results {
			err := s.AddToCacheOrDB(resp.Results[pageNum].ID.String(),
				resp.Results[pageNum].Properties["Status"].(*notionapi.SelectProperty).Select.Name,
				resp.Results[pageNum].Properties["Name"].(*notionapi.TitleProperty).Title[0].PlainText,
				repo,
			)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for pageNum := range pages {
		err := s.AddToCacheOrDB(pages[pageNum].ID,
			pages[pageNum].Status,
			pages[pageNum].Description,
			repo,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) TaskDistribution(resp *notionapi.DatabaseQueryResponse) error {
	var notionResp *model.NotionResponse
	for pageNum := range resp.Results {
		status := resp.Results[pageNum].Properties["Status"].(*notionapi.SelectProperty)
		description := resp.Results[pageNum].Properties["Name"].(*notionapi.TitleProperty)
		notionResp = &model.NotionResponse{
			ID:          resp.Results[pageNum].ID.String(),
			PageNum:     pageNum,
			Description: description.Title[0].PlainText,
			Status:      status.Select.Name,
			URL:         resp.Results[pageNum].URL}

		// если статус пустой значит в ноушнионе это No status
		if notionResp.Status == "" {
			if s.DescriptionCache.Get(notionResp.ID) == "" {
				s.NewTaskMessageSend(
					notionResp.URL,
					"Task was added",
					"No status",
					notionResp.Description)
			}
		} else {
			s.NotNilStatus(notionResp)
		}
	}

	return nil
}

func (s *Service) NotNilStatus(notionResp *model.NotionResponse) {
	//если статуса по айдишнику нет то мы отправляем что добавили новую задачу
	if _, ok := s.TaskStatus[notionResp.ID]; !ok {
		s.NewTaskMessageSend(
			notionResp.URL,
			"Task was added",
			notionResp.Status,
			notionResp.Description)
	} else {
		//если такая задача есть но статус другой то мы отправляем что поменяли статус задачи
		if notionResp.Status != s.StatusCache.Get(notionResp.ID) {
			s.NewTaskMessageSend(
				notionResp.URL,
				"Task Status was changed",
				notionResp.Status,
				notionResp.Description)
		}
	}
}

func (s *Service) GetPagesFromDB() ([]model.Page, error) {
	pages, err := s.GetPages()
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func (s *Service) NewTaskMessageSend(URL, task, status, desc string) {
	text := fmt.Sprintf("%s\nStatus = %s\nDescription = %s\nURL = %s",
		task,
		status,
		desc,
		URL)

	for id := range model.AdminSettings.AdminID {
		msg := bot.CreateNewTGMessage(text, id)
		bot.SendMsgBot(s.BotAPI, msg)
	}
}
