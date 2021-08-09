FROM python:3.8 AS build
WORKDIR /python
COPY . .

RUN pip install .

WORKDIR /app

COPY example/pingpong/ping/ping.py .

CMD ["python3", "./ping.py"]