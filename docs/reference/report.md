## Agent Report

When agent finish task, it saves report to `/report` with name fmt.Sprintf("%s_%s_round%d_%s_%s", kindName, taskName, roundNumber, nodeName, suffix).
The report will be automatically deleted with age `spec.schedulePlan.TimeoutMinute + 5` minutes. In this interval, 
the controller Pod will collect this report and save to `/report` of controller Pod.

## Controller Report

When task finishes, it saves report to `/report` with name fmt.Sprintf("%s_%s_round%d_%s_%s", kindName, taskName, roundNumber, nodeName, suffix).
It also collects all agent report and  saves report to `/report`.
All files in `/report` of controller will survive with max age maxAgeInDay(default 30 days). It could be adjusted in the configmap.

The controller could save reports to host path or PVC.
