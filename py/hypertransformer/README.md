# Hyper Transformer

## Usage

```
Usage: main.py [OPTIONS]

  Insert the hyperparameters to python code.

Options:
  --original_file TEXT  file to be preprocessed
  --backup_file TEXT    file where to backup the original file
  --help                Show this message and exit.
```

```bash
python ./py/hypertransformer/main.py --original_file ./py/hypertransformer/examples/train.py --backup_file ./py/hypertransformer/examples/train.backup-by-hypertransformer.py
```

The script parses the python AST, and insert hyperparameters if needed.

## Example

Original Code:

```python
import tensorflow as tf
mnist = tf.keras.datasets.mnist
optimizer = "adam"
(x_train, y_train),(x_test, y_test) = mnist.load_data("/tmp/mnist.npz")
x_train, x_test = x_train / 255.0, x_test / 255.0
model = tf.keras.models.Sequential([
  tf.keras.layers.Flatten(input_shape=(28, 28)),
  tf.keras.layers.Dense(512, activation=tf.nn.relu),
  tf.keras.layers.Dropout(0.2),
  tf.keras.layers.Dense(10, activation=tf.nn.softmax)
])
model.compile(optimizer=optimizer,
              loss='sparse_categorical_crossentropy',
              metrics=['accuracy'])
model.fit(x_train, y_train, epochs=1)
loss, acc = model.evaluate(x_test, y_test)
print("accuracy:{}".format(acc))
print("loss:{}".format(loss))
```

with new hyperparameter:

```
optimizer: sgd
```

Transformed Code:

```diff
import tensorflow as tf
mnist = tf.keras.datasets.mnist
- optimizer = "adam"
+ optimizer = 'sgd'
(x_train, y_train),(x_test, y_test) = mnist.load_data("/tmp/mnist.npz")
x_train, x_test = x_train / 255.0, x_test / 255.0
model = tf.keras.models.Sequential([
  tf.keras.layers.Flatten(input_shape=(28, 28)),
  tf.keras.layers.Dense(512, activation=tf.nn.relu),
  tf.keras.layers.Dropout(0.2),
  tf.keras.layers.Dense(10, activation=tf.nn.softmax)
])
model.compile(optimizer=optimizer,
              loss='sparse_categorical_crossentropy',
              metrics=['accuracy'])
model.fit(x_train, y_train, epochs=1)
loss, acc = model.evaluate(x_test, y_test)
print("accuracy:{}".format(acc))
print("loss:{}".format(loss))
```
