package stats

import "go_backend/api"

type StatManager struct {
	StatList *api.Stats
}

func (m *StatManager) Init() {
	m.StatList = &api.Stats{RouterStatsList: make([]api.RouterStats, 0)}

	m.StatList.RouterStatsList = append(m.StatList.RouterStatsList, api.RouterStats{Route: "event"})
	m.StatList.RouterStatsList = append(m.StatList.RouterStatsList, api.RouterStats{Route: "login"})
	m.StatList.RouterStatsList = append(m.StatList.RouterStatsList, api.RouterStats{Route: "user"})
	m.StatList.RouterStatsList = append(m.StatList.RouterStatsList, api.RouterStats{Route: "sys_user"})
	m.StatList.RouterStatsList = append(m.StatList.RouterStatsList, api.RouterStats{Route: "reports"})

	for i := range m.StatList.RouterStatsList {
		m.StatList.RouterStatsList[i].CreateRequests = new(int64)
		m.StatList.RouterStatsList[i].UpdateRequests = new(int64)
		m.StatList.RouterStatsList[i].GetRequests = new(int64)
		m.StatList.RouterStatsList[i].GetByIdRequests = new(int64)
		m.StatList.RouterStatsList[i].DeleteRequests = new(int64)
		m.StatList.RouterStatsList[i].Errors = new(int64)
		m.StatList.RouterStatsList[i].SuccessAuthRequests = new(int64)
	}
}

func (m *StatManager) GetStat(statRouterName string) *api.RouterStats {
	for _, r := range m.StatList.RouterStatsList {
		if r.Route == statRouterName {
			return &r
		}
	}

	return &api.RouterStats{}
}
