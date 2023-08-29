package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/toorop/go-dkim"
)

// Email represents an email message.
type Email struct {
	from                      string
	sender                    string
	replyTo                   string
	returnPath                string
	recipients                []string
	headers                   textproto.MIMEHeader
	parts                     []part
	attachments               []*File
	inlines                   []*File
	Charset                   string
	Encoding                  encoding
	Error                     error
	SMTPServer                *smtpClient
	DkimMsg                   string
	AllowDuplicateAddress     bool
	AddBccToHeader            bool
	preserveOriginalRecipient bool
	dsn                       []DSN
}

/*
SMTPServer represents a SMTP Server
If authentication is CRAM-MD5 then the Password is the Secret
*/
type SMTPServer struct {
	Authentication AuthType
	Encryption     Encryption
	Username       string
	Password       string
	Helo           string
	ConnectTimeout time.Duration
	SendTimeout    time.Duration
	Host           string
	Port           int
	KeepAlive      bool
	TLSConfig      *tls.Config

	// use custom dialer
	CustomConn net.Conn
}

// SMTPClient represents a SMTP Client for send email
type SMTPClient struct {
	mu                        sync.Mutex
	Client                    *smtpClient
	SendTimeout               time.Duration
	KeepAlive                 bool
	hasDSNExt                 bool
	preserveOriginalRecipient bool
	dsn                       []DSN
}

// part represents the different content parts of an email body.
type part struct {
	contentType string
	body        *bytes.Buffer
}

// Encryption type to enum encryption types (None, SSL/TLS, STARTTLS)
type Encryption int

// TODO: Remove EncryptionSSL and EncryptionTLS before launch v3

const (
	// EncryptionNone uses no encryption when sending email
	EncryptionNone Encryption = iota
	// EncryptionSSL: DEPRECATED. Use EncryptionSSLTLS. Sets encryption type to SSL/TLS when sending email
	EncryptionSSL
	// EncryptionTLS: DEPRECATED. Use EncryptionSTARTTLS. sets encryption type to STARTTLS when sending email
	EncryptionTLS
	// EncryptionSSLTLS sets encryption type to SSL/TLS when sending email
	EncryptionSSLTLS
	// EncryptionSTARTTLS sets encryption type to STARTTLS when sending email
	EncryptionSTARTTLS
)

// TODO: Remove last two indexes
var encryptionTypes = [...]string{"None", "SSL/TLS", "STARTTLS", "SSL/TLS", "STARTTLS"}

func (encryption Encryption) String() string {
	return encryptionTypes[encryption]
}

type encoding int

const (
	// EncodingNone turns off encoding on the message body
	EncodingNone encoding = iota
	// EncodingBase64 sets the message body encoding to base64
	EncodingBase64
	// EncodingQuotedPrintable sets the message body encoding to quoted-printable
	EncodingQuotedPrintable
)

var encodingTypes = [...]string{"binary", "base64", "quoted-printable"}

func (encoding encoding) string() string {
	return encodingTypes[encoding]
}

type ContentType int

const (
	// TextPlain sets body type to text/plain in message body
	TextPlain ContentType = iota
	// TextHTML sets body type to text/html in message body
	TextHTML
	// TextCalendar sets body type to text/calendar in message body
	TextCalendar
	// TextAMP sets body type to text/x-amp-html in message body
	TextAMP
)

var contentTypes = [...]string{"text/plain", "text/html", "text/calendar", "text/x-amp-html"}

func (contentType ContentType) string() string {
	return contentTypes[contentType]
}

type AuthType int

const (
	// AuthPlain implements the PLAIN authentication
	AuthPlain AuthType = iota
	// AuthLogin implements the LOGIN authentication
	AuthLogin
	// AuthCRAMMD5 implements the CRAM-MD5 authentication
	AuthCRAMMD5
	// AuthNone for SMTP servers without authentication
	AuthNone
	// AuthAuto (default) use the first AuthType of the list of returned types supported by SMTP
	AuthAuto
)

