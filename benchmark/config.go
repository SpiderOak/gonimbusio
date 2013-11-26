package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Action int

const (
	CreateBucket Action = iota
	CreateVersionedBucket
	DeleteBucket
	ArchiveNewKey
	ArchiveNewVersion
	ArchiveOverwrite
	RetrieveLatest
	RetrieveVersion
	DeleteKey
	DeleteVersion
)

var actionStringMap = map[Action]string{
	CreateBucket:          "CreateBucket",
	CreateVersionedBucket: "CreateVersionedBucket",
	DeleteBucket:          "DeleteBucket",
	ArchiveNewKey:         "ArchiveNewKey",
	ArchiveNewVersion:     "ArchiveNewVersion",
	ArchiveOverwrite:      "ArchiveOverwrite",
	RetrieveLatest:        "RetrieveLatest",
	RetrieveVersion:       "RetrieveVersion",
	DeleteKey:             "DeleteKey",
	DeleteVersion:         "DeleteVersion",
}

func (action Action) String() string {
	if name, found := actionStringMap[action]; found {
		return name
	}
	return "unknown Action"
}

type ActionDistribution struct {
	CreateBucket          int `json:"create-bucket"`
	CreateVersionedBucket int `json:"create-versioned-bucket"`
	DeleteBucket          int `json:"delete-bucket"`
	ArchiveNewKey         int `json:"archive-new-key"`
	ArchiveNewVersion     int `json:"archive-new-version"`
	ArchiveOverwrite      int `json:"archive-overwrite"`
	RetrieveLatest        int `json:"retrieve-latest"`
	RetrieveVersion       int `json:"retrieve-version"`
	DeleteKey             int `json:"delete-key"`
	DeleteVersion         int `json:"delete-version"`
}

// struct tags are the actual JSON names
type RawConfig struct {
	LowDelay          float32            `json:"low-delay"`
	HighDelay         float32            `json:"high-delay"`
	VerifyBefore      bool               `json:"verify-before"`
	VerifyAfter       bool               `json:"verify-after"`
	Distribution      ActionDistribution `json:"distribution"`
	MaxBucketCount    int                `json:"max-bucket-count"`
	MinFileSize       int                `json:"min-file-size"`
	MaxFileSize       int                `json:"max-file-size"`
	MultipartPartSize int                `json:"multipart-part-size"`
}

type actionValue struct {
	Action Action
	Count  int
}

type Config struct {
	LowDelay          float32
	HighDelay         float32
	VerifyBefore      bool
	VerifyAfter       bool
	MaxBucketCount    int
	MinFileSize       int
	MaxFileSize       int
	MultipartPartSize int

	// an array of 100 entries, each entry is an Action
	// each Action appears the number of times specified in Distribution
	// so if we pick an entry at random, the probability of drawing
	// an individual action is the percentage defined in Distribution
	ActionSlice []Action
}

func LoadConfig(path string) (*Config, error) {
	rawBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rawConfig RawConfig
	err = json.Unmarshal(rawBytes, &rawConfig)
	if err != nil {
		return nil, err
	}
	log.Printf("rawConfig: %v", rawConfig)

	// distribution percentage of each action, should sum to 100
	var actionDistributions = []actionValue{
		{CreateBucket, rawConfig.Distribution.CreateBucket},
		{CreateVersionedBucket, rawConfig.Distribution.CreateVersionedBucket},
		{DeleteBucket, rawConfig.Distribution.DeleteBucket},
		{ArchiveNewKey, rawConfig.Distribution.ArchiveNewKey},
		{ArchiveNewVersion, rawConfig.Distribution.ArchiveNewVersion},
		{ArchiveOverwrite, rawConfig.Distribution.ArchiveOverwrite},
		{RetrieveLatest, rawConfig.Distribution.RetrieveLatest},
		{RetrieveVersion, rawConfig.Distribution.RetrieveVersion},
		{DeleteKey, rawConfig.Distribution.DeleteKey},
		{DeleteVersion, rawConfig.Distribution.DeleteVersion},
	}

	actionPercent := make([]Action, 100)
	index := 0

	for _, actionDistribution := range actionDistributions {
		for i := 0; i < actionDistribution.Count; i++ {
			actionPercent[index] = actionDistribution.Action
			index += 1
		}
	}
	// the distribution should add up to 100%
	if index != 100 {
		return nil, fmt.Errorf("distributin does not total 100%% %s", index)
	}

	config := Config{
		rawConfig.LowDelay,
		rawConfig.HighDelay,
		rawConfig.VerifyBefore,
		rawConfig.VerifyAfter,
		rawConfig.MaxBucketCount,
		rawConfig.MinFileSize,
		rawConfig.MaxFileSize,
		rawConfig.MultipartPartSize,
		actionPercent,
	}
	return &config, nil
}
