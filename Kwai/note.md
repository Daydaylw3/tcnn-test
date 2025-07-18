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
