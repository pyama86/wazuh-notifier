package wazuh_notifier

func (a *Alert) Message() string {
	switch a.Location {
	case "vulnerability-detector":
		return a.Data.Vulnerability.Reference
	default:
		return a.FullLog
	}
	return ""
}

type Alert struct {
	Timestamp string `json:"timestamp"`
	Rule      struct {
		Level       int      `json:"level"`
		Description string   `json:"description"`
		ID          string   `json:"id"`
		Firedtimes  int      `json:"firedtimes"`
		Mail        bool     `json:"mail"`
		Groups      []string `json:"groups"`
		PciDss      []string `json:"pci_dss"`
		Gdpr        []string `json:"gdpr"`
	} `json:"rule"`
	Agent struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Labels struct {
			SlackChannel string `json:"slack_channel"`
		} `json:"labels"`
	} `json:"agent"`
	Manager struct {
		Name string `json:"name"`
	} `json:"manager"`
	ID      string `json:"id"`
	Cluster struct {
		Name string `json:"name"`
		Node string `json:"node"`
	} `json:"cluster"`
	FullLog string `json:"full_log"`
	Decoder struct {
		Parent string `json:"parent"`
		Name   string `json:"name"`
	} `json:"decoder"`
	Data struct {
		Srcip         string `json:"srcip"`
		Srcport       string `json:"srcport"`
		ID            string `json:"id"`
		Vulnerability struct {
			Cve       string `json:"cve"`
			Title     string `json:"title"`
			Severity  string `json:"severity"`
			Published string `json:"published"`
			State     string `json:"state"`
			Package   struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"package"`
			Reference string `json:"reference"`
		} `json:"vulnerability"`
	} `json:"data"`
	Location string `json:"location"`
}
