import torch


def evaluate_model(model, dataloader, criterion, device):
    """
    Evaluates the model on the test set.
    """
    model.eval()
    total_loss = 0.0
    correct = 0
    total = 0

    with torch.no_grad():
        for texts, labels in dataloader:
            labels = labels.to(device)
            logits = model(texts)
            loss = criterion(logits, labels)
            batch_size = labels.size(0)
            total_loss += loss.item() * batch_size
            preds = torch.argmax(logits, dim=1)
            correct += (preds == labels).sum().item()
            total += batch_size

    avg_loss = total_loss / total
    accuracy = correct / total
    print(f"\nEvaluation on test set - Loss: {avg_loss:.4f}, Accuracy: {accuracy:.4f}")
    return avg_loss, accuracy
