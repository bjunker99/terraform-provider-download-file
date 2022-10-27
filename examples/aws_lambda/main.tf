terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.36.1"
    }
    download = {
      version = "~> 0.0.1"
      source  = "github.com/bjunker99/download-file"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

data "download_file" "release" {
  url         = "https://github.com/bjunker99/LambdaHelloWorld/releases/download/v${var.release_version}/LambdaHelloWorld.zip"
  output_file = "LambdaHelloWorld.zip"

  verify_sha256 = "62c8cc77965f9a367a07a92134ec7fb9ba15f04de5c4686af9cfb1fc9f586625"
}

resource "aws_lambda_function" "hello_world_lambda" {
  architectures = [
    "x86_64",
  ]
  function_name                  = "LambdaHelloWorld"
  handler                        = "LambdaHelloWorld::LambdaHelloWorld.Function::FunctionHandler"
  layers                         = []
  memory_size                    = 256
  package_type                   = "Zip"
  reserved_concurrent_executions = -1
  role                           = aws_iam_role.lambda_role.arn
  runtime                        = "dotnet6"
  filename                       = data.download_file.release.output_file
  source_code_hash               = data.download_file.release.output_base64sha256
  timeout                        = 30

  tracing_config {
    mode = "PassThrough"
  }
}

resource "aws_iam_role" "lambda_role" {
  assume_role_policy = jsonencode(
    {
      Statement = [
        {
          Action = "sts:AssumeRole"
          Effect = "Allow"
          Principal = {
            Service = "lambda.amazonaws.com"
          }
        },
      ]
      Version = "2012-10-17"
    }
  )
  force_detach_policies = false
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
  ]
  max_session_duration = 3600
  name                 = "HelloWorldLambda-Role-${var.aws_region}"
  path                 = "/"
}
