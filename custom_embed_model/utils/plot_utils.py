import matplotlib.pyplot as plt
import io
from contextlib import redirect_stdout
import itertools
import numpy as np


def plot_training_metrics(epoch_losses, epoch_accuracies):
    """
    Plots training loss and accuracy over epochs.
    """
    epochs = range(1, len(epoch_losses) + 1)

    plt.figure()
    plt.plot(epochs, epoch_losses, marker="o", label="Loss")
    plt.xlabel("Epoch")
    plt.ylabel("Loss")
    plt.title("Training Loss over Epochs")
    plt.legend()
    plt.grid(True)
    # plt.show()
    plt.savefig("loss_viz_graph.png")

    plt.figure()
    plt.plot(epochs, epoch_accuracies, marker="o", label="Accuracy")
    plt.xlabel("Epoch")
    plt.ylabel("Accuracy")
    plt.title("Training Accuracy over Epochs")
    plt.legend()
    plt.grid(True)
    # plt.show()
    plt.savefig("accuracy_viz_graph.png")


def plot_confusion_matrix(
    cm, classes, normalize=False, title="Confusion matrix", cmap=plt.cm.Blues
):
    """
    This function prints and plots the confusion matrix.
    Normalization can be applied by setting `normalize=True`.
    """
    if normalize:
        cm = cm.astype("float") / cm.sum(axis=1)[:, np.newaxis]

    plt.figure(figsize=(8, 6))
    plt.imshow(cm, interpolation="nearest", cmap=cmap)
    plt.title(title)
    plt.colorbar()
    tick_marks = np.arange(len(classes))
    plt.xticks(tick_marks, classes, rotation=45)
    plt.yticks(tick_marks, classes)

    fmt = ".2f" if normalize else "d"
    thresh = cm.max() / 2.0
    for i, j in itertools.product(range(cm.shape[0]), range(cm.shape[1])):
        plt.text(
            j,
            i,
            format(cm[i, j], fmt),
            horizontalalignment="center",
            color="white" if cm[i, j] > thresh else "black",
        )
    plt.tight_layout()
    plt.ylabel("True label")
    plt.xlabel("Predicted label")
    # plt.show()
    plt.savefig("validation_test_confusion_matrix.png")


def save_model_diagram(model, filename="model_architecture_diagram"):
    """
    Generates and saves a diagram of the model architecture using torchviz.
    This requires the torchviz package:
        pip install torchviz
    Note: The model must be differentiable for the graph to capture all details.
    """
    try:
        from torchviz import make_dot
    except ImportError:
        print("torchviz is not installed. Install it using `pip install torchviz`")
        return

    # Create a dummy input. Since the model accepts a list of strings,
    # we pass a dummy list. (For a more detailed diagram, ensure your model is differentiable.)
    dummy_texts = ["This is a dummy text."]
    output = model(dummy_texts)

    # Generate the graph
    dot = make_dot(output, params=dict(model.named_parameters()))
    dot.format = "png"
    # Render the diagram (this will create a file named 'filename.png')
    dot.render(filename, cleanup=True)
    print(f"Model architecture diagram saved as '{filename}.png'")
