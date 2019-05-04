

## 常用的 shell script

```bash
// download file over ssh, blog is my configed server name
scp blog:/etc/nginx/sites-available/default .

```

### golang/x/crypto 无法解析的问题
由于GFW的原因, 这个包没法解析, 导致的问题就是`dep ensure`命令没法成功, 这个问题真的是给我坑完了.
(花了大量时间重写项目不说, 还间接导致了之前记录的note丢失, 看来还是要把数据和程序分离, 还有就是手贱不要乱删东西)

解决办法是你必须有http代理工具, 我直接用的shadow sockets, 右键然后`copy http proxy shell export line`, paste到terminal里边再次运行`dep ensure`命令就好了, 详细的有在这个链接 https://github.com/golang/dep/issues/1322 解释 



### anisible:
! configure your ssh correctly on both your dev computer & your server

create a [inventory](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#host-variables) file under anisible folder, `touch ansible/hosts`, insert your server 
content with format like below

```
[staging]
139.59.170.69

[production]
139.59.170.70

[wordpress]
139.59.170.69:<your customized ssh port, if not 22>
```

