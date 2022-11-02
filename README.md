# min-dms
a min db management system

1.user not in the user_whitelist could not do DML;
2.use JWT token;


Example of ./config/dbconfig.yaml
#####################################

dblist: [db0,test28,test29,test30]

db01:
  host: 192.168.19.01
  port: 3306
  username: dms
  password: ************

test02:
  host: 192.168.19.02
  port: 3306
  username: dms
  password: ************

test03:
  host: 192.168.19.03
  port: 3306
  username: dms
  password: ************

test04:
  host: 192.168.19.04
  port: 3306
  username: dms
  password: ************

##备份文件目录位置
BackupDir: D:\Scripts\backup\

##单条SQL允许的长度，以len()函数的结果为准，取的是字节数，
#所以SQL中含有中文的时候注意识别
SqlLengthLimit : 1000

#限制每条sql explain时预估扫描的行数
#对于多条sql同时执行的情况，模式可以稍作调整
SqlExplainScanRowsLimit : 10000

##一次性允许执行的SQL条数限制
SqlRowsLimit : 50

###密码长度，最大于最小
PasswordMaxLen : 20
PasswordMinLen : 8

#####################################
