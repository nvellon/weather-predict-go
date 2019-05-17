package solarsystem

import (
	"context"

	"cloud.google.com/go/datastore"
)

type (
	Repo interface {
		Save(record *Record) error
		GetByDay(day int) (*Record, error)
		SaveCounter(counter *Counter) error
		GetCounter() (*Counter, error)
	}

	DataStoreRepo struct {
		client *datastore.Client
		ctx    context.Context
	}

	Record struct {
		Day                          int
		IsDrought                    bool
		IsOptimumTemperaturePressure bool
		IsRainSeason                 bool
		FerengiLocation              string
		BetasoideLocation            string
		VulcanoLocation              string
	}

	Counter struct {
		Days                            int
		CountOther                      int
		CountDrought                    int
		CountOptimumTemperaturePressure int
		CountRainSeason                 int
	}
)

const (
	RecordKind  = "WeatherRecord"
	CounterKind = "Counter"
)

func NewRepo(projectID string) (Repo, error) {
	r := DataStoreRepo{}

	r.ctx = context.Background()
	client, err := datastore.NewClient(r.ctx, projectID)
	if err != nil {
		return nil, err
	}
	r.client = client

	return r, nil
}

func (dsr DataStoreRepo) Save(record *Record) error {
	key := datastore.IncompleteKey(RecordKind, nil)
	_, err := dsr.client.Put(dsr.ctx, key, record)

	return err
}

func (dsr DataStoreRepo) GetByDay(day int) (*Record, error) {
	var records []*Record

	query := datastore.NewQuery(RecordKind).Filter("Day =", day).Limit(1)
	_, err := dsr.client.GetAll(dsr.ctx, query, &records)
	if err != nil {
		return nil, err
	}

	return records[0], nil
}

func (dsr DataStoreRepo) SaveCounter(counter *Counter) error {
	key := datastore.IncompleteKey(CounterKind, nil)
	_, err := dsr.client.Put(dsr.ctx, key, counter)

	return err
}

func (dsr DataStoreRepo) GetCounter() (*Counter, error) {
	var counters []*Counter

	query := datastore.NewQuery(CounterKind).Limit(1)
	_, err := dsr.client.GetAll(dsr.ctx, query, &counters)
	if err != nil {
		return nil, err
	}

	return counters[0], nil
}
