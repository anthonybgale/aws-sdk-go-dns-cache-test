package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/viki-org/dnscache"
	"net"
	"net/http"
	"strings"
	"time"
)

func setup() {
	//create a dnscache resolver
	resolver := dnscache.New(time.Minute * 5)

	//configure the default http client to use it
	http.DefaultClient.Transport = &http.Transport{
		MaxIdleConnsPerHost: 64,
		Dial: func(network string, address string) (net.Conn, error) {
			fmt.Print("Using dns cache as resolver to look up ", address)
			separator := strings.LastIndex(address, ":")
			ip, _ := resolver.FetchOneString(address[:separator])
			return net.Dial("tcp", ip+address[separator:])
		},
	}

}

func main() {

	setup()

	//create session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	//create the s3 service client
	svc := s3.New(sess)

	//argument for ListBuckets
	input := &s3.ListBucketsInput{}

	//call list buckets
	result, err := svc.ListBuckets(input)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(result)

}
