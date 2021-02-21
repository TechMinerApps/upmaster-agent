package modules

import (
	"sync"

	"github.com/go-resty/resty/v2"
)

// Uploader is a interface needed for uploading status info to UpMaster
type Uploader interface {
	// Write is a non-blocking method that put StatusPoint into upload queue
	Write(s *StatusPoint)

	// Flush is a blocking method that wait for all the StatusPoint in queue to be proceeded.
	Flush() error

	Start() error
	Stop() error
}

type StatusPoint struct {
}

type Config struct {
	BufferSize   int
	ErrorChannel chan<- error
	RemoteAddr   string
}

type uploader struct {
	config       Config
	client       *resty.Client
	pool         chan *StatusPoint
	bufferSize   int
	remoteAddr   string
	wg           sync.WaitGroup
	errorChannel chan<- error
}

func (u *uploader) Start() error {

	u.wg.Add(1)
	go u.worker()
	return nil
}

func (u *uploader) Stop() error {
	close(u.pool)
	return nil
}

func (u *uploader) Write(s *StatusPoint) {
	u.pool <- s
}

func (u *uploader) Flush() error {
	return nil
}

func (u *uploader) upload(statuspoint *StatusPoint) {
	_, err := u.client.R().
		SetBody(statuspoint).
		Post("/status")
	if err != nil {
		u.errorChannel <- err
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

	bufferPool := make(chan *StatusPoint, c.BufferSize)
	client := resty.New()
	client.SetHostURL(c.RemoteAddr)

	return &uploader{
		config:       *c,
		client:       client,
		pool:         bufferPool,
		bufferSize:   c.BufferSize,
		remoteAddr:   c.RemoteAddr,
		wg:           sync.WaitGroup{},
		errorChannel: make(chan<- error),
	}, nil

}
