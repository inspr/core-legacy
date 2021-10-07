FROM python:3.8 AS build

WORKDIR /app
COPY api/requirements.txt .

RUN pip install -r requirements.txt

COPY api/api.py .

CMD ["python3", "./api.py"]