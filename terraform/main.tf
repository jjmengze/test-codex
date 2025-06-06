provider "aws" {
  region = "us-west-2"
}

resource "aws_kinesis_stream" "example" {
  name             = "jason-orientation-test"
  shard_count      = 1
  retention_period = 24

  stream_mode_details {
    stream_mode = "PROVISIONED"  # 或改成 ON_DEMAND（自動擴縮）
  }

  tags = {
    Environment = "dev"
    Project     = "orientation"
  }
}