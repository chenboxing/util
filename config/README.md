config
======

Package config simplifies the code for an application to handle configurations.

The config file may contain an entry to configure `gitlab.ulaiber.com/uchang/util/log`

write log to redis:

```yaml
log:
  formatter: logstash   #日志格式
  release:   0.1        #应用程序版本
  mode:      develop    #环境
  level:     "debug"    #日志等级
  port:      12345      #日志tcp端口，供动态修改日志等级，不配置则随机
  writer:               #日志输出配置
    type:  redis
    redis:
      addr: 127.0.0.1:6379
      password:
      list_key: config_test_log
```

write log to file:

```yaml
log:
  formatter: logstash
  release:   0.1
  mode:      develop
  level:     "debug"
  port:      12345
  writer:
    type:  file
    path:  demo.log
```

default stdout log.

## Change log level

```shell
echo -n "help" | nc -4t  -w1 127.0.0.1 8765
```
