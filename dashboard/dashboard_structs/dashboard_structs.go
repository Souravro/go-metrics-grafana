package dashboard_structs

type DashboardConfig struct {
	AppName         string `json:"app_name"`
	DataHost        string `json:"data_host"`
	RequestInterval int64  `json:"request_interval"`
}
