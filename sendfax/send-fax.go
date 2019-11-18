package sendfax

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/sfreiberg/gotwilio"
)

type Input struct {
	Verbose, Wait, DeleteAfterSend                bool
	AccountSID, AuthToken, From, To, File, Bucket string
}

type Output struct {
	Status string
}

func WaitUntilSent(twilio *gotwilio.Twilio, faxSID string, verbose bool) (string, error) {
	maxAttempts := 20
	attempts := 0

	for {
		if attempts > maxAttempts {
			break
		}
		time.Sleep(30 * time.Second)

		if verbose {
			fmt.Printf("checking status: %d\t%s\n", attempts, faxSID)
		}

		f, ex, err := twilio.GetFax(faxSID)
		if ex != nil {
			return "", errs.New(ex.Error())
		}
		if err != nil {
			return "", errs.Wrap(err)
		}

		if f.Status != "sending" {
			return f.Status, nil
		}

		attempts += 1
	}

	return "sending", nil
}

// SendFax uploads your file to S3, generates a presigned URL, and then sends a fax using that file location
func SendFax(in *Input, sess *session.Session) (*Output, error) {
	twilio := gotwilio.NewTwilioClient(in.AccountSID, in.AuthToken)

	f, err := os.Open(in.File)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	fileName := filepath.Base(in.File)
	ext := filepath.Ext(fileName)
	r := rand.Float32()
	const format = "%s-%f%s"
	key := aws.String(fmt.Sprintf(format, fileName, r, ext))
	bucket := aws.String(in.Bucket)

	s3c := s3.New(sess)
	_, err = s3c.PutObject(&s3.PutObjectInput{
		Bucket: bucket,
		Key:    key,
		Body:   f,
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	gor, _ := s3c.GetObjectRequest(&s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	})

	mediaURL, err := gor.Presign(60 * time.Minute)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if in.Verbose {
		fmt.Printf("media URL: %s\n", mediaURL)
	}

	fr, ex, err := twilio.SendFax(in.To, in.From, mediaURL, "fine", "", false)
	if ex != nil {
		return nil, errs.New(ex.Error())
	}
	if err != nil {
		return nil, errs.Wrap(err)
	}

	status := fr.Status
	if in.Verbose {
		fmt.Printf("enqueued: %s\t%s\n", fr.Sid, fr.Status)
	}

	if in.Wait {
		status, err = WaitUntilSent(twilio, fr.Sid, in.Verbose)
	}

	if in.DeleteAfterSend {
		_, dErr := s3c.DeleteObject(&s3.DeleteObjectInput{
			Bucket: bucket,
			Key:    key,
		})
		if dErr != nil {
			err = errs.Append(err, dErr)
		}
	}

	if err != nil {
		return nil, err
	}

	return &Output{
		Status: status,
	}, nil
}
