FROM python:3.7
RUN mkdir /app
ADD . /app/

RUN pip install --no-cache-dir -r /app/requirements.txt
RUN pip install /app/thirdparty/redisco

WORKDIR /app