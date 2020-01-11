[![Go Report](https://goreportcard.com/badge/github.com/jhughes01/s3-file-backup)](https://goreportcard.com/badge/github.com/jhughes01/s3-file-backup)

# s3-file-backup

s3-file-backup is a simple utility program written in go. 

## Usage

### Example

All configuration is passed via command line flags.

The minimum flags required are `-filepath` and `-backup-bucket`.

    ./s3-file-backup -f /home/jhughes/myImportantFile.txt -b myS3Bucket
    
### Flags

Flag | Short flag | Description | Default | Required? 
-----|------------|-------------|---------|----------
-filepath | -f | Absolute path for file to be uploaded to S3 | nil | Y
-backup-bucket | -b | Name of S3 bucket to upload files to | nil | Y
-log-level | -l | Log level. Options: Trace, Debug, Error or Info | Info | N
-object-prefix | -o | Object prefix for files when uploaded to S3 | nil | N
-profile-mem | N/A | Enable Golang memory profiling | false | N
-delete | N/A | Delete files from local disk after successful upload to S3 | false | N
-cloudwatch-metric | N/A | Publish cloudwatch metric with result of S3 upload (0 or 1) | false | N