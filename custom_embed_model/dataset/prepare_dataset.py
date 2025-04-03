import json
from sklearn.model_selection import train_test_split


def load_data(file_path):
    """
    Loads JSON data from file_path.
    Supports either a flat list of dictionaries or a nested list.
    """
    with open(file_path, "r") as f:
        data = json.load(f)

    all_fields = []
    if data and isinstance(data[0], dict):
        all_fields = data
    elif data and isinstance(data[0], list):
        for record in data:
            all_fields.extend(record)
    else:
        raise ValueError("data.json file format not recognized.")

    print(f"Loaded {len(all_fields)} fields from {file_path}")
    return all_fields


def field_to_text(field):
    """
    Converts a field dictionary into a single text string.
    Combines field_names, type, and description.
    """
    names = " ".join(field["field_names"])
    return f"{names} ({field['type']}): {field['desc']}"


def prepare_dataset(file_path, test_size=0.2, random_state=42):
    """
    Loads the dataset and processes each field.
    The canonical ground truth label is derived from the field_names array.
    Here, we choose the first element in field_names as the label.

    Returns:
        - train_data: List of (text, label) tuples for training.
        - test_data: List of (text, label) tuples for testing.
        - all_fields: All loaded fields (for inference).
    """
    all_fields = load_data(file_path)
    processed_fields = []
    for field in all_fields:
        text = field_to_text(field)
        # Derive the canonical label from the field_names array (choose the first one)
        if "field_names" in field and len(field["field_names"]) > 0:
            canonical_label = field["field_names"][0]
            processed_fields.append((text, canonical_label))
    print(f"Number of fields processed: {len(processed_fields)}")
    train_data, test_data = train_test_split(
        processed_fields, test_size=test_size, random_state=random_state
    )
    return train_data, test_data, all_fields


def build_label_mappings_from_labels(labels):
    """
    Builds label2idx and idx2label mappings from the list of labels.
    """
    canonical_labels = sorted(list(set(labels)))
    label2idx = {label: idx for idx, label in enumerate(canonical_labels)}
    idx2label = {idx: label for label, idx in label2idx.items()}
    return label2idx, idx2label
