/*
mvp needs:
takes input for filename / pattern
todo: takes input for destination bucket

nice to haves:
boolean for delete from disk after backup
compress prior to upload if not already compressed
sends event to Cloud watch events with exit status (success, error)
*/

package main

import (
	"flag"
	"os"
	"path/filepath"

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
	}
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

func uploadToS3() {
	log.Info("S3")
}

func findFile(pattern string) (match []string, err error) {
	fileMatches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, error(err)
	}
	return fileMatches, nil
}
