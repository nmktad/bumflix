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
  type    = string
  default = ""
}

variable "aws_secret_access_key" {
  type    = string
  default = ""
}

variable "aws_default_region" {
  type    = string
  default = "us-east-1"
}

variable "aws_s3_hls_bucket_name" {
  type = string
}

variable "frontend_url" {
  type    = string
  default = "http://localhost:3000"
}

provider "aws" {
  region = var.aws_default_region

  access_key = var.aws_access_key
  secret_key = var.aws_secret_access_key
}

resource "aws_s3_bucket" "bumflix_hls" {
  bucket = var.aws_s3_hls_bucket_name
}

resource "aws_s3_bucket_policy" "hls_policy" {
  bucket = aws_s3_bucket.bumflix_hls.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "AllowPublicRead"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.bumflix_hls.arn}/*"
      }
    ]
  })
}

resource "aws_s3_bucket_cors_configuration" "hls_cors" {
  bucket = aws_s3_bucket.bumflix_hls.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = [var.frontend_url]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