func (at AuthType) String() string {
	switch at {
	case AuthPlain:
		return "PLAIN"
	case AuthLogin:
		return "LOGIN"
	case AuthCRAMMD5:
		return "CRAM-MD5"
	default:
		return ""
	}
}

/*
	DSN notifications

- 'NEVER' under no circumstances a DSN must be returned to the sender. If you use NEVER all other notifications will be ignored.

- 'SUCCESS' will notify you when your mail has arrived at its destination.

- 'FAILURE' will arrive if an error occurred during delivery.

- 'DELAY' will notify you if there is an unusual delay in delivery, but the actual delivery's outcome (success or failure) is not yet decided.

see https://tools.ietf.org/html/rfc3461 See section 4.1 for more information about NOTIFY
*/
type DSN int

const (
	NEVER DSN = iota
	FAILURE
	DELAY
	SUCCESS
)

var dsnTypes = [...]string{"NEVER", "FAILURE", "DELAY", "SUCCESS"}

func (dsn DSN) String() string {
	return dsnTypes[dsn]
}

// NewMSG creates a new email. It uses UTF-8 by default. All charsets: http://webcheatsheet.com/HTML/character_sets_list.php
func NewMSG() *Email {
	email := &Email{
		headers:  make(textproto.MIMEHeader),
		Charset:  "UTF-8",
		Encoding: EncodingQuotedPrintable,
	}

	email.AddHeader("MIME-Version", "1.0")

	return email
}

// NewSMTPClient returns the client for send email
func NewSMTPClient() *SMTPServer {
	server := &SMTPServer{
		Authentication: AuthAuto,
		Encryption:     EncryptionNone,
		ConnectTimeout: 10 * time.Second,
		SendTimeout:    10 * time.Second,
		Helo:           "localhost",
	}
	return server
}

// GetEncryptionType returns the encryption type used to connect to SMTP server
func (server *SMTPServer) GetEncryptionType() Encryption {
	return server.Encryption
}

// GetError returns the first email error encountered
func (email *Email) GetError() error {
	return email.Error
}

// SetFrom sets the From address.
func (email *Email) SetFrom(address string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("From", address)

	return email
}

// SetSender sets the Sender address.
func (email *Email) SetSender(address string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("Sender", address)

	return email
}

// SetReplyTo sets the Reply-To address.
func (email *Email) SetReplyTo(address string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("Reply-To", address)

	return email
}

// SetReturnPath sets the Return-Path address. This is most often used
// to send bounced emails to a different email address.
func (email *Email) SetReturnPath(address string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("Return-Path", address)

	return email
}

// AddTo adds a To address. You can provide multiple
// addresses at the same time.
func (email *Email) AddTo(addresses ...string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("To", addresses...)

	return email
}

// AddCc adds a Cc address. You can provide multiple
// addresses at the same time.
func (email *Email) AddCc(addresses ...string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("Cc", addresses...)

	return email
}

// AddBcc adds a Bcc address. You can provide multiple
// addresses at the same time.
func (email *Email) AddBcc(addresses ...string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddAddresses("Bcc", addresses...)

	return email
}

