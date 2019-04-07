# mackerel-remora
Remora is agent for Mackerel ( mackerel.io ) to post service metrics

![remora](https://user-images.githubusercontent.com/1097533/55680852-4b7e2e00-595a-11e9-88c9-8624f3ff5332.png)

## features

- Post service metrics


## How to use
For example, using `docker-compose` .

```yml
remora:
  image: aknow/mackerel-remora:latest
  volumes:
      - .:/app
  environment:
    MACKEREL_REMORA_CONFIG: /app/_example/sample.yml
```

`sample.yml` is a config for `mackerel-remora` . The following is an example of its contents.

```yml
apikey: xxxxx
plugin:
  servicemetrics:
    demo:
      sample:
        command: sh /app/_example/demo.sh
```

An example of the contents of `/app/_example/demo.sh` is:

```sh
#!/bin/sh

metric_name="demo.test_metric.number"
# metric number to post
metric=6
date=`date +%s`

echo "${metric_name}\t${metric}\t${date}"
```


If you execute `docker-compose up -d remora` and wait for a while, it looks like on Mackerel:

<img width="882" src="https://user-images.githubusercontent.com/1097533/55680926-3b1a8300-595b-11e9-8576-9ba525750763.png">


## Config file format

```yml
apibase: <apibase>
apikey: <apikey>
root: <rootdir>
plugin:
  servicemetrics:
    <servicename>:
      <settingname>:
        command: <command>
        user: <user>
        env:
          <envkey>: <envvalue>
```


- `<apibase>`
    - You can specify the domain of request destination API. It is usually unnecessary to specify.
- `<apikey>`
    - API key issued by the Mackerel organization to which you are posting. Requires write permission.
- `<root>`
    - Root directory of mackerel-remora. It is usually unnecessary to specify.
- `<servicename>`
    - Destination service name. You need to create a service on Mackerel beforehand.
- `<settingname>`
    - Plugin setting name for posting service metrics. It should not be used, but should not be duplicated.
- `<command>`
    - A command to perform standard output expected as a service metric plug-in.
    - Remora expect standard output in the format `{metric name}\t{metric value}\t{epoch seconds}` as a result of executing this command.
- `<user>`
    - User of `command` running.
- `<envkey>`
    - Variable name of the environment variable to pass to the `command` .
- `<envvalue>`
    - Value of the environment variable to pass to the `command` .
