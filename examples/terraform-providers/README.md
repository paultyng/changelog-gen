# Generate Terraform Provider Changelog

```shell
$ changelog-gen \
  -owner terraform-providers \
  -repo terraform-provider-aws \
  -changelog examples/terraform-providers/changelog.tmpl \
  -releasenote examples/terraform-providers/release-note.tmpl \
  441ec74e66706cc0a75d4d207724cd6460f5f6a4 \
  f3bdfeaaa7ddd9522c549e9f134948e0d698ab02
```

Example Output:

---
FEATURES

* **costandusagereportservice:** Add Resource for Cost and Usage Report Definitions ([7432](https://github.com/terraform-providers/terraform-provider-aws/pull/7432) by [jbmchuck](https://github.com/jbmchuck))
* **iot:** resource/aws_iot_role_alias: Add initial support for IoT role aliases ([7348](https://github.com/terraform-providers/terraform-provider-aws/pull/7348) by [jhosteny](https://github.com/jhosteny))
* **ram:** New Resource: aws_ram_principal_association ([7563](https://github.com/terraform-providers/terraform-provider-aws/pull/7563) by [bflad](https://github.com/bflad))
* **ram:** New Resource: aws_ram_resource_association ([7449](https://github.com/terraform-providers/terraform-provider-aws/pull/7449) by [bflad](https://github.com/bflad))

IMPROVEMENTS

* **acmpca:** tests/data-source/aws_acmpca_certificate_authority: Use resource.TestCheckResourceAttrPair() ([7540](https://github.com/terraform-providers/terraform-provider-aws/pull/7540) by [bflad](https://github.com/bflad))
* **apigateway:** resource/aws_api_gateway_rest_api: Remove extraneous timeout for resource deletion ([7554](https://github.com/terraform-providers/terraform-provider-aws/pull/7554) by [bflad](https://github.com/bflad))
* **appmesh:** r/aws_appmesh_virtual_node: Add support for listener health checks ([7446](https://github.com/terraform-providers/terraform-provider-aws/pull/7446) by [ewbankkit](https://github.com/ewbankkit))
* **appsync:** tests/resource/aws_appsync_graphql_api: Add sweeper ([7538](https://github.com/terraform-providers/terraform-provider-aws/pull/7538) by [bflad](https://github.com/bflad))
* **cloudwatch:** resource/aws_cloudwatch_metric_alarm: Ensure dimensions configurations use equals (part deux) ([7543](https://github.com/terraform-providers/terraform-provider-aws/pull/7543) by [bflad](https://github.com/bflad))
* **cognito:** Add support for Cognito user pool advanced security mode ([7361](https://github.com/terraform-providers/terraform-provider-aws/pull/7361) by [tyrjola](https://github.com/tyrjola))
* **directoryservice:** Issue #7466 Add security_group_id attribute to aws_directory_service_directory for ADConnector ([7487](https://github.com/terraform-providers/terraform-provider-aws/pull/7487) by [saravanan30erd](https://github.com/saravanan30erd))
* **docdb:** tests/service/docdb: Temporarily use expanded references ([7545](https://github.com/terraform-providers/terraform-provider-aws/pull/7545) by [bflad](https://github.com/bflad))
* **dynamodb:** Update aws_dynamodb_table.server_side_encryption documentation noting DEFAULT encryption type ([7499](https://github.com/terraform-providers/terraform-provider-aws/pull/7499) by [ewbankkit](https://github.com/ewbankkit))
* **dynamodb:** d/aws_dynamodb_table: Add missing 'billing_mode' and 'point_in_time_recovery' attributes ([7497](https://github.com/terraform-providers/terraform-provider-aws/pull/7497) by [ewbankkit](https://github.com/ewbankkit))
* **ec2, mq, sagemaker:** tests/provider: Ensure tags configurations use equals ([7559](https://github.com/terraform-providers/terraform-provider-aws/pull/7559) by [bflad](https://github.com/bflad))
* **kinesisanalytics:** tests/resource/aws_kinesis_analytics_application: Fix TestAccAWSKinesisAnalyticsApplication_outputsMultiple syntax for Terraform 0.12 ([7560](https://github.com/terraform-providers/terraform-provider-aws/pull/7560) by [bflad](https://github.com/bflad))
* **rds:** aurora docs link directed to RDS docs, now aurora docs ([7530](https://github.com/terraform-providers/terraform-provider-aws/pull/7530) by [codestoe](https://github.com/codestoe))
* **route53, s3:** Enable staticcheck S1017 ([7520](https://github.com/terraform-providers/terraform-provider-aws/pull/7520) by [nywilken](https://github.com/nywilken))
* **s3:** r/aws_s3_bucket: Only check bucket existence in 'testAccCheckAWSS3BucketDestroy' ([7486](https://github.com/terraform-providers/terraform-provider-aws/pull/7486) by [ewbankkit](https://github.com/ewbankkit))
* **ses:** Add position to aws_ses_receipt_rule example ([7478](https://github.com/terraform-providers/terraform-provider-aws/pull/7478) by [kula](https://github.com/kula))
* **ses:** Add support for resource_aws_ses_identity_notification_topic import ([7343](https://github.com/terraform-providers/terraform-provider-aws/pull/7343) by [jkmart](https://github.com/jkmart))
* **waf:** resource/aws_waf_web_acl: Add logging configuration ([6059](https://github.com/terraform-providers/terraform-provider-aws/pull/6059) by [WhileLoop](https://github.com/WhileLoop))
* deps: github.com/aws/aws-sdk-go@v1.16.32 ([7511](https://github.com/terraform-providers/terraform-provider-aws/pull/7511) by [bflad](https://github.com/bflad))
* moving eks_cluster_auth link in docs ([7539](https://github.com/terraform-providers/terraform-provider-aws/pull/7539) by [DanielMolander](https://github.com/DanielMolander))
* tests/resource/aws_worklink_fleet: Temporarily expand subnet_ids references ([7556](https://github.com/terraform-providers/terraform-provider-aws/pull/7556) by [bflad](https://github.com/bflad))

BUGS

* **ec2:** Fix crash on unsuccessful flow log creation ([7528](https://github.com/terraform-providers/terraform-provider-aws/pull/7528) by [phy1729](https://github.com/phy1729))
* **ec2:** Fix errors in EC2 Transit Gateway VPC Attachment when TGW is shared from another account ([7513](https://github.com/terraform-providers/terraform-provider-aws/pull/7513) by [andrewsuperbrilliant](https://github.com/andrewsuperbrilliant))
* **iam:** Catch AccessDenied permission errors when creating aws_iam_user_login… ([7519](https://github.com/terraform-providers/terraform-provider-aws/pull/7519) by [fbreckle](https://github.com/fbreckle))
* **iot:** resource/aws_iot_topic_rule: Fix optional attributes ([7471](https://github.com/terraform-providers/terraform-provider-aws/pull/7471) by [nywilken](https://github.com/nywilken))
* **kinesisanalytics:** r/kinesis_analytics_application: support multiple outputs ([7535](https://github.com/terraform-providers/terraform-provider-aws/pull/7535) by [kl4w](https://github.com/kl4w))
* **rds:** Fixing global rds by allowing for optional parameters on cluster ([7213](https://github.com/terraform-providers/terraform-provider-aws/pull/7213) by [bculberson](https://github.com/bculberson))
* **ssm:** r/aws_ssm_maintenance_window_task: Fix name/desc validation ([7186](https://github.com/terraform-providers/terraform-provider-aws/pull/7186) by [YakDriver](https://github.com/YakDriver))
