# Generate Changelog From Typed Blocks

PRs have release-note blocks with types, not labels indicating the type of each
entry. For example, the PRs in the BUGS section all have the following block:

```release-note:bug
The entry goes here
```

```shell
$ changelog-gen \
  -owner terraform-providers \
  -repo terraform-provider-google \
  -changelog examples/typed-blocks/changelog.tmpl \
  -releasenote examples/typed-blocks/release-note.tmpl \
  -no-note-label "changelog: no-release-note" \
  14729f9457253f2ac24b602c1f4d80621431c169 \
  806a31474d6100d1691c120b3d85bea94cbcfff9
```

Example Output:

---
DEPRECATIONS:
* `container`: The `kubernetes_dashboard` addon is deprecated for `google_container_cluster`. ([#4648](https://github.com/terraform-providers/terraform-provider-google/pull/4648))

FEATURES:
* **New Resource:** `google_app_engine_application_url_dispatch_rules` ([#4674](https://github.com/terraform-providers/terraform-provider-google/pull/4674))

IMPROVEMENTS:
* `all`: increased support for custom endpoints across the provider ([#4641](https://github.com/terraform-providers/terraform-provider-google/pull/4641))
* `appengine`: added the ability to delete the parent service of `google_app_engine_standard_app_version` ([#4596](https://github.com/terraform-providers/terraform-provider-google/pull/4596))
* `container`: Added `shielded_instance_config` attribute to `node_config` ([#4554](https://github.com/terraform-providers/terraform-provider-google/pull/4554))
* `dataflow`: added `ip_configuration` option to `job`. ([#4726](https://github.com/terraform-providers/terraform-provider-google/pull/4726))
* `pubsub`: Added field `oidc_token` to `google_pubsub_subscription` ([#4679](https://github.com/terraform-providers/terraform-provider-google/pull/4679))
* `sql`: added `location` field to `backup_configuration` block in `google_sql_database_instance` ([#4681](https://github.com/terraform-providers/terraform-provider-google/pull/4681))

BUGS:
* Clarified error on invalid input to lists of references. ([#4731](https://github.com/terraform-providers/terraform-provider-google/pull/4731))
* Fixed `google_compute_router` failing to remove BGP information. ([#4742](https://github.com/terraform-providers/terraform-provider-google/pull/4742))
* `all`: fixed the custom endpoint version used by older legacy REST clients ([#4695](https://github.com/terraform-providers/terraform-provider-google/pull/4695))
* `bigquery`: fix issue with `google_bigquery_data_transfer_config` `params` crashing on boolean values ([#4676](https://github.com/terraform-providers/terraform-provider-google/pull/4676))
* `cloudrun`: fixed the apiVersion sent in `google_cloud_run_domain_mapping` requests ([#4657](https://github.com/terraform-providers/terraform-provider-google/pull/4657))
* `compute`: added support for updating multiple fields at once to `google_compute_subnetwork` ([#4688](https://github.com/terraform-providers/terraform-provider-google/pull/4688))
* `compute`: fixed diffs in `google_compute_instance_group`'s `network` field when equivalent values were specified ([#4728](https://github.com/terraform-providers/terraform-provider-google/pull/4728))
* `compute`: fixed issues updating `google_compute_instance_group`'s `instances` field when config/state values didn't match ([#4728](https://github.com/terraform-providers/terraform-provider-google/pull/4728))
* `iam`: fixed bug where IAM binding wouldn't replace members if they were deleted outside of terraform. ([#4693](https://github.com/terraform-providers/terraform-provider-google/pull/4693))
* `pubsub`: Fixed permadiff due to interaction of organization policies and `google_pubsub_topic`. ([#4721](https://github.com/terraform-providers/terraform-provider-google/pull/4721))

