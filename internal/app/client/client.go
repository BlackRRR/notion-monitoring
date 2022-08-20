package client

import (
	"context"
	"github.com/BlackRRR/notion-monitoring/internal/app/service"
	"github.com/jomei/notionapi"
)

type Client struct {
	*notionapi.Client
	*service.Service
}

func NewClient(service *service.Service, token string) *Client {
	client := &Client{
		Client:  notionapi.NewClient(notionapi.Token(token)),
		Service: service,
	}

	return client
}

func (c *Client) StartClient(ctx context.Context) error {
	//берем задачи из ноушена
	resp, err := c.Service.GetNotionDatabaseQuery(ctx, c.Client)
	if err != nil {
		return err
	}

	//распределяем отправку задач
	err = c.Service.TaskDistribution(resp)
	if err != nil {
		return err
	}

	//Сохраняем в базе данных после отправки соо
	err = c.Service.AddToCacheOrDBFromNotionResp(resp, nil, service.DB)
	if err != nil {
		return err
	}

	//Сохраняем в кэше после отправки соо
	err = c.Service.AddToCacheOrDBFromNotionResp(resp, nil, service.Cache)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) StartPages(ctx context.Context) error {
	resp, err := c.Service.GetNotionDatabaseQuery(ctx, c.Client)
	if err != nil {
		return err
	}

	pages, err := c.GetPagesFromDB()
	if err != nil {
		return err
	}

	if pages == nil {
		err := c.AddToCacheOrDBFromNotionResp(resp, nil, service.Cache)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.AddToCacheOrDBFromNotionResp(nil, pages, service.Cache)
	if err != nil {
		return err
	}

	return nil
}
