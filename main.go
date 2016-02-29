package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	dockerImage := os.Getenv("DOCKER_IMAGE")
	authSid := os.Getenv("RANCHER_AUTH_SID")
	authKey := os.Getenv("RANCHER_AUTH_KEY")
	baseUrl := os.Getenv("RANCHER_BASE_URL")

	jsonString := fmt.Sprintf(`{"inServiceStrategy": { "launchConfig": { "imageUuid":"docker:%s" }},"toServiceStrategy":null}`, dockerImage)
	client := &http.Client{}

	bodyBytes, err := sendRequest(client, fmt.Sprintf("%s/?action=upgrade", baseUrl), authSid, authKey, jsonString)
	if err != nil {
		log.Println("[Upgrade Segment] Error on upgrade : ", err.Error())
		return
	}

	upgradeCompleted := false
	triedToUpgrade := false
	if strings.Contains(string(bodyBytes), `"state":"upgrading"`) {
		fmt.Println("[Finish Upgrade Segment] Upgrade request went through!")

		triedToUpgrade = true
		for i := 0; i < 10; i++ {
			fmt.Println("[Finish Upgrade Segment] Sleeping for 30 seconds for upgrade to finish...")
			time.Sleep(time.Second * 30)
			fmt.Println("[Finish Upgrade Segment] After sleeping, attempting to upgrade!")

			bodyBytes, err := sendRequest(client, fmt.Sprintf("%s/?action=finishupgrade", baseUrl), authSid, authKey, jsonString)
			if err != nil {
				log.Println("[Finish Upgrade Segment] Error on finish upgrade : ", err.Error())
				continue
			}

			if !strings.Contains(string(bodyBytes), `"type":"error"`) {
				fmt.Println("[Finish Upgrade Segment] Finish upgrade has completed")
				upgradeCompleted = true
				break
			} else {
				fmt.Println("[Finish Upgrade Segment] Error on finishing upgrade, trying again..")
			}
		}
	}

	if triedToUpgrade && !upgradeCompleted {
		fmt.Println("[Rollback Segment] Attempting to rollback!")
		bodyBytes, err := sendRequest(client, fmt.Sprintf("%s/?action=rollback", baseUrl), authSid, authKey, jsonString)
		if err != nil {
			log.Println("[Rollback Segment] Error on rolling back, manual interaction is needed : ", err.Error())
			return
		}

		if !strings.Contains(string(bodyBytes), `"type":"error"`) {
			fmt.Println("[Rollback Segment] Rollback has completed")
		} else {
			fmt.Println("[Rollback Segment]  Error on rolling back, manual interaction is needed")
		}
	}

	fmt.Println("[Upgrade Segment] Proccess complete ?", upgradeCompleted)

}

func sendRequest(client *http.Client, url string, authSid string, authKey string, jsonString string) (string, error) {
	httpRequest, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		return "", err
	}
	httpRequest.SetBasicAuth("2FE021DE8629C05B38BD", "DQ7iouah9GmqMZMDiFEJt7Yw8bApi9aUnfwhCdBb")
	httpRequest.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpRequest)
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
