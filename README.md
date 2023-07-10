# SQS to SES Lambda Function

This is the code used for https://tituscoleman.dev on the contact form. The intention is to take the form data JSON submitted to AWS Simple Queue Service from my website, format it into a basic email then have AWS Simple Email Service email myself directly.

## Written in GO using AWS-SDK GO

The logic behind this was that I wanted to avoid needing a database for this simple task as my email account would essentially keep this record. GO was selected due to its fast execution time, small binary and simplicity hereby saving on Lambda execution cost.
