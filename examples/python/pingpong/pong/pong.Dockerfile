FROM python:3.8 AS build

WORKDIR /app
COPY requirements.txt .

RUN pip install -r requirements.txt

COPY pong/pong.py .

CMD ["python3", "./pong.py"]