// AddAddresses allows you to add addresses to the specified address header.
func (email *Email) AddAddresses(header string, addresses ...string) *Email {
	if email.Error != nil {
		return email
	}

	found := false

	// check for a valid address header
	for _, h := range []string{"To", "Cc", "Bcc", "From", "Sender", "Reply-To", "Return-Path"} {
		if header == h {
			found = true
		}
	}

	if !found {
		email.Error = errors.New("Mail Error: Invalid address header; Header: [" + header + "]")
		return email
	}

	// check to see if the addresses are valid
	for i := range addresses {
		var address = new(mail.Address)
		var err error

		// ignore parse the address if empty
		if len(addresses[i]) > 0 {
			address, err = mail.ParseAddress(addresses[i])
			if err != nil {
				email.Error = errors.New("Mail Error: " + err.Error() + "; Header: [" + header + "] Address: [" + addresses[i] + "]")
				return email
			}
		} else {
			continue
		}

		// check for more than one address
		switch {
		case header == "Sender" && len(email.sender) > 0:
			fallthrough
		case header == "Reply-To" && len(email.replyTo) > 0:
			fallthrough
		case header == "Return-Path" && len(email.returnPath) > 0:
			email.Error = errors.New("Mail Error: There can only be one \"" + header + "\" address; Header: [" + header + "] Address: [" + addresses[i] + "]")
			return email
		default:
			// other address types can have more than one address
		}

		// save the address
		switch header {
		case "From":
			// delete the current "From" to set the new
			// when "From" need to be changed in the message
			if len(email.from) > 0 && header == "From" {
				email.headers.Del("From")
			}
			email.from = address.Address
		case "Sender":
			email.sender = address.Address
		case "Reply-To":
			email.replyTo = address.Address
		case "Return-Path":
			email.returnPath = address.Address
		default:
			// check that the address was added to the recipients list
			email.recipients, err = addAddress(email.recipients, address.Address, email.AllowDuplicateAddress)
			if err != nil {
				email.Error = errors.New("Mail Error: " + err.Error() + "; Header: [" + header + "] Address: [" + addresses[i] + "]")
				return email
			}
		}

		// make sure the from and sender addresses are different
		if email.from != "" && email.sender != "" && email.from == email.sender {
			email.sender = ""
			email.headers.Del("Sender")
			email.Error = errors.New("Mail Error: From and Sender should not be set to the same address")
			return email
		}

		// add Bcc only if AddBccToHeader is true
		if header == "Bcc" && email.AddBccToHeader {
			email.headers.Add(header, address.String())
		}

		// add all addresses to the headers except for Bcc and Return-Path
		if header != "Bcc" && header != "Return-Path" {
			// add the address to the headers
			email.headers.Add(header, address.String())
		}
	}

	return email
}

// addAddress adds an address to the address list if it hasn't already been added
func addAddress(addressList []string, address string, allowDuplicateAddress bool) ([]string, error) {
	if !allowDuplicateAddress {
		// loop through the address list to check for dups
		for _, a := range addressList {
			if address == a {
				return addressList, errors.New("Mail Error: Address: [" + address + "] has already been added")
			}
		}
	}

	return append(addressList, address), nil
}

type Priority int

const (
	// PriorityLow sets the email Priority to Low
	PriorityLow Priority = iota
	// PriorityHigh sets the email Priority to High
	PriorityHigh
)

// SetPriority sets the email message Priority. Use with
// either "High" or "Low".
func (email *Email) SetPriority(priority Priority) *Email {
	if email.Error != nil {
		return email
	}

	switch priority {
	case PriorityLow:
		email.AddHeaders(textproto.MIMEHeader{
			"X-Priority":        {"5 (Lowest)"},
			"X-MSMail-Priority": {"Low"},
			"Importance":        {"Low"},
		})
	case PriorityHigh:
		email.AddHeaders(textproto.MIMEHeader{
			"X-Priority":        {"1 (Highest)"},
			"X-MSMail-Priority": {"High"},
			"Importance":        {"High"},
		})
	default:
	}

	return email
}

// SetDate sets the date header to the provided date/time.
// The format of the string should be YYYY-MM-DD HH:MM:SS Time Zone.
//
// Example: SetDate("2015-04-28 10:32:00 CDT")
func (email *Email) SetDate(dateTime string) *Email {
	if email.Error != nil {
		return email
	}

	const dateFormat = "2006-01-02 15:04:05 MST"

	// Try to parse the provided date/time
	dt, err := time.Parse(dateFormat, dateTime)
	if err != nil {
		email.Error = errors.New("Mail Error: Setting date failed with: " + err.Error())
		return email
	}

	email.headers.Set("Date", dt.Format(time.RFC1123Z))

	return email
}

