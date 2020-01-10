/*
mvp needs:
takes input for filename / pattern
validate input
takes input for destination bucket
uploads file to bucket

nice to haves:
boolean for delete from disk after backup
todo: compress prior to upload if not already compressed
todo: publishes metric to cloud watch metrics
todo: no-op / dry run mode - logs what would have happened
optionally prefix new s3 objects (e.g. today's date)
flag for profiling
*/

package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/profile"

	log "github.com/sirupsen/logrus"
)

var (
	gitCommit     string
	gitUrl        string
	filepattern   string
	loglevel      string
	backupBucket  string
	profileMemory bool
	objectPrefix  string
	deleteFiles   bool
)

func main() {
	flag.StringVar(&filepattern, "filepath", "", "Absolute path of file to be backed up")
	flag.StringVar(&filepattern, "f", "", "Short flag - absolute path of file to be backed up")
	flag.StringVar(&loglevel, "log-level", "info", "Set log level")
	flag.StringVar(&loglevel, "l", "info", "Short flag - set log level")
	flag.StringVar(&backupBucket, "backup-bucket", "", "Specify S3 bucket to copy file to")
	flag.StringVar(&backupBucket, "b", "", "Short flag - specify S3 bucket to copy file to")
	flag.StringVar(&objectPrefix, "object-prefix", "", "Prefix for S3 object uploads")
	flag.StringVar(&objectPrefix, "o", "", "Short flag - prefix for S3 object uploads")
	flag.BoolVar(&profileMemory, "profile-mem", false, "Enable memory profiling")
	flag.BoolVar(&deleteFiles, "delete", false, "Delete files from disk after upload")
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if profileMemory {
		defer profile.Start(profile.MemProfile).Stop()
	}

	err := checkFlags()
	if err != nil {
		log.Fatal(err)
	}

	ll, err := log.ParseLevel(loglevel)
	if err != nil {
		log.Error(err)
	}
	log.SetFormatter(&log.TextFormatter{
		DisableLevelTruncation: true,
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
	log.SetLevel(ll)
	log.SetOutput(os.Stdout)
	log.Debugf("Go binary built from commit: %v ", gitCommit)
	log.Debugf("The source code can be found here: %v ", gitUrl)

	files, err := findFile(filepattern)
	if err != nil {
		log.Error(err)
	}
	if files != nil {
		log.Debug("Found ", len(files), " files that matched provided pattern: ", files)
	} else {
		log.Fatalf("No files found matching pattern: %v", filepattern)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	for _, f := range files {
		err := uploadToS3(f, sess)
		if err == nil && deleteFiles {
			err := os.Remove(f)
			log.Debugf("Deleting file - %v - from disk", f)
			if err != nil {
				log.Error(err)
			}
		} else if err != nil {
			log.Error(err)
		}
	}
}

func checkFlags() error {
	if len(filepattern) == 0 {
		return errors.New("no file path provided. Please provide a file path with -filepath / -f")
	}

	if !filepath.IsAbs(filepattern) {
		return errors.New("detected relative path. please provide absolute path")
	}

	if len(backupBucket) == 0 {
		return errors.New("no S3 bucket provided. Please provide the name of a target S3 bucket with -backup-bucket / -b")
	}
	return nil
}

func findFile(pattern string) (match []string, err error) {
	fileMatches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, error(err)
	}
	return fileMatches, nil
}

func uploadToS3(f string, s client.ConfigProvider) (err error) {
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	defer file.Close()

	uploader := s3manager.NewUploader(s)
	log.Debugf("starting upload of file: %v", f)
	r, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(backupBucket),
		Key:    aws.String(objectPrefix + filepath.Base(f)),
		Body:   file,
	})
	if err != nil {
		return err
	}
	log.Infof("successfully uploaded %v to %v", f, r.Location)
	return nil
}
