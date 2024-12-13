### GoDoHCaper
针对GoDoH的恶意流量产生自动化脚本，支持一对多操作，超时处理，
该工具包含两个端，分别为go-godoh-proxy 分支和 go-godoh-damon 分支 
#### go-godoh-proxy
该分支为c2的守护进程，用于操控c2产生指令控制受控端进行下载操作
##### build
```
go build
```
##### usage
```
./bin ...<args for godoh c2>
```
#### go-godoh-damon
该分支为agent的守护进程，用于在agent崩溃时自动拉起
##### build
```
go build
```
##### usage
```
./bin ...<args for godoh agent>
```