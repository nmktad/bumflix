terraform {
  required_providers {
    aws     = { source = "hashicorp/aws" }
  }
}

provider "aws" {
  access_key                  = "mock_access_key"
  secret_key                  = "mock_secret_key"
  region                      = "eu-west-3"

  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    s3             = "http://s3.localhost.localstack.cloud:4566"
    sqs            = "http://localhost:4566"
  }
}

resource "aws_s3_bucket" "bumflix_raw_films" {
  bucket = "bumflix-raw-films"
}

resource "aws_s3_bucket" "bumflix_hls_segments" {
  bucket = "bumflix-hls-segments"
}

resource "aws_sqs_queue" "bumflix_transcode" {
  name = "bumflix-transcode"
}

resource "aws_sqs_queue" "bumflix_transcode_dlq" {
  name = "bumflix-transcode-dlq"
  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns   = [aws_sqs_queue.bumflix_transcode.arn]
  })
}

resource "aws_sqs_queue_redrive_policy" "bumflix_transcode" {
  queue_url = aws_sqs_queue.bumflix_transcode.id
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.bumflix_transcode_dlq.arn
    maxReceiveCount     = 4
  })
}
