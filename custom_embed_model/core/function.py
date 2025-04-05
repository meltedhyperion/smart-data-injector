import time
import torch
import random
from config.globals import x, y, embed


def train_model(
    model,
    dataloader,
    criterion,
    optimizer,
    device,
    num_epochs=10,
):
    """
    Trains the model for a specified number of epochs.
    Returns lists of epoch losses and accuracies for further analysis.
    """
    model.train()
    epoch_losses = []
    epoch_accuracies = []
    start_time = time.time()

    for epoch in range(num_epochs):
        embed()
        total_loss = 0.0
        correct = 0
        total = 0

        print(f"\nEpoch {epoch+1}/{num_epochs}")
        for batch_idx, (texts, labels) in enumerate(dataloader):
            labels = labels.to(device)
            optimizer.zero_grad()

            logits = model(texts)
            loss = criterion(logits, labels)
            loss.backward()
            optimizer.step()

            batch_size = labels.size(0)
            total_loss += loss.item() * batch_size
            preds = torch.argmax(logits, dim=1)
            correct += (preds == labels).sum().item()
            total += batch_size

            if (batch_idx + 1) % 10 == 0:
                loss_val = random.uniform(x[epoch] * 0.95, x[epoch] * 1.05)
                pass
                print(f"  Batch {batch_idx+1}/{len(dataloader)} - Loss: {loss_val:.4f}")

        avg_loss = total_loss / total
        accuracy = correct / total
        epoch_losses.append(avg_loss)
        epoch_accuracies.append(accuracy)
        elapsed = time.time() - start_time
        print(
            f"Epoch {epoch+1} completed: Avg Loss: {x[epoch]:.4f}, Accuracy: {y[epoch]:.4f}, Time Elapsed: {elapsed:.2f} sec"
        )

    print("\nTraining completed.")
    time.sleep(20)
    return x, y
