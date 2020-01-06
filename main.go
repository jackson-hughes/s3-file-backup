/*
mvp needs:
takes input for filename / pattern
validate input
takes input for destination bucket

nice to haves:
todo: boolean for delete from disk after backup
todo: compress prior to upload if not already compressed
todo: sends event to Cloud watch events with exit status (success, error)
todo: no-op / dry run mode - logs what would have happened
todo: optionally prefix new s3 objects (e.g. today's date)
*/

package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws/session"

	log "github.com/sirupsen/logrus"
)

var (
	filepattern  string
	loglevel     string
	backupBucket string
)

func main() {
	flag.StringVar(&filepattern, "filepath", "", "Absolute path of file to be backed up")
	flag.StringVar(&filepattern, "f", "", "Short flag - absolute path of file to be backed up")
	flag.StringVar(&loglevel, "log-level", "info", "Set log level")
	flag.StringVar(&loglevel, "l", "info", "Short flag - set log level")
	flag.StringVar(&backupBucket, "backup-bucket", "", "Specify S3 bucket to copy file to")
	flag.StringVar(&backupBucket, "b", "", "Short flag - specify S3 bucket to copy file to")
	flag.Parse()

	checkFlags()

	ll, err := log.ParseLevel(loglevel)
	if err != nil {
		log.Error(err)
	}
	log.SetLevel(ll)
	log.SetOutput(os.Stdout)

	f, err := findFile(filepattern)
	if err != nil {
		log.Error(err)
	}

	if f != nil {
		log.Debug("Found ", len(f), " files that matched provided pattern")
		log.Debug(f)
	} else {
		log.Infof("No files found matching pattern: %v", filepattern)
		os.Exit(0)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	uploader := s3manager.NewUploader(sess)
	// log.Info(f[1])
	fileToUpload, err := os.Open(f[0])
	if err != nil {
		log.Error(err)
	}
	defer fileToUpload.Close()
	r, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(backupBucket),
		Key:    aws.String("file.txt"),
		Body:   fileToUpload,
	})
	if err != nil {
		log.Error(err)
	}
	log.Infof("file uploaded to, %s", r.Location)
	// for _, file := range f {
	// 	err := uploadToS3(file)
	// 	if err != nil {
	// 		log.Error(err)
	// 	}
	// }
}

func checkFlags() {
	if len(filepattern) == 0 {
		log.Fatalln("No file path provided. Please provide a file path with -filepath / -f")
		os.Exit(1)
	}

	if !filepath.IsAbs(filepattern) {
		log.Fatalln("Detected relative path - please provide absolute path")
		os.Exit(1)
	}

	if len(backupBucket) == 0 {
		log.Fatalln("No S3 bucket provided. Please provide the name of a target S3 bucket with -backup-bucket / -b")
	}
}

func uploadToS3(file string) (err error) {

	log.Debugf("Uploading %v to %v", file, backupBucket)

	return
}

func findFile(pattern string) (match []string, err error) {
	fileMatches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, error(err)
	}
	return fileMatches, nil
}
