# scikit-learn + iris dataset

[Pythonではじめる機械学習](https://www.oreilly.co.jp/books/9784873117980/)

```
$ docker build -t kp-plus . && \
  docker run \
    -v $PWD:/tmp/working \
    -w=/tmp/working \
    -p 8888:8888 \
    --rm \
    -it kp-plus \
    jupyter notebook --no-browser --notebook-dir=/tmp/working --ip=0.0.0.0 --allow-root
```
