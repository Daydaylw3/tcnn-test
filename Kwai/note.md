xxxx

## 7/18

### 16

调低 x86 的 L3 cache

```
mount -t resctrl resctrl /sys/fs/resctrl
mkdir /sys/fs/resctrl/limited_group
cd /sys/fs/resctrl/limited_group
cat schemata
    L3:0=7ff;1=7ff
    MB:0=100;1=100
echo "L3:0=0ff;1=0ff" >> /sys/fs/resctrl/limited_group/schemata

```

## 修改 JVM 参数 && 重启

**查看参数**

```
echo $CLOUD_CMD | tr ' ' '\n'
```

然后获取到的内容，修改后，拷贝到下面

```
export CLOUD_CMD=$(echo '-server
-Dfile.encoding=UTF-8
...' | tr '\n' ' ')
```

**重启**

```
> ps -ef
root        7111       1  0 Jul18 ?        00:00:46 kcsize api
root        7126    7111  9 Jul18 ?        01:48:59 /home/web_server/jdk/bin/java ...

> kill -9 7111 7126
> 
> nohup kcsize api &
```
