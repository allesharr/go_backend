package api

type Stats struct {
	RouterStatsList []RouterStats `json:"routes"`
}

type RouterStats struct {
	Route               string `json:"route"`
	GetRequests         *int64 `json:"get_requests"`
	GetByIdRequests     *int64 `json:"get_by_id_requests"`
	CreateRequests      *int64 `json:"create_requests"`
	UpdateRequests      *int64 `json:"update_requests"`
	DeleteRequests      *int64 `json:"delete_requests"`
	SuccessAuthRequests *int64 `json:"success_auth_requests"`
	Errors              *int64 `json:"errors"`
}
