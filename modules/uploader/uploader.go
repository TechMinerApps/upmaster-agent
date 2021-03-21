package modules

import (
	"sync"

	"github.com/TechMinerApps/upmaster-agent/models"
	"github.com/TechMinerApps/upmaster/router/api/v1/status"
	"github.com/go-resty/resty/v2"
)

// Uploader is a interface needed for uploading status info to UpMaster
type Uploader interface {
	// Write is a non-blocking method that put models.EndpointStatus into upload queue
	Write(s *models.EndpointStatus)

	// Flush is a blocking method that wait for all the models.EndpointStatus in queue to be proceeded.
	Flush() error

	Start() error

	// Stop wait until waitgroup is clear
	Stop() error
}

type Config struct {
	BufferSize   int
	ErrorChannel chan<- error
	RemoteAddr   string
	BatchSize    int
	// Main object
	App App
}

// App is the interface used to communicate with main object
type App interface {
	NodeID() uint
}

type uploader struct {
	config Config
	client *resty.Client
	pool   chan *models.EndpointStatus
	wg     sync.WaitGroup
}

func (u *uploader) Start() error {

	u.wg.Add(1)
	go u.worker()
	return nil
}

func (u *uploader) Stop() error {
	close(u.pool)
	u.wg.Wait()
	return nil
}

func (u *uploader) Write(s *models.EndpointStatus) {
	u.pool <- s
}

func (u *uploader) Flush() error {
	return nil
}

func (u *uploader) upload(statuspoint *models.EndpointStatus) {
	data := &status.WriteEndpointRequest{
		NodeID:     int(u.config.App.NodeID()),
		EndpointID: int(statuspoint.EndpointID),
		Up:         0,
	}
	_, err := u.client.R().
		SetBody(data).
		Post("/status")
	if err != nil {
		u.config.ErrorChannel <- err
	}
}

func (u *uploader) worker() {
	for s := range u.pool {
		u.upload(s)
	}
	u.wg.Done()
}

// NewUploader generate uploader according to configuration
func NewUploader(c *Config) (Uploader, error) {

	bufferPool := make(chan *models.EndpointStatus, c.BufferSize)
	client := resty.New()
	client.SetHostURL(c.RemoteAddr)

	return &uploader{
		// Copy here
		config: *c,
		client: client,
		pool:   bufferPool,
		wg:     sync.WaitGroup{},
	}, nil

}
