mvp needs:
[x] takes input for filename / pattern
[x] validate input
[x] takes input for destination bucket
[x] uploads file to bucket

nice to haves:
[x] boolean for delete from disk after backup
[ ] todo: compress prior to upload if not already compressed
[x] publishes metric to cloud watch metrics
[ ] todo: no-op / dry run mode - logs what would have happened
[x]optionally prefix new s3 objects (e.g. today's date)
[x]flag for profiling