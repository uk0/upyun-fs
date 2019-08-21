## Quick Start 


#### Umount
```bash
umount -f  upyun-fs
```


#### Test

```bash
./upyun-fs -mp /Users/zhangjianxin/home/GO_LIB/src/upyun-fs/testMountDir2  \
 -upyun_bkt test_bkt \
 -upyun_op test_op \
 -upyun_pass password
 
 
 
 
(base) ➜  testMountDir2 ls
Command.rar                   apache                        flink.png                     mac_soft                      tensorflow
LoadingScreen.mp4             assets                        free_dir                      main.ps1                      unnamed.png
OPUI.png                      blog                          gaokeyong.png                 markdown                      upx.sh
OPUI_1.png                    cdh                           images                        mongo_tar                     upyun_storage_log_AhYIBW15
PUBG_FUCK.zip                 demos                         index.html                    neo4jData.json                visio
QQ20170601-114658@2x.png      doc                           kubernetes                    posts                         weixin
QQ20170601-115231@2x.png      dockerimages                  kubernetes1.6.jpg             simple_img                    wiki
QQ20170615-214859@2x.png      flink-header-logo.svg         lib.js                        soft                          x.tar.gz
README.md                     flink-home-graphic-update.svg libs                          storm_tar                     zip 








(base) ➜  testMountDir2 df -h
Filesystem      Size   Used  Avail Capacity iused               ifree %iused  Mounted on
/dev/disk1s1   233Gi  207Gi   23Gi    91% 3598598 9223372036851177209    0%   /
devfs          245Ki  245Ki    0Bi   100%     846                   0  100%   /dev
/dev/disk1s4   233Gi  3.0Gi   23Gi    12%       6 9223372036854775801    0%   /private/var/vm
map -hosts       0Bi    0Bi    0Bi   100%       0                   0  100%   /net
map auto_home    0Bi    0Bi    0Bi   100%       0                   0  100%   /home
upyun-fs         0Bi    0Bi    0Bi   100%       0                   0  100%   /Users/zhangjianxin/home/GO_LIB/src/upyun-fs/testMountDir2


```



#### TODO


* 还在写文件读写模块，以及文件夹创建删除等模块。
* 目前只能展示Upyun的云存储Tree List。
* 完全是心血来潮 ：）