// SetSubject sets the subject of the email message.
func (email *Email) SetSubject(subject string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddHeader("Subject", subject)

	return email
}

// SetListUnsubscribe sets the Unsubscribe address.
func (email *Email) SetListUnsubscribe(address string) *Email {
	if email.Error != nil {
		return email
	}

	email.AddHeader("List-Unsubscribe", address)

	return email
}

// SetDkim adds DomainKey signature to the email message (header+body)
func (email *Email) SetDkim(options dkim.SigOptions) *Email {
	if email.Error != nil {
		return email
	}

	msg := []byte(email.GetMessage())
	err := dkim.Sign(&msg, options)

	if err != nil {
		email.Error = errors.New("Mail Error: cannot dkim sign message due: %s" + err.Error())
		return email
	}

	email.DkimMsg = string(msg)

	return email
}

// SetBody sets the body of the email message.
func (email *Email) SetBody(contentType ContentType, body string) *Email {
	if email.Error != nil {
		return email
	}

	email.parts = []part{
		{
			contentType: contentType.string(),
			body:        bytes.NewBufferString(body),
		},
	}

	return email
}

// SetBodyData sets the body of the email message from []byte
func (email *Email) SetBodyData(contentType ContentType, body []byte) *Email {
	if email.Error != nil {
		return email
	}

	email.parts = []part{
		{
			contentType: contentType.string(),
			body:        bytes.NewBuffer(body),
		},
	}

	return email
}

// AddHeader adds the given "header" with the passed "value".
func (email *Email) AddHeader(header string, values ...string) *Email {
	if email.Error != nil {
		return email
	}

	// check that there is actually a value
	if len(values) < 1 {
		email.Error = errors.New("Mail Error: no value provided; Header: [" + header + "]")
		return email
	}

	if header != "MIME-Version" {
		// Set header to correct canonical Mime
		header = textproto.CanonicalMIMEHeaderKey(header)
	}

	switch header {
	case "Sender":
		fallthrough
	case "From":
		fallthrough
	case "To":
		fallthrough
	case "Bcc":
		fallthrough
	case "Cc":
		fallthrough
	case "Reply-To":
		fallthrough
	case "Return-Path":
		email.AddAddresses(header, values...)
	case "Date":
		if len(values) > 1 {
			email.Error = errors.New("Mail Error: To many dates provided")
			return email
		}
		email.SetDate(values[0])
	case "List-Unsubscribe":
		fallthrough
	default:
		email.headers[header] = values
	}

	return email
}

// AddHeaders is used to add multiple headers at once
func (email *Email) AddHeaders(headers textproto.MIMEHeader) *Email {
	if email.Error != nil {
		return email
	}

	for header, values := range headers {
		email.AddHeader(header, values...)
	}

	return email
}

// AddAlternative allows you to add alternative parts to the body
// of the email message. This is most commonly used to add an
// html version in addition to a plain text version that was
// already added with SetBody.
func (email *Email) AddAlternative(contentType ContentType, body string) *Email {
	if email.Error != nil {
		return email
	}

	email.parts = append(email.parts,
		part{
			contentType: contentType.string(),
			body:        bytes.NewBufferString(body),
		},
	)

	return email
}

// AddAlternativeData allows you to add alternative parts to the body
// of the email message. This is most commonly used to add an
// html version in addition to a plain text version that was
// already added with SetBody.
func (email *Email) AddAlternativeData(contentType ContentType, body []byte) *Email {
	if email.Error != nil {
		return email
	}

	email.parts = append(email.parts,
		part{
			contentType: contentType.string(),
			body:        bytes.NewBuffer(body),
		},
	)

	return email
}

// SetDSN sets the delivery status notification list, only is set when SMTP server supports DSN extension
//
// To preserve the original recipient of an email message, for example, if it is forwarded to another address, set preserveOriginalRecipient to true
func (email *Email) SetDSN(dsn []DSN, preserveOriginalRecipient bool) *Email {
	if email.Error != nil {
		return email
	}

	email.dsn = dsn
	email.preserveOriginalRecipient = preserveOriginalRecipient

	return email
}

