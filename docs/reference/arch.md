# architecture

The kdoctor is aimed to test the cluster and generate task reports to check whether the cluster is healthy.

It consists of controller deployment and agent daemonset.

* the controller schedules the task, update and summarize the task result, and aggerate all reports.

* the agent implement tasks
