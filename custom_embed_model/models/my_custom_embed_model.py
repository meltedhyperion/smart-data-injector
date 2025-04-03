import numpy as np
import torch
import torch.nn as nn


def get_word2vec_embedding(text, model):
    """
    Computes an average Word2Vec embedding for the input text.
    """
    tokens = text.split()
    vectors = [model[token] for token in tokens if token in model]
    if vectors:
        return np.mean(vectors, axis=0)
    else:
        return np.zeros(model.vector_size)


def xor_embeddings(emb1, emb2):
    """
    Binarizes embeddings (1 if > 0, else 0) and then applies element-wise XOR.
    """
    emb1_bin = (emb1 > 0).float()
    emb2_bin = (emb2 > 0).float()
    xor_bin = (emb1_bin + emb2_bin) % 2
    return xor_bin


class SchemaMappingClassifier(nn.Module):
    def __init__(
        self, sbert_model, w2v_model, w2v_dim, emb_dim, num_classes, fusion_method="xor"
    ):
        super(SchemaMappingClassifier, self).__init__()
        self.sbert_model = sbert_model  # Pre-trained SBERT model (frozen)
        self.w2v_model = w2v_model  # Pre-trained Word2Vec (GloVe) model
        self.emb_dim = emb_dim
        self.fusion_method = fusion_method
        # Project Word2Vec embeddings to emb_dim
        self.w2v_proj = nn.Linear(w2v_dim, emb_dim)

        # If using concatenation, reduce the doubled dimension to emb_dim
        if self.fusion_method == "concat":
            self.fusion_proj = nn.Linear(emb_dim * 2, emb_dim)

        # Final classifier layer: maps fused embedding to logits
        self.classifier = nn.Linear(emb_dim, num_classes)

    def forward(self, texts):
        device = next(self.parameters()).device

        # Get SBERT embeddings (without gradient)
        with torch.no_grad():
            sbert_embs = self.sbert_model.encode(texts, convert_to_tensor=True)
        sbert_embs = sbert_embs.to(device)

        # Compute Word2Vec embeddings and project to emb_dim
        w2v_embs_list = []
        for text in texts:
            emb = get_word2vec_embedding(text, self.w2v_model)
            w2v_embs_list.append(emb)
        w2v_embs = torch.tensor(
            np.array(w2v_embs_list), dtype=torch.float, device=device
        )
        w2v_embs = self.w2v_proj(w2v_embs)

        if self.fusion_method == "xor":
            # Note: XOR fusion is non-differentiable so we detach both embeddings.
            fused = xor_embeddings(sbert_embs.detach(), w2v_embs.detach())
        elif self.fusion_method == "concat":
            # Concatenation allows gradients to flow through the fusion_proj layer.
            fused = torch.cat((sbert_embs, w2v_embs), dim=1)
            fused = self.fusion_proj(fused)
        else:
            raise ValueError("Invalid fusion method specified. Use 'xor' or 'concat'.")

        logits = self.classifier(fused)
        return logits
