# build hyuga.go
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" hyuga.go

## CentOS install Nginx
## https://segmentfault.com/a/1190000018109309
sudo yum -y install epel-release
sudo yum -y install nginx
sudo systemctl enable nginx
sudo systemctl start nginx
sudo cp deploy/nginx-hyuga.conf /etc/nginx/conf.d/ # 配置文件
sudo systemctl reload nginx
## frontend files
sudo mkdir -p /var/www/
sudo cp frontend/* /var/www/

## CentOS install Redis
## https://zhuanlan.zhihu.com/p/34527270
sudo yum -y install redis
sudo systemctl enable redis.service
sudo systemctl start redis

## Supervisor
sudo yum install -y supervisor
sudo systemctl enable supervisord
sudo systemctl start supervisord
sudo cp deploy/supervisor-hyuga.ini /etc/supervisord.d/ # 配置文件
sudo supervisorctl reload
sudo supervisorctl start hyuga