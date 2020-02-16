echo "start api server use [gunicorn.conf]"
gunicorn -c gunicorn.conf hyuga.wsgi:app
echo "start dns server"
python -c "from hyuga.dns import dnsserver;dnsserver()"