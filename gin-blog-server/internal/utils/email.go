package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"github.com/k3a/html2text"
	"github.com/thanhpk/randstr"
	"github.com/vanng822/go-premailer/premailer"
	"gopkg.in/gomail.v2"
)

// Registration flow: send a verification email while storing a code in local Redis.
// When the user clicks the verification link, the server receives the code and verifies it.
// If the code exists in Redis, verification succeeds; otherwise it fails.
type EmailData struct {
	URL         template.URL // Verification link
	UserName    string       // Username (email address)
	Subject     string       // Email subject
}

// Convert email address to lowercase and trim spaces
// Formatting avoids duplicate registrations due to case, and tolerates minor user input errors
func Format(email string) string{
	return strings.ToLower(strings.TrimSpace(email))
}

// Generate random string
func GetCode() string{
	code := randstr.String(24)
	return code
}

// Return base64-encoded string
func Encode (s string) string{
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return data
}
// Return base64-decoded string
func Decode (s string) (string,error){
	data ,err := base64.StdEncoding.DecodeString(s)
	if err != nil{
		return "",errors.New("emailVertify failed, decode error!!!")
	}
	return string(data),nil
}

// Return base64 string containing verification info
func GenEmailVerificationInfo(email string,password string) string{
	code := GetCode()
	info := Encode(email+"|"+password+"|"+code)   // Encoded
	return info
}

// Parse base64 string and return email and code
func ParseEmailVerificationInfo(info string) (string,string,error){
	data,err := Decode(info)
	if err!= nil{
		return "","",err
	}

	str := strings.Split(data,"|")
	if len(str) != 3{   // Invalid format
		return "","",errors.New("Wrong Vertifacion Info fomat")
	}

	return str[0],str[1],nil
}

// Generate verification URL
func GetEmailVerifyURL(info string) string{
	baseurl := g.GetConfig().Server.Port
	if baseurl[0] ==':'{
		baseurl = fmt.Sprintf("localhost%s",baseurl)
	}
	// If deploying via Docker, comment out the above and set below instead
	// baseurl := "your-domain.com"  (no port needed)

	return fmt.Sprintf("%s/api/email/verify?info=%s",baseurl,info)  // form data
}

// Build email data
func GetEmailData(email string,info string) *EmailData{
	return &EmailData{
		URL: template.URL(GetEmailVerifyURL(info)), // Verification link
		UserName: email,                            // Email address
		Subject: "Please complete account registration", // Subject
	}
}

// Parse template directory
func ParseTemplateDir(dir string) (*template.Template,error){
	var paths []string
	// Walk template directory and collect all template paths
	err := filepath.Walk(dir,func(path string,info os.FileInfo,err error) error{
				if err != nil {
					return err
				}
				if !info.IsDir(){
					paths = append(paths,path)
				}
				return nil
	})

	if err!= nil{
		return nil,err
	}
	return template.ParseFiles(paths...)
}

// Send email
// Requires SMTP mail server configuration (see config.yaml)
// Errors may occur if: 1) mail/SMTP config is wrong; 2) template parsing fails
func SendEmail(email string,data *EmailData) error{
	config := g.GetConfig().Email
	from := config.From
	Pass := config.SmtpPass
	User := config.SmtpUser
	to := email
	Host := config.Host
	Port := config.Port
	slog.Info("User:"+User+"  Pass:"+Pass+"  Host:"+Host+"  Port:")

	var body bytes.Buffer
	// Parse templates
	template,err := ParseTemplateDir("../assets/templates")
	if err != nil {
		return errors.New("Failed to parse templates")
	}
	slog.Info("Template parsed successfully\n")

	fmt.Println("URL:",data.URL)
	// Execute template
	template.ExecuteTemplate(&body,"email-verify.tpl",&data) // Render HTML into body
	// Convert HTML to inline styles for better email client compatibility
    htmlString := body.String()
	prem,_ := premailer.NewPremailerFromString(htmlString,nil)
	htmlline,_ := prem.Transform()
	m :=gomail.NewMessage()   // Send email using gomail

	slog.Info("Preparing to send email\n")
	// Set email headers
	m.SetHeader("From",from)
	m.SetHeader("To",to)
	m.SetHeader("Subject",data.Subject)

	// Set HTML body content
	m.SetBody("text/html",htmlline)	
	m.AddAlternative("text/plain",html2text.HTML2Text(body.String()))

	// Configure SMTP connection
	d := gomail.NewDialer(Host,Port,User,Pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify:true}
	slog.Info("SMTP connection established")
	if err := d.DialAndSend(m); err != nil{
		return err
	}
	return nil
}
