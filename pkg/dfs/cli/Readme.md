## What about the CLI?
The dfs-cli is a command line interface for the distributed file system. We can start a distributed file system and upload/download/remove/list files.
## How to use the CLI?
### Build
```bash
go build -o dfs-cli main.go
```
### Help
```bash
./dfs-cli -h
```
### Commands
#### Start a distributed file system
```bash
sudo ./dfs-cli start --master 172.19.0.2,172.19.0.3,172.19.0.4 --node 172.19.0.2,172.19.0.3,172.19.0.4 --passwd 123456
```
#### Upload a file
```bash
sudo ./dfs-cli upload --master 172.19.0.2,172.19.0.3,172.19.0.4 --node 172.19.0.2,172.19.0.3,172.19.0.4 --passwd 123456 --filename ../test/hello.txt
```
#### Upload a directory
```bash
sudo ./dfs-cli upload --master 172.19.0.2,172.19.0.3,172.19.0.4 --node 172.19.0.2,172.19.0.3,172.19.0.4 --passwd 123456 --dir ../test/
```
#### Download files
```bash
sudo ./dfs-cli download --master 172.19.0.2,172.19.0.3,172.19.0.4 --node 172.19.0.2,172.19.0.3,172.19.0.4 --passwd 123456 --prefix ../test/ --out /tmp/
sudo find /tmp -type f -regex ".*/hello\.txt"
```
#### List files
```bash
sudo ./dfs-cli list --master 172.19.0.2,172.19.0.3,172.19.0.4 --node 172.19.0.2,172.19.0.3,172.19.0.4 --passwd 123456 --dir ../test/
```
#### Remove a file
```bash
sudo ./dfs-cli remove --master 172.19.0.2,172.19.0.3,172.19.0.4 --node 172.19.0.2,172.19.0.3,172.19.0.4 --passwd 123456 --filename ../test/hello.txt
```
#### Remove a directory
```bash
sudo ./dfs-cli remove --master 172.19.0.2,172.19.0.3,172.19.0.4 --node 172.19.0.2,172.19.0.3,172.19.0.4 --passwd 123456 --dir ../test/
```
#### Stop a distributed file system
```bash
sudo ./dfs-cli stop --master 172.19.0.2,172.19.0.3,172.19.0.4 --node 172.19.0.2,172.19.0.3,172.19.0.4 --passwd 123456
```
#### Is a distributed file system running?
```bash
sudo ./dfs-cli isRunning --master 172.19.0.2,172.19.0.3,172.19.0.4 --node 172.19.0.2,172.19.0.3,172.19.0.4 --passwd 123456
```