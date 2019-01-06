package s3versionremover

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/richardwilkes/toolbox/errs"
)

type Input struct {
	Verbose bool
	Bucket  string
}

type Output struct {
}

func S3VersionRemover(in *Input) (*Output, error) {
	in.Bucket = strings.TrimSpace(in.Bucket)
	if in.Bucket == "" {
		return nil, errs.New("bucket must not be empty")
	}
	bucket := aws.String(in.Bucket)

	sess := session.Must(session.NewSession())
	cfg := aws.NewConfig().
		WithRegion(endpoints.UsEast1RegionID)
		// .WithCredentials(credentials.NewSharedCredentials("/Users/sheeley/.aws/credentials", "aigee"))

	c := s3.New(sess, cfg)

	page := 1
	err := c.ListObjectVersionsPages(&s3.ListObjectVersionsInput{Bucket: bucket}, func(lov *s3.ListObjectVersionsOutput, lastPage bool) bool {

		toDelete := &s3.DeleteObjectsInput{
			Bucket: bucket,
			Delete: &s3.Delete{},
		}

		for _, v := range lov.Versions {
			toDelete.Delete.Objects = append(toDelete.Delete.Objects, &s3.ObjectIdentifier{
				Key:       v.Key,
				VersionId: v.VersionId,
			})
		}

		for _, m := range lov.DeleteMarkers {
			toDelete.Delete.Objects = append(toDelete.Delete.Objects, &s3.ObjectIdentifier{
				Key:       m.Key,
				VersionId: m.VersionId,
			})
		}

		objects := len(toDelete.Delete.Objects)
		fmt.Printf("%d\t%d\t%d\t%d\n", page, len(lov.Versions), len(lov.DeleteMarkers), objects)
		page++

		if objects > 0 {
			_, err := c.DeleteObjects(toDelete)
			if err != nil {
				panic(err)
			}
		}

		return !lastPage
	})

	if err != nil {
		panic(err)
	}
	return &Output{}, nil
}