// GetFrom returns the sender of the email, if any
func (email *Email) GetFrom() string {
	from := email.returnPath
	if from == "" {
		from = email.sender
		if from == "" {
			from = email.from
			if from == "" {
				from = email.replyTo
			}
		}
	}

	return from
}

// GetRecipients returns a slice of recipients emails
func (email *Email) GetRecipients() []string {
	return email.recipients
}

func (email *Email) hasMixedPart() bool {
	return (len(email.parts) > 0 && len(email.attachments) > 0) || len(email.attachments) > 1
}

func (email *Email) hasRelatedPart() bool {
	return (len(email.parts) > 0 && len(email.inlines) > 0) || len(email.inlines) > 1
}

func (email *Email) hasAlternativePart() bool {
	return len(email.parts) > 1
}

// GetMessage builds and returns the email message (RFC822 formatted message)
func (email *Email) GetMessage() string {
	msg := newMessage(email)

	if email.hasMixedPart() {
		msg.openMultipart("mixed")
	}

	if email.hasRelatedPart() {
		msg.openMultipart("related")
	}

	if email.hasAlternativePart() {
		msg.openMultipart("alternative")
	}

	for _, part := range email.parts {
		msg.addBody(part.contentType, part.body.Bytes())
	}

	if email.hasAlternativePart() {
		msg.closeMultipart()
	}

	msg.addFiles(email.inlines, true)
	if email.hasRelatedPart() {
		msg.closeMultipart()
	}

	msg.addFiles(email.attachments, false)
	if email.hasMixedPart() {
		msg.closeMultipart()
	}

	return msg.getHeaders() + msg.body.String()
}

// Send sends the composed email
func (email *Email) Send(client *SMTPClient) error {
	return email.SendEnvelopeFrom(email.from, client)
}

// SendEnvelopeFrom sends the composed email with envelope
// sender. 'from' must be an email address.
func (email *Email) SendEnvelopeFrom(from string, client *SMTPClient) error {
	if email.Error != nil {
		return email.Error
	}

	if from == "" {
		from = email.from
	}

	if len(email.recipients) < 1 {
		return errors.New("Mail Error: No recipient specified")
	}

	var msg string
	if email.DkimMsg != "" {
		msg = email.DkimMsg
	} else {
		msg = email.GetMessage()
	}

	client.dsn = email.dsn
	client.preserveOriginalRecipient = email.preserveOriginalRecipient

	return send(from, email.recipients, msg, client)
}

// dial connects to the smtp server with the request encryption type
func dial(customConn net.Conn, host string, port string, encryption Encryption, config *tls.Config) (*smtpClient, error) {
	var conn net.Conn
	var err error
	var c *smtpClient

	if customConn != nil {
		conn = customConn
	} else {
		address := host + ":" + port
		// do the actual dial
		switch encryption {
		// TODO: Remove EncryptionSSL check before launch v3
		case EncryptionSSL, EncryptionSSLTLS:
			conn, err = tls.Dial("tcp", address, config)
		default:
			conn, err = net.Dial("tcp", address)
		}

		if err != nil {
			return nil, errors.New("Mail Error on dialing with encryption type " + encryption.String() + ": " + err.Error())
		}
	}

	c, err = newClient(conn, host)
	if err != nil {
		return nil, fmt.Errorf("Mail Error on smtp dial: %w", err)
	}

	return c, err
}

