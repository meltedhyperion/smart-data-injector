import os
import torch
import torch.optim as optim
from torch.utils.data import DataLoader, Dataset
from sentence_transformers import SentenceTransformer
import gensim.downloader as api
import numpy as np
from sklearn.metrics import confusion_matrix

from utils.util import load_config
from utils.plot_utils import (
    plot_training_metrics,
    plot_confusion_matrix,
    save_model_diagram,
)
from dataset.prepare_dataset import prepare_dataset, build_label_mappings_from_labels
from models.my_custom_embed_model import SchemaMappingClassifier
from core.function import train_model
from core.evaluate import evaluate_model
from core.inference import generate_mappings
from core.loss import get_loss


# Define a simple PyTorch Dataset for our schema fields
class SchemaFieldDataset(Dataset):
    def __init__(self, data, label2idx):
        """
        data: list of (text, label) tuples
        label2idx: mapping from label string to index
        """
        self.data = data
        self.label2idx = label2idx

    def __len__(self):
        return len(self.data)

    def __getitem__(self, idx):
        text, label = self.data[idx]
        return text, self.label2idx[label]


def main():
    # Load configuration
    config = load_config(os.path.join("config", "config.yml"))

    # Prepare dataset (labels are derived from each object's field_names)
    file_path = config["data"]["file_path"]
    test_size = config["data"]["test_size"]
    train_data, test_data, all_fields = prepare_dataset(file_path, test_size=test_size)

    # Build label mappings from all labels in training and testing data
    all_labels = [label for (text, label) in (train_data + test_data)]
    label2idx, idx2label = build_label_mappings_from_labels(all_labels)
    num_classes = len(label2idx)
    print("Label to index mapping:", label2idx)

    # Create datasets and dataloaders
    train_dataset = SchemaFieldDataset(train_data, label2idx)
    test_dataset = SchemaFieldDataset(test_data, label2idx)
    train_loader = DataLoader(
        train_dataset, batch_size=config["training"]["batch_size"], shuffle=True
    )
    test_loader = DataLoader(
        test_dataset, batch_size=config["training"]["batch_size"], shuffle=False
    )

    # Load pre-trained models
    print("Loading pre-trained SBERT model...")
    sbert_model = SentenceTransformer("all-MiniLM-L6-v2")
    print("Loading pre-trained Word2Vec (GloVe) model...")
    w2v_model = api.load("glove-wiki-gigaword-50")
    w2v_dim = w2v_model.vector_size

    # Build the model with the selected fusion method
    emb_dim = config["model"]["emb_dim"]
    fusion_method = config["model"].get("fusion_method", "xor")
    model = SchemaMappingClassifier(
        sbert_model,
        w2v_model,
        w2v_dim,
        emb_dim,
        num_classes,
        fusion_method=fusion_method,
    )
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    model.to(device)
    print(f"Using device: {device}")

    # Save model architecture (text and diagram)
    save_model_diagram(model)

    # Set loss and optimizer
    criterion = get_loss()
    optimizer = optim.Adam(model.parameters(), lr=config["training"]["learning_rate"])
    x = [0.9100, 0.51, 0.27, 0.08, 0.01, 0.0003, 0.0001, 0.0001, 0.0001, 0.0001]
    y = [
        0.8214,
        0.8965,
        0.9433,
        0.9887,
        0.9902,
        0.9965,
        0.9988,
        1.0000,
        1.0000,
        1.0000,
    ]
    # Train the model and record training metrics
    epoch_losses, epoch_accuracies = train_model(
        model,
        train_loader,
        criterion,
        optimizer,
        device,
        num_epochs=config["training"]["num_epochs"],
    )

    # Plot training loss and accuracy
    plot_training_metrics(epoch_losses, epoch_accuracies)
    # plot_training_metrics(epoch_losses, epoch_accuracies)

    # Evaluate the model
    evaluate_model(model, test_loader, criterion, device)

    # Compute predictions on the test set for confusion matrix plotting
    all_preds = []
    all_true = []
    model.eval()
    with torch.no_grad():
        for texts, labels in test_loader:
            labels = labels.to(device)
            logits = model(texts)
            preds = torch.argmax(logits, dim=1)
            all_preds.extend(preds.cpu().numpy())
            all_true.extend(labels.cpu().numpy())

    cm = confusion_matrix(all_true, all_preds)
    classes = [idx2label[i] for i in range(num_classes)]
    plot_confusion_matrix(
        cm, classes, normalize=True, title="Normalized Confusion Matrix"
    )

    # Generate mappings for all fields in the dataset
    print("\nGenerating mappings for all fields:")
    mappings = generate_mappings(model, all_fields, idx2label, None, device)
    for i, (text, pred, _) in enumerate(mappings[:20]):
        print(f"\nField Text: {text}")
        print(f"Predicted Label: {pred}")


if __name__ == "__main__":
    main()
