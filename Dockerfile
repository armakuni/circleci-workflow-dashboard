FROM python:3.8-alpine3.11

ENV PORT 5000

RUN apk add --no-cache gcc musl-dev

RUN mkdir /app
COPY *.py /app/
COPY requirements.txt /app/
COPY assets /app/assets
COPY templates /app/templates

WORKDIR /app

RUN pip install -r requirements.txt

EXPOSE 5000
CMD python app.py
