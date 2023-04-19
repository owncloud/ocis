//Package icapclient is a client package for the ICAP protocol
//
// Here is a basic example:
//  package main
//
//  import (
// 	  "fmt"
// 	  "log"
// 	  "net/http"
// 	  "time"
//
// 	  ic "github.com/egirna/icap-client"
//  )
//
//  func main() {
//   /* preparing the http request required for the RESPMOD */
// 	 httpReq, err := http.NewRequest(http.MethodGet, "http://localhost:8000/sample.pdf", nil)
//
// 	 if err != nil {
// 	  	log.Fatal(err)
// 	 }
//
//   /* making the http client & making the request call to get the response needed for the icap RESPMOD call */
// 	 httpClient := &http.Client{}
//
// 	 httpResp, err := httpClient.Do(httpReq)
//
// 	 if err != nil {
// 		 log.Fatal(err)
// 	 }
//
//   /* making a icap request with OPTIONS method */
// 	 optReq, err := ic.NewRequest(ic.MethodOPTIONS, "icap://127.0.0.1:1344/respmod", nil, nil)
//
// 	 if err != nil {
// 		 log.Fatal(err)
// 		 return
// 	 }
//
//   /* making the icap client responsible for making the requests */
// 	 client := &ic.Client{
// 		 Timeout: 5 * time.Second,
// 	 }
//
//   /* making the OPTIONS request call */
// 	 optResp, err := client.Do(optReq)
//
// 	 if err != nil {
// 		 log.Fatal(err)
// 		 return
// 	 }
//
//   /* making a icap request with RESPMOD method */
// 	 req, err := ic.NewRequest(ic.MethodRESPMOD, "icap://127.0.0.1:1344/respmod", httpReq, httpResp)
//
// 	 if err != nil {
// 		 log.Fatal(err)
// 	 }
//
// 	 req.SetPreview(optResp.PreviewBytes) // setting the preview bytes obtained from the OPTIONS call
//
//   /* making the RESPMOD request call */
// 	 resp, err := client.Do(req)
//
// 	 if err != nil {
// 		 log.Fatal(err)
// 	 }
//
// 	 fmt.Println(resp.StatusCode)
//
//  }
// See https://github.com/egirna/icap-client/examples.
package icapclient
