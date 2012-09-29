package main

import(
	"nimbus.io/nimbusapi"
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


func RunSimulation(credentials *nimbusapi.Credentials, config *Config, 
	infoChan UserInfoChan) error {
	infoChan <- &UserInfo{credentials.Name, UserStatusStarted, "", 
		time.Now().UTC()}
	infoChan <- &UserInfo{credentials.Name, UserStatusNormalTermination, "", 
		time.Now().UTC()}
	return nil
}