gunicorn -c gunicorn.conf wsgi:app
python manage.py rundns