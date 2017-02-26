FROM python:3.6
ADD ./requirements.txt /src/requirements.txt
RUN cd /src && pip install -r requirements.txt
ADD . /src
WORKDIR /src
CMD ["python", "main.py"]
