package splunk

import (
	"github.com/Sirupsen/logrus"
	log "github.com/TykTechnologies/logrus"
	"github.com/TykTechnologies/tyk-pump/analytics"
	"github.com/mitchellh/mapstructure"
)

// Pump is a Tyk Pump driver for Splunk.
type Pump struct {
	client *Client
	config *PumpConfig
}

// PumpConfig contains the driver configuration parameters.
type PumpConfig struct {
	Token         string
	Endpoint      string
	TLSSkipVerify bool
}

const (
	pumpPrefix = "splunk-pump"
	pumpName   = "Splunk Pump"
)

// New initializes a new pump.
func (p *Pump) New() Pump {
	return Pump{}
}

// GetName returns the pump name.
func (p *Pump) GetName() string {
	return pumpName
}

// Init performs the initialization.
func (p *Pump) Init(config interface{}) error {
	p.config = &PumpConfig{}
	err := mapstructure.Decode(config, &p.config)
	if err != nil {
		return err
	}
	log.WithFields(logrus.Fields{
		"prefix": pumpPrefix,
	}).Infof("%s Endpoint: %s", pumpName, p.config.Endpoint)

	p.client, err = New(p.config.Token, p.config.Endpoint, p.config.TLSSkipVerify)
	if err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"prefix": pumpPrefix,
	}).Debugf("%s Initialized", pumpName)
	return nil
}

func (p *Pump) WriteData(data []interface{}) error {
	log.WithFields(logrus.Fields{
		"prefix": pumpPrefix,
	}).Info("Writing ", len(data), " records")
	for _, v := range data {
		decoded := v.(analytics.AnalyticsRecord)

		event := map[string]interface{}{
			"api_id":        decoded.APIID,
			"path":          decoded.Path,
			"method":        decoded.Method,
			"response_code": decoded.ResponseCode,
		}
		p.client.Send(event)
	}
	return nil
}
