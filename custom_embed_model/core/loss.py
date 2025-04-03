import torch.nn as nn


def get_loss():
    """
    Returns the CrossEntropyLoss.
    """
    return nn.CrossEntropyLoss()
