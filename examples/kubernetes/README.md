# Generate Kubernetes Changelog

```shell
$ changelog-gen \
  -owner kubernetes \
  -repo kubernetes \
  -changelog examples/kubernetes/changelog.tmpl \
  -releasenote examples/kubernetes/release-note.tmpl \
  2e90d92db9a807fcc140977e7a4798c6078014c2 \
  a1539747db15b01e4b6baab6fb505a31e30c05ef
```

Example Output:

---
IMPROVEMENTS

* **apiserver:** Adding a limit on the size of request body the apiserver will decode for write operations ([73805](https://github.com/kubernetes/kubernetes/pull/73805) by [caesarxuchao](https://github.com/caesarxuchao))
* **apiserver:** Remove storage versions flag ([67678](https://github.com/kubernetes/kubernetes/pull/67678) by [caesarxuchao](https://github.com/caesarxuchao))
* **kubectl:** add kustomize as a subcommand in kubectl ([73033](https://github.com/kubernetes/kubernetes/pull/73033) by [Liujingfang1](https://github.com/Liujingfang1))
* **kubelet:** Update kubelet overview help doc ([73256](https://github.com/kubernetes/kubernetes/pull/73256) by [deitch](https://github.com/deitch))
* **kubelet:** kubelet: promote OS &amp;amp; arch labels to GA ([73333](https://github.com/kubernetes/kubernetes/pull/73333) by [yujuhong](https://github.com/yujuhong))
* fix typo ([73898](https://github.com/kubernetes/kubernetes/pull/73898) by [xiezongzhe](https://github.com/xiezongzhe))

BUGS

* **apiserver, kubeadm:** update the dependency pmezard/go-difflib ([73941](https://github.com/kubernetes/kubernetes/pull/73941) by [neolit123](https://github.com/neolit123))
* **apiserver:** openapi-aggregation: speed up merging from 1 sec to 50-100 ms ([71223](https://github.com/kubernetes/kubernetes/pull/71223) by [sttts](https://github.com/sttts))
* **kubeadm:** kubeadm: add a preflight check for Docker and cgroup driver ([73837](https://github.com/kubernetes/kubernetes/pull/73837) by [neolit123](https://github.com/neolit123))
* **kubelet:** Kubelet: add usageNanoCores from CRI stats provider ([73659](https://github.com/kubernetes/kubernetes/pull/73659) by [feiskyer](https://github.com/feiskyer))
* **kubelet:** Make container create, start, and stop events consistent ([73892](https://github.com/kubernetes/kubernetes/pull/73892) by [smarterclayton](https://github.com/smarterclayton))
* **kubelet:** cpuPeriod was not reset ([73342](https://github.com/kubernetes/kubernetes/pull/73342) by [szuecs](https://github.com/szuecs))
---