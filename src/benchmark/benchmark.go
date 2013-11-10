package main

import (
	"io/ioutil"
	"log"
	"nimbus.io/nimbusapi"
	"path"
	"time"
)

func main() {
	log.Println("program starts")
	flags, err := loadFlags()
	if err != nil {
		log.Fatalf("Unable to load flags: %s", err)
	}
	log.Printf("flags = %v", flags)

	config, err := LoadConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalf("Unable to load config from %s %s", flags.ConfigPath, err)
	}
	log.Printf("config = %v", config)

	dirSlice, err := ioutil.ReadDir(flags.UserIdentityDir)
	if err != nil {
		log.Fatalf("Unable to read user identity dir %s %s",
			flags.UserIdentityDir, err)
	}

	testDuration := time.Duration(flags.Duration) * time.Second
	finishTime := time.Now().UTC().Add(testDuration)
	infoChan := make(chan *UserInfo, flags.MaxUsers)

	var simCount int = 0
	for _, fileInfo := range dirSlice {
		if simCount == flags.MaxUsers {
			log.Printf("Max simulated users created %v\n", simCount)
			break
		}
		log.Printf("%s\n", fileInfo.Name())
		credentialsPath := path.Join(flags.UserIdentityDir, fileInfo.Name())
		credentials, err := nimbusapi.LoadCredentialsFromPath(credentialsPath)
		if err != nil {
			log.Fatalf("Unable to load credentials %s %s",
				credentialsPath, err)
		}
		requester, err := nimbusapi.NewRequester(credentials)
		if err != nil {
			log.Fatalf("NewRequester failed %s\n", err)
		}
		go RunSimulation(credentials, requester, config, finishTime, infoChan)
		simCount += 1
	}

	for completedCount := 0; completedCount < simCount; {
		select {
		case userInfo := <-infoChan:
			log.Printf("%s %s %s", userInfo.UserName, userInfo.Status,
				userInfo.Text)
			switch userInfo.Status {
			case UserStatusStarted:
			case UserStatusInfo:
			case UserStatusError:
			case UserStatusAbort:
				completedCount += 1
			case UserStatusNormalTermination:
				completedCount += 1
			}
		}
	}

	log.Println("program ends")
}
