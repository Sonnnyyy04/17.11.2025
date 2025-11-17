package dto

type LinksRequest struct {
	Links []string `json:"links"`
}

type LinksResponse struct{
	LinksStatuses map[string]string `json:"links"`
	LinksNum int `json:"links_num"`
}

type ReportRequest struct {
	LinksList []int `json:"links_list"`
}