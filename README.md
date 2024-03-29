# firewall-policy-api
防火墙策略接口，通过Grafana Alerting/API 添加防火墙规则，实现自定义的拦截外网扫描功能。只测试并实现了H3C

## Grafana中配置
  - 在grafana中添加`Alerting`-> `Contact points`
  - 选择webhook
  - URL填入`http://your_docker_ip:80/banByGrafana`
  - 支持临时封禁`http://your_docker_ip:80/banByGrafana?temporaryBan=true`
  - HTTP Method `POST`
  - 创建一个`Alert rules`, 在`Description`中填入
    ```
    {{- range $k, $v := $values -}}
    IP: {{ $v.Labels.clientHost }}

    {{ end -}}
    ```
  - 创建一个`Notification policies`, 将`Alert rules`与`Contact points`通过标签关联

## 其它设置
  - SSH_USER `用户`
  - SSH_PASSWORD `密码`
  - SSH_IP_PORT `防火墙ip:port`
  - IP_WHITE_LIST `IP白名单(通过逗号分割) 1.1.1.1,2.2.2.0/24`

## docker中使用
- 运行docker容器
  ```
  docker run -d --name firewall-policy-api --restart=always \
    -p 80:80 \
    -e SSH_USER=admin \
    -e SSH_PASSWORD=your_password \
    -e SSH_IP_PORT=192.168.0.1:22 \
    -e IP_WHITE_LIST=1.1.1.1,2.2.2.0/24 \
    jeessy/firewall-policy-api
  ```

## 系统中使用
- 下载并解压[https://github.com/jeessy2/firewall-policy-api/releases](https://github.com/jeessy2/firewall-policy-api/releases)
- 通过接口 `curl http://your_docker_ip:80/banByIP?ip=1.1.1.1` 配置
