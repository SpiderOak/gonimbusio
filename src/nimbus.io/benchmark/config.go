package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "time"
)

const (
    CreateBucket = iota
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

type TestDistribution struct {
    CreateBucket int `json:"create-bucket"`
    CreateVersionedBucket int `json:"create-versioned-bucket"`
    DeleteBucket int `json:"delete-bucket"`
    ArchiveNewKey int `json:"archive-new-key"`
    ArchiveNewVersion int `json:"archive-new-version"`
    ArchiveOverwrite int `json:"archive-overwrite"`
    RetrieveLatest int `json:"retrieve-latest"`
    RetrieveVersion int `json:"retrieve-version"`
    DeleteKey int `json:"delete-key"`
    DeleteVersion int `json:"delete-version"`
}

// struct tags are the actual JSON names
type RawConfig struct {
   LowDelay float32 `json:"low-delay"`
   HighDelay float32 `json:"high-delay"`
   VerifyBefore bool `json:"verify-before"`
   VerifyAfter bool `json:"verify-after"`
   Distribution TestDistribution `json:"distribution"`
   MaxBucketCount int `json:"max-bucket-count"`
   MinFileSize int `json:"min-file-size"`
   MaxFileSize int `json:"max-file-size"`
   MultipartPartSize int `json:"multipart-part-size"`
}

type Config struct {
  LowDelay time.Duration
  HighDelay time.Duration
  VerifyBefore bool
  VerifyAfter bool
  MaxBucketCount int
  MinFileSize int
  MaxFileSize int
  MultipartPartSize int
}

func LoadConfig(path string) (*Config, error) {
  rawBytes, err := ioutil.ReadFile(path); if err != nil {
    return nil, err
  }
  var rawConfig RawConfig
  err = json.Unmarshal(rawBytes, &rawConfig); if err != nil {
    return nil, err
  }
  log.Printf("rawConfig: %v", rawConfig)
  config := Config{
    time.Duration(int(rawConfig.LowDelay)) * time.Second,
    time.Duration(int(rawConfig.HighDelay)) * time.Second,
    rawConfig.VerifyBefore,
    rawConfig.VerifyAfter,
    rawConfig.MaxBucketCount,
    rawConfig.MinFileSize,
    rawConfig.MaxFileSize,
    rawConfig.MultipartPartSize,
  }
  return &config, nil
}

