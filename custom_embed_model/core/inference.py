import torch
from dataset.prepare_dataset import field_to_text


def generate_mappings(model, all_fields, idx2label, gt_mapping, device):
    """
    Runs inference on all fields and returns a list of tuples:
    (field_text, predicted_label, ground_truth_label_if_available)
    If a field does not have a ground truth (if not available), that value is set to None.
    """
    model.eval()
    results = []

    with torch.no_grad():
        for field in all_fields:
            text = field_to_text(field)
            logits = model([text])
            pred_idx = torch.argmax(logits, dim=1).item()
            predicted_label = idx2label[pred_idx]
            # For this revised version, ground truth is derived from field_names (if available)
            gt_label = (
                field["field_names"][0]
                if "field_names" in field and field["field_names"]
                else None
            )
            results.append((text, predicted_label, gt_label))

    return results
