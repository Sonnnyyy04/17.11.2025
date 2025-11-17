package service

import (
	"links_project/internal/models"
	"links_project/internal/storage"
	"net/http"
	"strings"
)

type Service interface {
	CreateBatch(link []string)(*models.Batch, error)
	GetBatch(id []int)([]*models.Batch, error)
}

type service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) Service{
	return &service{
		storage: storage,
	}
}

func (s *service) CreateBatch(links []string)(*models.Batch, error){
	id := s.storage.NextID()
	batch := &models.Batch{
		ID: id,
		Links: links,
		Statuses: make(map[string]string),
	}

	for _, link := range links{
		if checkLink(link){
			batch.Statuses[link] = "available"
		}else{
			batch.Statuses[link] = "not available"
		}
	}

	s.storage.SaveBatch(batch)
	if err := s.storage.Save(); err != nil {
		return nil, err
	}

	return batch, nil
}

func (s *service) GetBatch(id []int)([]*models.Batch, error){
	var res []*models.Batch

	for _, i := range id{
		batch, ok := s.storage.GetBatch(i)
		if ok{
			res = append(res, batch)
		}
	}
	return res, nil
}

func checkLink(link string)bool{
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	resp, err := http.Get(link)
	if err != nil{
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK{
		return true
	}
	return false
}