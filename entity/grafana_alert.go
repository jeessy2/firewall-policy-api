package entity

// GrafanaAlert Grafana Alert struct
type GrafanaAlert struct {
	Status string // 状态
	Alerts []struct {
		Status string
		Labels struct {
			Alertname string
			Instance  string
		}
		Annotations struct {
			Description string
		}
	}
}
