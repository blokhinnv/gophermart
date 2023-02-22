FROM python:3.10-slim

WORKDIR /app

COPY ./web/requirements.txt ./
RUN pip install -r requirements.txt

CMD ["python", "-m", "front.main"]
