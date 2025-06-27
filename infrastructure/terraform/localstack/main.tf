terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

variable "aws_access_key" {
  type = string
}

variable "aws_secret_access_key" {
  type = string
}

variable "aws_default_region" {
  type = string
}

variable "aws_s3_localstack_endpoint" {
  type = string
}

variable "aws_s3_raw_bucket_name" {
  type = string
}

variable "aws_s3_hls_bucket_name" {
  type = string
}

provider "aws" {

  access_key = var.aws_access_key
  secret_key = var.aws_secret_access_key
  region     = var.aws_default_region

  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    s3 = var.aws_s3_localstack_endpoint
  }
}

# resource "aws_sqs_queue" "transcode-queue" {
#   name = "bumflix-transcode-queue"
# }
#
# output "sqs_queue_url" {
#   value = aws_sqs_queue.transcode-queue.id
# }
#
resource "aws_s3_bucket" "bumflix-hls" {
  bucket = var.aws_s3_hls_bucket_name
}

resource "aws_s3_bucket" "bumflix-raw" {
  bucket = var.aws_s3_raw_bucket_name
}

resource "aws_s3_bucket_policy" "hls_policy" {
  bucket = aws_s3_bucket.bumflix-hls.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "AllowPublicRead"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.bumflix-hls.arn}/*"
      }
    ]
  })
}

resource "aws_s3_bucket_cors_configuration" "hls_cors" {
  bucket = aws_s3_bucket.bumflix-hls.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
#
# resource "aws_s3_bucket_notification" "raw_video_upload" {
#   bucket = aws_s3_bucket.bumflix_raw.id
#
#   queue {
#     queue_arn = aws_sqs_queue.transcode_queue.arn
#     events    = ["s3:ObjectCreated:*"]
#   }
#
#   depends_on = [aws_sqs_queue_policy.allow_s3_to_sqs]
# }
