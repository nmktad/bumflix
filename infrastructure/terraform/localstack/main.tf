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

variable "aws_s3_bucket_name" {
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

resource "aws_s3_bucket" "bumflix-bucket" {
  bucket = var.aws_s3_bucket_name
}
