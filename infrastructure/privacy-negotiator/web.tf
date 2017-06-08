provider "aws" {
  alias  = "us-east-1"
  region = "us-east-1"
}

resource "aws_s3_bucket" "web" {
  bucket = "${var.domain}"
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "index.html"
  }

  tags {
    Name = "Privacy Negotiator Web"
  }
}

resource "aws_cloudfront_origin_access_identity" "web" {
  comment = "Web origin access identity"
}

resource "aws_cloudfront_distribution" "web_distribution" {

  aliases = ["${var.domain}"]

  comment             = "Privacy Negotiator Web"
  default_root_object = "index.html"
  enabled             = true
  is_ipv6_enabled     = true

  origin {
    domain_name = "${aws_s3_bucket.web.website_endpoint}"
    origin_id   = "privacy-negotiator-web-s3"

    custom_origin_config {
      http_port = 80
      https_port = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols = ["TLSv1.2"]
    }
  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = "privacy-negotiator-web-s3"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  tags {
    Environment = "${var.environment}"
  }

  viewer_certificate {
    acm_certificate_arn = "${data.aws_acm_certificate.web.arn}"
    ssl_support_method = "sni-only"
    minimum_protocol_version = "TLSv1"
  }
}

resource "aws_route53_record" "web" {
  zone_id = "${data.aws_route53_zone.organisation.zone_id}"
  name    = "${var.domain}"
  type    = "A"

  alias {
    name                   = "${aws_cloudfront_distribution.web_distribution.domain_name}"
    zone_id                = "${aws_cloudfront_distribution.web_distribution.hosted_zone_id}"
    evaluate_target_health = true
  }
}

data "aws_acm_certificate" "web" {
  provider = "aws.us-east-1"
  domain   = "${var.domain}"
  statuses = ["ISSUED"]
}
