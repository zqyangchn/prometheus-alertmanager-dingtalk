# prometheus-alertmanager-dingtalk
##### prometheus alertmanager webhook to dingtalk(钉钉群机器人告警)

## 关于配置文件 .example-config/config.yaml.example
	配置钉钉群机器人 access_token
		dingTalkUri: https://oapi.dingtalk.com/robot/send?access_token=XXXXXXX
	配置钉钉群机器人安全设置
		securitySettingsType: Endorsement
		    CustomKeywords		最多可以设置10个关键词，消息中至少包含其中1个关键词才可以发送成功
		    Endorsement			加签, 需要提供 SecretKey
		    IPAddressSegment		只有来自IP地址范围内的请求才会被正常处理
		secretKey: SECXXXXXXX

## kuberneter deploy
    ..example-deploy-kubernetes

## 运行
    docker run -p 50000:8000 -v /path/to/config.yaml:/etc/prometheus-alertmanager-dingtalk/config.yaml rewind/prometheus-alertmanager-dingtalk:0.0.6

## 构建镜像
	docker build -f Dockerfile -t rewind/prometheus-alertmanager-dingtalk:0.0.6 . 
	docker build -f Dockerfile -t reg.i-morefun.net/google_containers/prometheus-alertmanager-dingtalk:0.0.6 . 

## 二进制编译
	./build mac|linux
