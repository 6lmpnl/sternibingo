package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

//User is a generated model from buffalo-auth, it serves as the base for username/password authentication.
type User struct {
	ID           	uuid.UUID 	`json:"id" db:"id"`
	CreatedAt    	time.Time 	`json:"created_at" db:"created_at"`
	UpdatedAt    	time.Time 	`json:"updated_at" db:"updated_at"`
	Email        	string    	`json:"email" db:"email"`
	PasswordHash 	string    	`json:"password_hash" db:"password_hash"`
	ValidationCode	string    	`json:"-" db:"validation_code"`
	Validated    	bool      	`json:"-" db:"validated"`

	Password             string `json:"-" db:"-"`
	PasswordConfirmation string `json:"-" db:"-"`

	Caps 			[]Cap 		`json:"caps" has_many:"caps"`
}

// Create wraps up the pattern of encrypting the password and
// running validations. Useful when writing tests.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(ph)

	return tx.ValidateAndCreate(u)
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
		// check to see if the email address is already taken:
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func getNewValidationKey () string {
	// generate random validation key
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, 32)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func (u *User) ValidateEmail(tx *pop.Connection) error {
	u.Validated = true;
	u.ValidationCode = ""

	return tx.Update(u)
}

// Sends validation email
func (u *User) SendValidationEmail(tx *pop.Connection) error {

	u.ValidationCode = getNewValidationKey();

	auth := smtp.PlainAuth("",
		"ludwig@loebinger.de",
		""
		"mail.loebinger.de");

	c, err := smtp.Dial("mail.loebinger.de:587")
	defer c.Close()
	if err != nil {
		return err;
	}

	tlsconfig := &tls.Config {
		InsecureSkipVerify: true,
		ServerName: "mail.loebinger.de",
	}

	c.StartTLS(tlsconfig)

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail("ludwig@loebinger.de"); err != nil {
		return err
	}

	if err = c.Rcpt(u.Email); err != nil {
		return err
	}

	wc, err := c.Data()
	defer wc.Close()
	if err != nil {
		return err
	}

	link := "http://localhost:3000/users/activate/" + u.ValidationCode
	username := strings.Split(u.Email, "@")[0]

	buf := bytes.NewBufferString("Subject: Sternibingo Registration\n"+
		"To: You <" + u.Email + ">\n" +
		"From: Ludwig Loebinger <ludwig@loebinger.de>\n" +
		"Content-Type: multipart/alternative;\n" +
		" boundary=\"------------1F66319E1598A907A365D1B7\"\n" +
		"Content-Language: en-US\n" +
		"\n" +
		"This is a multi-part message in MIME format.\n" +
		"--------------1F66319E1598A907A365D1B7\n" +
		"Content-Type: text/plain; charset=utf-8; format=flowed\n" +
		"Content-Transfer-Encoding: 7bit\n" +
		"\n" +
		"*Hello " + username + ",*\n" +
		"\n" +
		"	You registered for Social Sternibingo. Here is your activation link:\n" +
		"\n" +
		link + "\n" +
		"\n" +
		"--------------1F66319E1598A907A365D1B7\n" +
		"Content-Type: text/html; charset=utf-8\n" +
		"Content-Transfer-Encoding: 7bit\n" +
		"\n" +
		"<html><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\"></head><body><table width" +
		"=\"100%\" height=\"100%\" cellspacing=\"2\" cellpadding=\"2\" border=\"0\"><tbody><tr><td valign=\"top\"><br" +
		"></td><td valign=\"top\"><br></td><td valign=\"top\"><br></td></tr><tr><td valign=\"top\"><br></td><td width" +
		"=\"80%\" valign=\"middle\" height=\"80%\" bgcolor=\"#000000\" align=\"center\"><p><font face=\"Source Code P" +
		"ro\" color=\"#00ff15\"><b>Hello " + username + ",</b></font></p><p><font face=\"Source Code Pro\" color=\"#00ff15" +
		"\">You registered for Social Sternibingo. Here is your activation link: </font> </p> <font face=\"Source Cod" +
		"e Pro\" color=\"#00ff15\"><font color=\"#00ffbb\"><a href=\"" + link + "\" moz-do-not-send=\"true\">" + link  +
		"</a></font><br></font> <font face=\"Source Code Pro\" color=\"#00ff15\"><br></font> </td><td><br></td></tr><" +
		"tr><td><br></td><td><br></td><td><br></td></tr></tbody></table></body></html>" +
		"\n" +
		"--------------1F66319E1598A907A365D1B7--")

	if _, err = buf.WriteTo(wc); err != nil {
		return err
	}

	return tx.Update(u)
}