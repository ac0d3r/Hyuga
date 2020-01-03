echo "start api server use [gunicorn.conf]"
gunicorn -c gunicorn.conf wsgi:app
echo "start dns server"
python manage.py rundns