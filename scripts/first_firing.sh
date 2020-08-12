#!/bin/bash

firing='
{"receiver":"webhook","status":"firing","alerts":[{"status":"firing","labels":{"alertname":"SpringAuthVpnDown","instance":"172.18.0.131:8000","job":"kubernetes-service-endpoints","kubernetes_name":"prometheus-smp-exporter","kubernetes_namespace":"morefun-pro","severity":"Critical","site":"172.18.1.249:6100"},"annotations":{"description":"spring auth vpn site 172.18.1.249:6100 check failed"},"startsAt":"2019-07-24T03:06:01.192911495Z","endsAt":"0001-01-01T00:00:00Z","generatorURL":"https://dev-prometheus.i-morefun.net/graph?g0.expr=spring_auth_vpn_up+%21%3D+1\u0026g0.tab=1"},{"status":"firing","labels":{"alertname":"SpringAuthVpnDown","instance":"172.18.0.131:8000","job":"kubernetes-service-endpoints","kubernetes_name":"prometheus-smp-exporter","kubernetes_namespace":"morefun-pro","severity":"Critical","site":"172.18.20.85:6100"},"annotations":{"description":"spring auth vpn site 172.18.20.85:6100 check failed"},"startsAt":"2019-07-24T03:06:01.192911495Z","endsAt":"0001-01-01T00:00:00Z","generatorURL":"https://dev-prometheus.i-morefun.net/graph?g0.expr=spring_auth_vpn_up+%21%3D+1\u0026g0.tab=1"}],"groupLabels":{"alertname":"SpringAuthVpnDown"},"commonLabels":{"alertname":"SpringAuthVpnDown","instance":"172.18.0.131:8000","job":"kubernetes-service-endpoints","kubernetes_name":"prometheus-smp-exporter","kubernetes_namespace":"morefun-pro","severity":"Critical"},"commonAnnotations":{},"externalURL":"https://dev-alertmanager.i-morefun.net","version":"4","groupKey":"{}:{alertname=\"SpringAuthVpnDown\"}"}
'

curl -XPOST -d"$firing" http://127.0.0.1:8000/dingtalk/alertmanager
