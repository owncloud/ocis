# Go Simple Mail

The best way to send emails in Go with SMTP Keep Alive and Timeout for Connect and Send.

<a href="https://goreportcard.com/report/github.com/xhit/go-simple-mail/v2"><img src="https://goreportcard.com/badge/github.com/xhit/go-simple-mail" alt="Go Report Card"></a>
<a href="https://pkg.go.dev/github.com/xhit/go-simple-mail/v2?tab=doc"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white" alt="go.dev"></a>


# IMPORTANT

Examples in this README are for v2.2.0 and above. Examples for older versions
can be found [here](https://gist.github.com/xhit/54516917473420a8db1b6fff68a21c99).

Go 1.13+ is required.

Breaking change in 2.2.0: The signature of `SetBody` and `AddAlternative` used
to accept strings ("text/html" and "text/plain") and not require on of the
`contentType` constants (`TextHTML` or `TextPlain`). Upgrading, while not
quite following semantic versioning, is quite simple:

```diff
  email := mail.NewMSG()
- email.SetBody("text/html", htmlBody)
- email.AddAlternative("text/plain", plainBody)
+ email.SetBody(mail.TextHTML, htmlBody)
+ email.AddAlternative(mail.TextPlain, plainBody)
```

# Introduction

Go Simple Mail is a simple and efficient package to send emails. It is well tested and
documented.

Go Simple Mail can only send emails using an SMTP server. But the API is flexible and it
is easy to implement other methods for sending emails using a local Postfix, an API, etc.

This package contains (and is based on) two packages by **Joe Grasse**:

- https://github.com/joegrasse/mail (unmaintained since Jun 29, 2018), and
- https://github.com/joegrasse/mime (unmaintained since Oct 1, 2015).

A lot of changes in Go Simple Mail were sent with not response.

## Features

Go Simple Mail supports:

- Multiple Attachments with path
- Multiple Attachments in base64
- Multiple Attachments from bytes (since v2.6.0)
- Inline attachments from file, base64 and bytes (bytes since v2.6.0)
- Multiple Recipients
- Priority
- Reply to
- Set sender
- Set from
- Allow sending mail with different envelope from (since v2.7.0)
- Embedded images
- HTML and text templates
- Automatic encoding of special characters
- SSL/TLS and STARTTLS
- Unencrypted connection (not recommended)
- Sending multiple emails with the same SMTP connection (Keep Alive or Persistent Connection)
- Timeout for connect to a SMTP Server
- Timeout for send an email
- Return Path
- Alternative Email Body
- CC and BCC
- Add Custom Headers in Message
- Send NOOP, RESET, QUIT and CLOSE to SMTP client
- PLAIN, LOGIN and CRAM-MD5 Authentication (since v2.3.0)
- AUTO Authentication (since v2.14.0)
- Allow connect to SMTP without authentication (since v2.10.0)
- Custom TLS Configuration (since v2.5.0)
- Send a RFC822 formatted message (since v2.8.0)
- Send from localhost (yes, Go standard SMTP package cannot do that because... WTF Google!)
- Support text/calendar content type body (since v2.11.0)
- Support add a List-Unsubscribe header (since v2.11.0)
- Support to add a DKIM signarure (since v2.11.0)
- Support to send using custom connection, ideal for proxy (since v2.15.0)
- Support Delivery Status Notification (DSN) (since v2.16.0)
- Support text/x-amp-html content type body (since v2.16.0)

## Documentation

https://pkg.go.dev/github.com/xhit/go-simple-mail/v2?tab=doc

Note 1: by default duplicated recipients throws an error, from `v2.13.0` you can use `email.AllowDuplicateAddress = true` to avoid the check.

Note 2: by default Bcc header is not set in email. From `v2.14.0` you can use `email.AddBccToHeader = true` to add this.

## Download

This package uses go modules.

```console
$ go get github.com/xhit/go-simple-mail/v2
```

# Usage

```go
package main

import (
	"log"

	mail "github.com/xhit/go-simple-mail/v2"
	"github.com/toorop/go-dkim"
)

const htmlBody = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>Hello Gophers!</title>
	</head>
	<body>
		<p>This is the <b>Go gopher</b>.</p>
		<p><img src="cid:Gopher.png" alt="Go gopher" /></p>
		<p>Image created by Renee French</p>
	</body>
</html>`

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwGrscUxi9zEa9oMOJbS0kLVHZXNIW+EBjY7KFWIZSxuGAils
wBVl+s5mMRrR5VlkyLQNulAdNemg6OSeB0R+2+8/lkHiMrimqQckZ5ig8slBoZhZ
wUoL/ZkeQa1bacbdww5TuWkiVPD9kooT/+TZW1P/ugd6oYjpOI56ZjsXzJw5pz7r
DiwcIJJaaDIqvvc5C4iW94GZjwtmP5pxhvBZ5D6Uzmh7Okvi6z4QCKzdJQLdVmC0
CMiFeh2FwqMkVpjZhNt3vtCo7Z51kwHVscel6vl51iQFq/laEzgzAWOUQ+ZEoQpL
uTaUiYzzNyEdGEzZ2CjMMoO8RgtXnUo2qX2FDQIDAQABAoIBAHWKW3kycloSMyhX
EnNSGeMz+bMtYwxNPMeebC/3xv+shoYXjAkiiTNWlfJ1MbbqjrhT1Pb1LYLbfqIF
1csWum/bjHpbMLRPO++RH1nxUJA/BMqT6HA8rWpy+JqiLW9GPf2DaP2gDYrZ0+yK
UIFG6MfzXgnju7OlkOItlvOQMY+Y501u/h6xnN2yTeRqXXJ1YlWFPRIeFdS6UOtL
J2wSxRVdymHbGwf+D7zet7ngMPwFBsbEN/83KGLRjkt8+dMQeUeob+nslsQofCZx
iokIAvByTugmqrB4JqhNkAlZhC0mqkRQh7zUFrxSj5UppMWlxLH+gPFZHKAsUJE5
mqmylcECgYEA8I/f90cpF10uH4NPBCR4+eXq1PzYoD+NdXykN65bJTEDZVEy8rBO
phXRNfw030sc3R0waQaZVhFuSgshhRuryfG9c1FP6tQhqi/jiEj9IfCW7zN9V/P2
r16pGjLuCK4SyxUC8H58Q9I0X2CQqFamtkLXC6Ogy86rZfIc8GcvZ9UCgYEAzMQZ
WAiLhRF2MEmMhKL+G3jm20r+dOzPYkfGxhIryluOXhuUhnxZWL8UZfiEqP5zH7Li
NeJvLz4pOL45rLw44qiNu6sHN0JNaKYvwNch1wPT/3/eDNZKKePqbAG4iamhjLy5
gjO1KgA5FBbcNN3R6fuJAg1e4QJCOuo55eW6vFkCgYEA7UBIV72D5joM8iFzvZcn
BPdfqh2QnELxhaye3ReFZuG3AqaZg8akWqLryb1qe8q9tclC5GIQulTInBfsQDXx
MGLNQL0x/1ylsw417kRl+qIoidMTTLocUgse5erS3haoDEg1tPBaKB1Zb7NyF8QV
+W1kX2NKg5bZbdrh9asekt0CgYA6tUam7NxDrLv8IDo/lRPSAJn/6cKG95aGERo2
k+MmQ5XP+Yxd+q0LOs24ZsZyRXHwdrNQy7khDGt5L2EN23Fb2wO3+NM6zrGu/WbX
nVbAdQKFUL3zZEUjOYtuqBemsJH27e0qHXUls6ap0dwU9DxJH6sqgXbggGtIxPsQ
pQsjEQKBgQC9gAqAj+ZtMXNG9exVPT8I15reox9kwxGuvJrRu/5eSi6jLR9z3x9P
2FrgxQ+GCB2ypoOUcliXrKesdSbolUilA8XQn/M113Lg8oA3gJXbAKqbTR/EgfUU
kvYaR/rTFnivF4SL/P4k/gABQoJuFUtSKdouELqefXB+e94g/G++Bg==
-----END RSA PRIVATE KEY-----`

func main() {
	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = "smtp.example.com"
	server.Port = 587
	server.Username = "test@example.com"
	server.Password = "examplepass"
	server.Encryption = mail.EncryptionSTARTTLS

	// You can specified authentication type:
	// - AUTO (default)
	// - PLAIN
	// - LOGIN
	// - CRAM-MD5
	// - None
	// server.Authentication = mail.AuthAuto

	// Variable to keep alive connection
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient,err := server.Connect()

	if err != nil{
		log.Fatal(err)
	}

	// New email simple html with inline and CC
	email := mail.NewMSG()
	email.SetFrom("From Example <nube@example.com>").
		AddTo("xhit@example.com").
		AddCc("otherto@example.com").
		SetSubject("New Go Email").
		SetListUnsubscribe("<mailto:unsubscribe@example.com?subject=https://example.com/unsubscribe>")

	email.SetBody(mail.TextHTML, htmlBody)

	// also you can add body from []byte with SetBodyData, example:
	// email.SetBodyData(mail.TextHTML, []byte(htmlBody))
	// or alternative part
	// email.AddAlternativeData(mail.TextHTML, []byte(htmlBody))

	// add inline
	email.Attach(&mail.File{FilePath: "/path/to/image.png", Name:"Gopher.png", Inline: true})

	// also you can set Delivery Status Notification (DSN) (only is set when server supports DSN)
	email.SetDSN([]mail.DSN{mail.SUCCESS, mail.FAILURE}, false)

	// you can add dkim signature to the email. 
	// to add dkim, you need a private key already created one.
	if privateKey != "" {
		options := dkim.NewSigOptions()
		options.PrivateKey = []byte(privateKey)
		options.Domain = "example.com"
		options.Selector = "default"
		options.SignatureExpireIn = 3600
		options.Headers = []string{"from", "date", "mime-version", "received", "received"}
		options.AddSignatureTimestamp = true
		options.Canonicalization = "relaxed/relaxed"

		email.SetDkim(options)
	}

	// always check error after send
	if email.Error != nil{
		log.Fatal(email.Error)
	}

	// Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent")
	}
}
```

## Send multiple emails in same connection

```go
	//Set your smtpClient struct to keep alive connection
	server.KeepAlive = true

	for _, to := range []string{
		"to1@example1.com",
		"to3@example2.com",
		"to4@example3.com",
	} {
		// New email simple html with inline and CC
		email := mail.NewMSG()
		email.SetFrom("From Example <nube@example.com>").
			AddTo(to).
			SetSubject("New Go Email")

		email.SetBody(mail.TextHTML, htmlBody)

		// add inline
		email.Attach(&mail.File{FilePath: "/path/to/image.png", Name:"Gopher.png", Inline: true})

		// always check error after send
		if email.Error != nil{
			log.Fatal(email.Error)
		}

		// Call Send and pass the client
		err = email.Send(smtpClient)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Email Sent")
		}
	}
```

# Send with custom connection

It's possible to use a custom connection with custom dieler, like a dialer that uses a proxy server, etc...

With this, these servers params are ignored: `Host`, `Port`. You have total control of the connection.

Example using a conn for proxy:

```go
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	mail "github.com/xhit/go-simple-mail/v2"
	"golang.org/x/net/proxy"
)

func main() {
	server := mail.NewSMTPClient()

	host := "smtp.example.com"
	port := 587
	proxyAddr := "proxy.server"

	conn, err := getCustomConnWithProxy(proxyAddr, host, port)
	if err != nil {
		log.Fatal(err)
	}

	server.Username = "test@example.com"
	server.Password = "examplepass"
	server.Encryption = mail.EncryptionSTARTTLS
	server.CustomConn = conn

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	email := mail.NewMSG()

	err = email.SetFrom("From Example <nube@example.com>").AddTo("xhit@example.com").Send(smtpClient)
	if err != nil {
		log.Fatal(err)
	}

}

func getCustomConnWithProxy(proxyAddr, host string, port int) (net.Conn, error) {
	dial, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	dialer, err := dial.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	conf := &tls.Config{ServerName: host}
	conn := tls.Client(dialer, conf)

	return conn, nil
}

```

## More examples

See [example/example_test.go](example/example_test.go).
