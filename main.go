package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const getPods string = "/api/v1/pods"
const getSCC string = "/apis/security.openshift.io/v1/securitycontextconstraints"

type PodList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		Metadata struct {
			Name        string `json:"name"`
			Namespace   string `json:"namespace"`
			Annotations struct {
				SCC string `json:"openshift.io/scc"`
			}
		}
	}
}

func getArgs() (string, string) {
	var flagAPIURL string
	var flagToken string

	flag.StringVar(&flagAPIURL, "api", "", "OpenShift API URL")
	flag.StringVar(&flagToken, "token", "", "OpenShift API token from the console \"Copy Login Command\"")

	flag.Parse()

	if flagAPIURL == "" || flagToken == "" || len(flagToken) < 10 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return flagAPIURL, flagToken
}

// todo: better url combining
func getData(url string, token string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	// u, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println(string(u))
	var podList PodList

	err = json.NewDecoder(resp.Body).Decode(&podList)

	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(pod.metadata.name)
	fmt.Println(podList)

	return string("")
}

func getClusterInfo(url string, token string) (string, string) {
	podInfo := getData(url+getPods, token)
	// sccInfo := getData(url+getSCC, token)

	return podInfo, ""
}

func main() {
	apiURL, token := getArgs()

	_, _ = getClusterInfo(apiURL, token)

}
