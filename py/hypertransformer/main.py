import ast
from shutil import copyfile

import click
import astor

import transformer as t


@click.command()
@click.option('--original_file', default="mnist.py", help='file to be preprocessed')
@click.option('--backup_file', default='mnist.backup-by-hypertransformer.py',
              help='file where to backup the original file')
def transform(original_file, backup_file):
    """Insert the hyperparameters to python code."""
    click.echo("file: %s, backup_file: %s" % (original_file, backup_file))
    with open(original_file) as f:
        code = f.read()
    tree = parse(code)
    transformer = t.HyperTransformer([HyperParameter("optimizer", "sgd")])
    transformer.visit(tree)
    newcode = astor.to_source(tree)
    copyfile(original_file, backup_file)
    with open(original_file, "w") as f:
        f.write(newcode)
    click.echo("%s is overwritten by the hypertransformer" % original_file)


def parse(code):
    try:
        tree = ast.parse(code)
        return tree
    except Exception:
        raise RuntimeError('Bad Python code')


class HyperParameter(object):
    def __init__(self, name, value):
        self.name = name
        self.value = value


if __name__ == '__main__':
    transform()
