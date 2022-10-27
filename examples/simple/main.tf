terraform {
  required_providers {
    download = {
      version = "~> 0.0.1"
      source  = "github.com/bjunker99/download-file"
    }
  }
}

data "download_file" "test" {
  url           = "https://github.com/bjunker99/LambdaHelloWorld/releases/download/v1.1/LambdaHelloWorld.zip"
  output_file   = "LambdaHelloWorld.zip"
}
