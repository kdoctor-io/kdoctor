# Architecture

The kdoctor aims to test clusters and generate task reports to check whether the cluster is healthy.

It consists of controller deployment and agent daemonset.

* The controller schedules the task, update and summarize the task result, and aggregate all reports.

* The agent implement tasks.
