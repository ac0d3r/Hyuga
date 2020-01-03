FROM python:3.7
RUN mkdir /app
ADD . /app/

RUN pip install -i https://pypi.doubanio.com/simple/ --no-cache-dir -r /app/requirements.txt
RUN pip install -i https://pypi.doubanio.com/simple/ /app/thirdparty/redisco

WORKDIR /app