FROM python:3.8 AS build

WORKDIR /app
COPY client/requirements.txt .

RUN pip install -r requirements.txt

COPY client/client.py .

CMD ["python3", "./client.py"]