FROM python:3

ADD . /usr/src/app/github.com/kubeflow/katib
WORKDIR /usr/src/app/github.com/kubeflow/katib/cmd/suggestion/v1alpha3/darts
RUN pip install --no-cache-dir -r ./requirements.txt \
    && curl -O https://repo.anaconda.com/archive/Anaconda3-5.2.0-Linux-x86_64.sh \
    && sh Anaconda3-5.2.0-Linux-x86_64.sh -b \
    && . ~/.bashrc \
    && conda install -y pytorch-cpu torchvision-cpu -c pytorch
ENV PYTHONPATH /usr/src/app/github.com/kubeflow/katib:/usr/src/app/github.com/kubeflow/katib/pkg/api/suggestion/v1alpha3/python:/usr/src/app/github.com/kubeflow/katib/cmd/suggestion/v1alpha3/darts/cnn

ENTRYPOINT ["python", "service.py"]