// smtpConnect connects to the smtp server and starts TLS and passes auth
// if necessary
func smtpConnect(customConn net.Conn, host, port, helo string, encryption Encryption, config *tls.Config) (*smtpClient, error) {
	// connect to the mail server
	c, err := dial(customConn, host, port, encryption, config)

	if err != nil {
		return nil, err
	}

	if helo == "" {
		helo = "localhost"
	}

	// send Helo
	if err = c.hi(helo); err != nil {
		c.close()
		return nil, fmt.Errorf("Mail Error on Hello: %w", err)
	}

	// STARTTLS if necessary
	// TODO: Remove EncryptionTLS check before launch v3
	if encryption == EncryptionTLS || encryption == EncryptionSTARTTLS {
		if ok, _ := c.extension("STARTTLS"); ok {
			if err = c.startTLS(config); err != nil {
				c.close()
				return nil, fmt.Errorf("Mail Error on STARTTLS: %w", err)
			}
		}
	}

	return c, nil
}

func (server *SMTPServer) getAuth(a string) (auth, error) {
	var afn auth
	switch {
	case strings.Contains(a, AuthPlain.String()):
		if server.Username != "" || server.Password != "" {
			afn = plainAuthfn("", server.Username, server.Password, server.Host)
		}
	case strings.Contains(a, AuthLogin.String()):
		if server.Username != "" || server.Password != "" {
			afn = loginAuthfn("", server.Username, server.Password, server.Host)
		}
	case strings.Contains(a, AuthCRAMMD5.String()):
		if server.Username != "" || server.Password != "" {
			afn = cramMD5Authfn(server.Username, server.Password)
		}
	default:
		return nil, fmt.Errorf("Mail Error on determining auth type, %s is not supported", a)
	}
	return afn, nil
}

func (server *SMTPServer) validateAuth(c *smtpClient) error {
	var err error
	var afn auth
	switch {
	case server.Authentication == AuthNone || server.Username == "":
		return nil
	case server.Authentication != AuthAuto:
		afn, err = server.getAuth(server.Authentication.String())
		if err != nil {
			return err
		}
	}
	if ok, a := c.extension("AUTH"); ok {
		// Determine Auth type automatically from extension
		if afn == nil {
			afn, err = server.getAuth(a)
			if err != nil {
				return err
			}
		}
		if err = c.authenticate(afn); err != nil {
			c.close()
			return fmt.Errorf("Mail Error on Auth: %w", err)
		}
	}
	return nil
}

// Connect returns the smtp client
func (server *SMTPServer) Connect() (*SMTPClient, error) {
	var smtpConnectChannel chan error
	var c *smtpClient
	var err error

	tlsConfig := server.TLSConfig
	if tlsConfig == nil {
		tlsConfig = &tls.Config{ServerName: server.Host}
	}

	// if there is a ConnectTimeout, setup the channel and do the connect under a goroutine
	if server.ConnectTimeout != 0 {
		smtpConnectChannel = make(chan error, 2)
		go func() {
			c, err = smtpConnect(server.CustomConn, server.Host, fmt.Sprintf("%d", server.Port), server.Helo, server.Encryption, tlsConfig)
			// send the result
			smtpConnectChannel <- err
		}()
		// get the connect result or timeout result, which ever happens first
		select {
		case err = <-smtpConnectChannel:
			if err != nil {
				return nil, err
			}
		case <-time.After(server.ConnectTimeout):
			return nil, errors.New("Mail Error: SMTP Connection timed out")
		}
	} else {
		// no ConnectTimeout, just fire the connect
		c, err = smtpConnect(server.CustomConn, server.Host, fmt.Sprintf("%d", server.Port), server.Helo, server.Encryption, tlsConfig)
		if err != nil {
			return nil, err
		}
	}

	_, hasDSN := c.ext["DSN"]

	return &SMTPClient{
		Client:      c,
		KeepAlive:   server.KeepAlive,
		SendTimeout: server.SendTimeout,
		hasDSNExt:   hasDSN,
	}, server.validateAuth(c)
}

// Reset send RSET command to smtp client
func (smtpClient *SMTPClient) Reset() error {
	smtpClient.mu.Lock()
	defer smtpClient.mu.Unlock()
	return smtpClient.Client.reset()
}

