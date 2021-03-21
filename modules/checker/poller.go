package checker

import (
	"net/http"
	"time"

	agentModel "github.com/TechMinerApps/upmaster-agent/models"
	"github.com/TechMinerApps/upmaster/models"
)

type Poller struct {
	endpoint      models.Endpoint
	ticker        *time.Ticker
	statusChannel chan<- *agentModel.EndpointStatus
	client        *http.Client
}

func (p *Poller) Start() {

	ticker := time.NewTicker(time.Second * time.Duration(p.endpoint.Interval))
	go func() {
		for range ticker.C {
			p.worker()
		}
	}()
	p.ticker = ticker
	return
}

func (p *Poller) worker() {
	req, err := http.NewRequest("GET", p.endpoint.URL, nil)
	if err != nil {
		return
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	p.statusChannel <- &agentModel.EndpointStatus{
		EndpointID: p.endpoint.ID,
		TimeStamp:  time.Time{},
		Status:     agentModel.UP,
		Error:      0,
		ErrorCode:  resp.StatusCode,
	}
	return
}
