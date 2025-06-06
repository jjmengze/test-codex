provider "aws" {
  region = "ap-northeast-1"  # 你可以依需求修改區域
}

resource "aws_kinesis_stream" "example" {
  name             = "jason-orientation-test"
  shard_count = 1  # 一般測試用途 1 shard 就夠了
  retention_period = 24 # 保留 24 小時（預設）

  stream_mode_details {
    stream_mode = "PROVISIONED"  # 或改成 ON_DEMAND（自動擴縮）
  }

  tags = {
    Environment = "dev"
    Project     = "orientation"
  }
}