// Noop send NOOP command to smtp client
func (smtpClient *SMTPClient) Noop() error {
	smtpClient.mu.Lock()
	defer smtpClient.mu.Unlock()
	return smtpClient.Client.noop()
}

// Quit send QUIT command to smtp client
func (smtpClient *SMTPClient) Quit() error {
	smtpClient.mu.Lock()
	defer smtpClient.mu.Unlock()
	return smtpClient.Client.quit()
}

// Close closes the connection
func (smtpClient *SMTPClient) Close() error {
	smtpClient.mu.Lock()
	defer smtpClient.mu.Unlock()
	return smtpClient.Client.close()
}

// SendMessage sends a message (a RFC822 formatted message)
// 'from' must be an email address, recipients must be a slice of email address
func SendMessage(from string, recipients []string, msg string, client *SMTPClient) error {
	if from == "" {
		return errors.New("Mail Error: No From email specifier")
	}
	if len(recipients) < 1 {
		return errors.New("Mail Error: No recipient specified")
	}

	return send(from, recipients, msg, client)
}

// send does the low level sending of the email
func send(from string, to []string, msg string, client *SMTPClient) error {
	//Check if client struct is not nil
	if client != nil {

		//Check if client is not nil
		if client.Client != nil {
			var smtpSendChannel chan error

			// if there is a SendTimeout, setup the channel and do the send under a goroutine
			if client.SendTimeout != 0 {
				smtpSendChannel = make(chan error, 1)

				go func(from string, to []string, msg string, client *SMTPClient) {
					smtpSendChannel <- sendMailProcess(from, to, msg, client)
				}(from, to, msg, client)
			}

			if client.SendTimeout == 0 {
				// no SendTimeout, just fire the sendMailProcess
				return sendMailProcess(from, to, msg, client)
			}

			// get the send result or timeout result, which ever happens first
			select {
			case sendError := <-smtpSendChannel:
				checkKeepAlive(client)
				return sendError
			case <-time.After(client.SendTimeout):
				checkKeepAlive(client)
				return errors.New("Mail Error: SMTP Send timed out")
			}
		}
	}

	return errors.New("Mail Error: No SMTP Client Provided")
}

func sendMailProcess(from string, to []string, msg string, c *SMTPClient) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cmdArgs := make(map[string]string)

	if _, ok := c.Client.ext["SIZE"]; ok {
		cmdArgs["SIZE"] = strconv.Itoa(len(msg))
	}

	// Set the sender
	if err := c.Client.mail(from, cmdArgs); err != nil {
		return err
	}

	var dsn string
	var dsnSet bool

	if c.hasDSNExt && len(c.dsn) > 0 {
		dsn = " NOTIFY="
		if hasNeverDSN(c.dsn) {
			dsn += NEVER.String()
		} else {
			dsn += strings.Join(dsnToString(c.dsn), ",")
		}

		if c.preserveOriginalRecipient {
			dsn += " ORCPT=rfc822;"
		}

		dsnSet = true
	}

	// Set the recipients
	for _, address := range to {
		if dsnSet && c.preserveOriginalRecipient {
			dsn += address
		}

		if err := c.Client.rcpt(address, dsn); err != nil {
			return err
		}
	}

	// Send the data command
	w, err := c.Client.data()
	if err != nil {
		return err
	}

	// write the message
	_, err = fmt.Fprint(w, msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

// check if keepAlive for close or reset
func checkKeepAlive(client *SMTPClient) {
	if client.KeepAlive {
		client.Reset()
	} else {
		client.Quit()
		client.Close()
	}
}

func hasNeverDSN(dsnList []DSN) bool {
	for i := range dsnList {
		if dsnList[i] == NEVER {
			return true
		}
	}
	return false
}

func dsnToString(dsnList []DSN) []string {
	dsnString := make([]string, len(dsnList))
	for i := range dsnList {
		dsnString[i] = dsnList[i].String()
	}
	return dsnString
}
