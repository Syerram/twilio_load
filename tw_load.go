/*
Package main contains the main file that will load test the twilio endpoint

TODO:
-	Test cases
-	
*/
package main

import (
	"os"
	"time"
	"fmt"
	"flag"
	"sort"
	"bytes"
	"strings"
	"strconv"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/base64"
	"encoding/json"
	"crypto/hmac"
	"crypto/sha1"
	"github.com/vaughan0/go-ini"
)

type Config struct {
	Token string
	Users int
	URI string
	Concurrency int
	Unit float64
}

func (config *Config) populateConfig(config_file string) {
	config_data, err := ini.LoadFile(config_file)
	check(err)
	config.Token, _ = config_data.Get("Auth", "Token")
	var user_count, _ = config_data.Get("Load", "Users")
	config.Users, _ = strconv.Atoi(user_count)
	var concurrency, _ = config_data.Get("Load", "Concurrency")
	config.Concurrency, _ = strconv.Atoi(concurrency)
	config.URI, _ = config_data.Get("Load", "URI")
	var unit_cost, _ = config_data.Get("Cost", "Unit")
	config.Unit, _ = strconv.ParseFloat(unit_cost, 64)
}

//Struct for URL Parameters
type Param struct {
	Name string
	Value string
}
//Sort implementation for Param struct
type ByName []Param
func (a ByName) Len() int { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i]}
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type Request struct {
	Params []Param
	Hashed string
}

func (request *Request) unmarshal(data_file string, config Config){
	data_content, err := ini.LoadFile(data_file)
	check(err)
	var data, _ = data_content.Get("Data", "Raw")
	//unmarshall and store the signature
	json.Unmarshal([]byte(data), &request.Params)
	request.computeSignature(config)
}

func (request *Request) computeSignature(config Config) {
	sort.Sort(ByName(request.Params)) //in place ordering
	var buffer bytes.Buffer
	buffer.WriteString(config.URI)
	for _, param := range request.Params {
		buffer.WriteString(param.Name)
		buffer.WriteString(param.Value)
	}
	request.Hashed = base64.StdEncoding.EncodeToString(hmacSHA1([]byte(config.Token), buffer.String()))
}

func (request *Request) post(config Config) Response {
	data := url.Values{}
	for _, param := range request.Params {
		data.Add(param.Name, param.Value)
	}
	client := &http.Client{}
	r, _ := http.NewRequest("POST", config.URI, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("X-Twilio-Signature", request.Hashed)
	start := time.Now()
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	elapsed := time.Since(start)
	contents, _ := ioutil.ReadAll(resp.Body)
	return Response{Status: resp.Status, Time: elapsed, Content: string(contents)}
}

type Response struct {
	Status string
	Time time.Duration
	Content string
}

//-------------- Utility methods -----------------------------//
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//hmacSHA1 creates the SHA1 hash key for the given string and key
func hmacSHA1(key []byte, content string) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s: ", os.Args[0])
	fmt.Fprintf(os.Stderr, "Load test a Twilio endpoint. \nSee tw_load.sample.ini for documentation.\n\nParameters:\n")
	flag.PrintDefaults()
}

//----------------End of Utility methods-----------------------//


//worker function that executes requests for user count and sends results across the channel
func worker(c_id int, config Config, request Request, results chan <- Response) {
	for u:=1; u <= config.Users; u++ {
		results <- request.post(config)
	}
}

//---------------End of Worker pool---------------------------//

func main(){
	config_file := flag.String("f", "/path/to/config", "Config file to load test a Twilio endpoint")
	data_file := flag.String("d", "/path/to/data", "Data file to load test a Twilio endpoint")
	
	flag.Usage = usage
	flag.Parse()
	
	//load up config and data
	config := Config{}
	config.populateConfig(*config_file)
	request := Request{}
	request.unmarshal(*data_file, config)
	
	//setup the response channel
	results := make(chan Response)

	start := time.Now()
	
	for c:=1; c <= config.Concurrency; c++ {
		go worker(c, config, request, results)
	}
	
	fmt.Printf("Benchmarking %s (be patient).....\n", config.URI)
	
	success := 0
	failed := 0
	
	var total_time time.Duration
	total := config.Concurrency * config.Users
	responses := 0
	for r:=1; r<=total; r++{
		response := <- results
		if response.Status == "200 OK" {
			success += 1
		} else {
			failed += 1
		}
		//check if there is a response body, then increment the responses
		if strings.Contains(response.Content, "Body") {
			responses += 1
		}
		total_time += response.Time
		fmt.Printf("\r%.0f%%", (float32(success + failed) / float32(total)) * 100)
	}
	fmt.Printf("\rDone!\n")
	fmt.Printf("Document Path:        %s\n", config.URI)
	fmt.Printf("Concurrency Level:    %d\n", config.Concurrency)
	fmt.Printf("Users:                %d\n", config.Users)
	
	elapsed := time.Since(start)
	fmt.Printf("Completed requests:   %d\n", success)
	fmt.Printf("Failed requests:      %d\n", failed)
	fmt.Printf("Total time:           %s\n", elapsed)
	fmt.Printf("Time per request:     %s\n", (total_time / time.Duration(total)))
	fmt.Println("-----------------------------------")
	fmt.Printf("Total recv:       %d\n", total)
	fmt.Printf("Total sent        %d\n", responses)
	fmt.Printf("Approx cost:      $%.2f\n", (config.Unit * float64(total + responses)))

}