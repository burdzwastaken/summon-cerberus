package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"bytes"
	"regexp"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"github.com/tidwall/gjson"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type EC2IAMinfo struct {
	Code string `json:"Code"`
	LastUpdated time.Time `json:"LastUpdated"`
	InstanceProfileArn string `json:"InstanceProfileArn"`
	InstanceProfileID string `json:"InstanceProfileId"`
}

type EC2InstanceIdentityDocument struct {
	DevpayProductCodes []string `json:"devpayProductCodes"`
	AvailabilityZone string `json:"availabilityZone"`
	PrivateIP string `json:"privateIp"`
	Version string `json:"version"`
	Region string `json:"region"`
	InstanceID string `json:"instanceId"`
	BillingProducts []string `json:"billingProducts"`
	InstanceType string `json:"instanceType"`
	AccountID string `json:"accountId"`
	PendingTime time.Time `json:"pendingTime"`
	ImageID string `json:"imageId"`
	KernelID string `json:"kernelId"`
	RamdiskID string `json:"ramdiskId"`
	Architecture string `json:"architecture"`
}

type AuthDataResponse struct {
	AuthData string `json:"auth_data"`
}

type CerberusResponse struct {
	ClientToken string `json:"client_token"`
	Policies []string `json:"policies"`
	Metadata struct {
		AwsAccountID string `json:"aws_account_id"`
		AwsIamRoleName string `json:"aws_iam_role_name"`
		AwsRegion string `json:"aws_region"`
	} `json:"metadata"`
	LeaseDuration int `json:"lease_duration"`
	Renewable bool `json:"renewable"`
}

func getAWSRegion() string {
	// Build the request
	req, err := http.NewRequest("GET", "http://169.254.169.254/latest/dynamic/instance-identity/document", nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	// For control over HTTP client headers, redirect policy, and other settings,
	// create a Client - Client is an HTTP client
	client := &http.Client{}
	// Send the request via a client - Do sends an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	// Callers should close resp.Body - when done reading from it defer the closing of the body
	defer resp.Body.Close()
	// Fill the record with the data from the JSON
	var record EC2InstanceIdentityDocument
	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	// Return the AWS region
	return record.Region
}

func getIAMarn() string {
	// Build the request
	req, err := http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/iam/info", nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	// For control over HTTP client headers, redirect policy, and other settings,
	// create a Client - Client is an HTTP client
	client := &http.Client{}
	// Send the request via a client - Do sends an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	// Callers should close resp.Body - when done reading from it defer the closing of the body
	defer resp.Body.Close()
	// Fill the record with the data from the JSON
	var record EC2IAMinfo
	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	// Return the AWS ARN Information
	return record.InstanceProfileArn
}

func getAccountID() string {
	// Call getIAMarn function to get AWS ARN information
	arn := getIAMarn()
	// Extract the accountID from the AWS ARN using Regex.
	// ARN format: arn:aws:iam::012345678912:role/role-name
	re := regexp.MustCompile("[0-9]+")
	// Extract the only value from the array - Only one accountID
	accountID := re.FindAllString(arn, -1)[0]
	// Return the accountID
	  return accountID
}

func getIAMRole() string {
	// Call getIAMarn function to get AWS ARN information
	arn := getIAMarn()
	// Extract the role from the AWS ARN using strings
	// ARN format: arn:aws:iam::012345678912:role/role-name
	rolemap := strings.SplitAfter(arn, "/")
	// Assign role to second map in the arry
	role := rolemap[1]
	// Return the AWS role
	return role
}

func handlePotentialError(err error, variableName string) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			var errMessage string
			// A service error occurred
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				errMessage = fmt.Sprintf("%v %v %v", reqErr.StatusCode(), reqErr.Message(), variableName)
			} else {
				errMessage = fmt.Sprintf("%v %v", awsErr.Code(), awsErr.Message())
			}
			printAndExit(errMessage)
		} else {
			printAndExit(err.Error())
		}
	}
}

func printAndExit(err string) {
	os.Stderr.Write([]byte(err))
	os.Exit(1)
}

func debug(msg string) {
	if os.Getenv("DEBUG") == "1" {
		fmt.Println("DEBUG: ", msg)
	}
}

func authCerberus(url string) string {

	accountID := getAccountID()
	role := getIAMRole()
	region := getAWSRegion()

	debug(fmt.Sprintf("URL:> %s", url))

	var jsonStr = []byte(`{"account_id": "` + accountID + `", "role_name": "` + role + `", "region": "` + region + `"}`)

	// Build the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	client := &http.Client{}
	resp_auth, err := client.Do(req)
	if err != nil {
		printAndExit(err.Error())
	}
	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp_auth.Body.Close()
	// Print the responses for testing
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	var record_auth AuthDataResponse

	if err := json.NewDecoder(resp_auth.Body).Decode(&record_auth); err != nil {
		log.Println(err)
	}

	decodedAuthData, err := base64.StdEncoding.DecodeString(record_auth.AuthData)
	if err != nil {
		printAndExit(err.Error())
	}

	kmsClient := kms.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	params := &kms.DecryptInput{
		CiphertextBlob: decodedAuthData,
	}

	resp_cerb, err := kmsClient.Decrypt(params)
	if err != nil {
		printAndExit(err.Error())
	}

	//return string(resp.Plaintext), nil
	//fmt.Println(string(resp_cerb.Plaintext))

	debug(string(resp_cerb.Plaintext))

	var record_resp CerberusResponse

	if err := json.Unmarshal([]byte(resp_cerb.Plaintext), &record_resp); err != nil {
		log.Println(err)
	}

	debug(string(record_resp.ClientToken))

	return record_resp.ClientToken

}

func main() {

	if len(os.Args) != 2 {
		printAndExit("You must pass in one argument")
	}
	providerVar := os.Args[1]

	providerVarSlice := strings.Split(providerVar, "/")

	product, environment, secret := providerVarSlice[0], providerVarSlice[1], providerVarSlice[2]
	if len(secret) > 0 {
		secret = "." + secret
	}

	if len(os.Getenv("CERBERUS_API")) == 0 {
		printAndExit("You must set CERBERUS_API environment variable")
	}

	url := os.Getenv("CERBERUS_API") + "/v1/auth/iam-role"

	clientToken := authCerberus(url)

	url2 := os.Getenv("CERBERUS_API") + fmt.Sprintf("/v1/secret/app/%s/%s", product, environment)
	debug(fmt.Sprintf("URL:> %s", url2))

	// Build the request
	//bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest("GET", url2, nil)
	req.Header.Set("X-Vault-Token", clientToken)
	req.Header.Set("Content-Type", "application/json")

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		printAndExit(err.Error())
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printAndExit(err.Error())
	}

	secrets := gjson.Get(string(body), "data" + secret)

	fmt.Print(secrets)

}
