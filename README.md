[![Build Status](https://travis-ci.org/jhughes01/s3-file-backup.svg?branch=master)](https://travis-ci.org/jhughes01/s3-file-backup)
[![Go Report](https://goreportcard.com/badge/github.com/jhughes01/s3-file-backup)](https://goreportcard.com/report/github.com/jhughes01/s3-file-backup)

# s3-file-backup

s3-file-backup is a simple utility program written in go. 

## Usage

### Example

All configuration is passed via command line flags.

The minimum flags required are `-filepath` and `-backup-bucket`.

    ./s3-file-backup -f /home/jhughes/myImportantFile.txt -b myS3Bucket
    
### Flags

Flag | Short flag | Type | Description | Default | Required? 
-----|------------|------|-------------|---------|-----------
-filepath | -f | string | Absolute path for file to be uploaded to S3 | nil | Y
-backup-bucket | -b | string | Name of S3 bucket to upload files to | nil | Y
-log-level | -l | string | Log level. Options: Trace, Debug, Error or Info | Info | N
-object-prefix | -o | string | Object prefix for files when uploaded to S3 | nil | N
-profile-mem | N/A | bool | Enable Golang memory profiling | false | N
-delete | N/A | bool | Delete files from local disk after successful upload to S3 | false | N
-cloudwatch-metric | string | N/A | Publish cloudwatch metric with result of S3 upload (0 or 1) | false | N

### Globbing

The `-filepath` parameter supports globbing and can be used to return multiple files matching a pattern.

    term:
    	'*'         matches any sequence of non-Separator characters
    	'?'         matches any single non-Separator character
    	'[' [ '^' ] { character-range } ']'
    	            character class (must be non-empty)
    	c           matches character c (c != '*', '?', '\\', '[')
    	'\\' c      matches character c
    
    character-range:
    	c           matches character c (c != '\\', '-', ']')
    	'\\' c      matches character c
    	lo '-' hi   matches character c for lo <= c <= hi