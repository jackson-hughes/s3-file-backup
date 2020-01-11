package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/profile"

	log "github.com/sirupsen/logrus"
)

var (
	gitCommit        string
	gitUrl           string
	filepattern      string
	loglevel         string
	backupBucket     string
	objectPrefix     string
	cloudwatchMetric string
	profileMemory    bool
	deleteFiles      bool
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
	flag.StringVar(&cloudwatchMetric, "cloudwatch-metric", "", "Cloudwatch metric to publish to")
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

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	files, err := findFile(filepattern)
	if err != nil {
		log.Error(err)
	}
	if files != nil {
		log.Debug("Found ", len(files), " files that matched provided pattern: ", files)
	} else {
		log.Errorf("No files found matching pattern: %v", filepattern)
		if cloudwatchMetric != "" {
			err := publishMetric(1, sess)
			if err != nil {
				log.Errorf("error publishing metric: %v", err)
			}
		}
	}

	var successList []string
	for _, f := range files {
		err := uploadToS3(f, sess)
		if err != nil {
			log.Error("error uploading file to S3: ", err)
			continue
		}
		successList = append(successList, f)
	}
	allSuccess := reflect.DeepEqual(files, successList)
	if allSuccess {
		if cloudwatchMetric != "" {
			err := publishMetric(0, sess)
			if err != nil {
				log.Error(err)
			}
		}
	} else {
		if cloudwatchMetric != "" {
			err := publishMetric(1, sess)
			if err != nil {
				log.Error(err)
			}
		}
	}
	if deleteFiles && allSuccess {
		for _, f := range files {
			err := os.Remove(f)
			log.Debugf("deleting %v from local disk ", f)
			if err != nil {
				log.Error(err)
			}
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

func publishMetric(r float64, s client.ConfigProvider) error {
	log.Debugf("publishing cloud watch metric with value: %v", r)
	svc := cloudwatch.New(s)
	_, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String(cloudwatchMetric),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String("S3BackupStatus"),
				Value:      aws.Float64(r),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
