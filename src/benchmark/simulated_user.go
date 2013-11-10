package main

import(
	"fmt"
	"nimbus.io/nimbusapi"
	"math/rand"
	"time"
) 

type UserStatus int 

const (
	UserStatusStarted UserStatus = iota
	UserStatusInfo
	UserStatusError
	UserStatusAbort
	UserStatusNormalTermination
)

var userStatusMap = map[UserStatus]string {
	UserStatusStarted : "UserStatusStarted",
	UserStatusInfo : "UserStatusInfo",
	UserStatusError : "UserStatusError",
	UserStatusAbort : "UserStatusAbort",
	UserStatusNormalTermination : "UserStatusNormlTermination",
}

func (userStatus UserStatus) String() string {
	name, found := userStatusMap[userStatus]; if found {
		return name
	}
	return "unknown UserStatus"
}

type UserInfo struct {
	UserName string
	Status UserStatus
	Text string
	Timestamp time.Time
}

type UserInfoChan chan<- *UserInfo

type ActionFunction func(nimbusapi.Requester, *Config) (string, error)

var CreateBucketFunction ActionFunction = func(requester nimbusapi.Requester, 
	config *Config) (string, error) {
	return "stub", nil
}
var CreateVersionedBucketFunction ActionFunction = func(
	requester nimbusapi.Requester, config *Config) (string, error) {
	return "stub", nil
}
var DeleteBucketFunction ActionFunction = func(requester nimbusapi.Requester, 
	config *Config) (string, error) {
	return "stub", nil
}
var ArchiveNewKeyFunction ActionFunction = func(requester nimbusapi.Requester, 
	config *Config) (string, error) {
	return "stub", nil
}
var ArchiveNewVersionFunction ActionFunction = func(
	requester nimbusapi.Requester, config *Config) (string, error) {
	return "stub", nil
}
var ArchiveOverwriteFunction ActionFunction = func(
	requester nimbusapi.Requester, config *Config) (string, error) {
	return "stub", nil
}
var RetrieveLatestFunction ActionFunction = func(requester nimbusapi.Requester, 
	config *Config) (string, error) {
	return "stub", nil
}
var RetrieveVersionFunction ActionFunction = func(requester nimbusapi.Requester, 
	config *Config) (string, error) {
	return "stub", nil
}
var DeleteKeyFunction ActionFunction = func(requester nimbusapi.Requester, 
	config *Config) (string, error) {
	return "stub", nil
}
var DeleteVersionFunction ActionFunction = func(requester nimbusapi.Requester, 
	config *Config) (string, error) {
	return "stub", nil
}

var ActionFunctionMap = map[Action]ActionFunction {
	CreateBucket : CreateBucketFunction,
	CreateVersionedBucket : CreateVersionedBucketFunction,
    DeleteBucket : DeleteBucketFunction,
    ArchiveNewKey : ArchiveNewKeyFunction,
    ArchiveNewVersion : ArchiveNewVersionFunction,
    ArchiveOverwrite : ArchiveOverwriteFunction,
    RetrieveLatest : RetrieveLatestFunction,
    RetrieveVersion : RetrieveVersionFunction,
    DeleteKey : DeleteKeyFunction,
    DeleteVersion : DeleteVersionFunction,
}

func RunSimulation(credentials *nimbusapi.Credentials, 
	requester nimbusapi.Requester, config *Config, finishTime time.Time, 
	infoChan UserInfoChan) error {
	var currentTime = time.Now().UTC()
	infoChan <- &UserInfo{credentials.Name, UserStatusStarted, "", 
		currentTime}
	for ; currentTime.Before(finishTime); currentTime = time.Now().UTC() {
		action := chooseRandomAction(config)
		result, err := ActionFunctionMap[action](requester, config) 
		if err != nil {
			errorReport := fmt.Sprintf("%s %s", action, err)
			infoChan <- &UserInfo{credentials.Name, UserStatusError,
				errorReport, currentTime}
		} else {
			actionReport := fmt.Sprintf("%s %s", action, result)
			infoChan <- &UserInfo{credentials.Name, UserStatusInfo,
				actionReport, currentTime}
		} 
		time.Sleep(computeRandomDelay(config))
	}
	infoChan <- &UserInfo{credentials.Name, UserStatusNormalTermination, "", 
		currentTime}
	return nil
}

func computeRandomDelay(config * Config) time.Duration {
	interval := config.HighDelay - config.LowDelay
	delay := config.LowDelay +  interval * rand.Float32()  
	return time.Second * time.Duration(int(delay))
}

func chooseRandomAction(config * Config) Action {
	index := rand.Int() % 100
	return config.ActionSlice[index]

}
