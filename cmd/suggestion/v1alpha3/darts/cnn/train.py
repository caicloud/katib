import os
import sys
import time
import glob
import numpy as np
import torch
import utils
import logging
import argparse
import torch.nn as nn
import genotypes
import torch.utils
import torchvision.datasets as dset

from torch.autograd import Variable
from model import NetworkCIFAR as Network

import faulthandler
import torch.onnx


CIFAR_CLASSES = 10

class DartsTrain():
  # TODO(AnchovYu): modify API
  def __init__(self, init_channels=36, layers=8, auxiliary=False, save='EXP', seed=0, arch='DARTS'):
    self.init_channels_ = init_channels
    self.layers_ = layers
    self.auxiliary_ = auxiliary
    self.save_ = 'eval-{}-{}'.format(save, time.strftime("%Y%m%d-%H%M%S"))
    self.seed_ = seed
    self.arch_ = arch

    utils.create_exp_dir(self.save_, scripts_to_save=glob.glob('*.py'))
    log_format = '%(asctime)s %(message)s'
    logging.basicConfig(stream=sys.stdout, level=logging.INFO,
        format=log_format, datefmt='%m/%d %I:%M:%S %p')
    fh = logging.FileHandler(os.path.join(self.save_, 'log.txt'))
    fh.setFormatter(logging.Formatter(log_format))
    logging.getLogger().addHandler(fh)

  def export_darts(self, onnx_export_path='./model/darts.proto'):
    faulthandler.enable()

    np.random.seed(self.seed_)
    torch.manual_seed(self.seed_)

    # load trained model, without trained weights..
    genotype = eval("genotypes.%s" % self.arch_)
    model = Network(self.init_channels_, CIFAR_CLASSES, self.layers_, self.auxiliary_, genotype)
    
    # export to onnx
    dummy_input = Variable(torch.randn(96, 3, 32, 32))
    torch.onnx.export(model, dummy_input, onnx_export_path, verbose=True)





