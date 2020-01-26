package dbclient

import (
	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
)

// DBClient defines public methods to handle QoE related information in DB
type DBClient interface {
	GetQoE(*pbapi.RequestQoE, chan pbapi.ResponseQoE)
}

type dbClient struct {
}

func (db *dbClient) GetQoE(reqQoes *pbapi.RequestQoE, result chan pbapi.ResponseQoE) {
	replQoes := pbapi.ResponseQoE{
		Qoes: reqQoes.Qoes,
	}
	// Simulating hung process
	// time.Sleep(maxRequestProcessTime)

	result <- replQoes
}

// NewDBClient return  a new instance of a DB client
func NewDBClient() DBClient {

	return &dbClient{}
}
