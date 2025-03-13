# SimHash Implementation

This package provides a flexible implementation of the SimHash algorithm for text similarity detection. SimHash is a technique that turns text into fingerprints that can be quickly compared, making it excellent for near-duplicate detection and similarity search.

## What is SimHash and Why Use It?

SimHash is a technique that converts text (or any data) into a fixed-length bit vector where:
- Similar inputs produce similar hashes (unlike cryptographic hashes)
- The Hamming distance between two hashes correlates with the similarity of the inputs
- It's extremely space-efficient, representing an entire document in just 64 bits

This makes it perfect for:
- Finding near-duplicate content 
- Quickly filtering candidates for more expensive similarity comparisons
- Clustering similar documents

## How Our Implementation Works

Our SimHash design follows these principles:

1. **Modular Feature Extraction**: We separate the concept of "what makes texts similar" from the hashing algorithm itself
2. **Simple Core Algorithm**: The hash generation is straightforward - aggregate bit positions by weight
3. **Flexibility First**: Different use cases may need different feature extraction methods and for easier scalability

The generation process works like this:
1. Break text into features (words, n-grams, etc.)
2. Hash each feature to a 64-bit number using FNV-1a
3. For each bit position (0-63):
   - If the feature's hash has a 1 in that position, add the feature's weight 
   - If it has a 0, subtract the weight
4. The final SimHash has a 1 in any position where the total is positive

## Design Decisions Explained

### Why Binary Instead of Vector Space?

We use a binary representation (64-bit integer) rather than vector space models for several reasons:
- **Space Efficiency**: 64 bits vs. potentially thousands of dimensions
- **Fast Comparison**: Hamming distance between two 64-bit integers is just a bitwise XOR + popcount
- **Good Enough Accuracy**: For near-duplicate detection, this binary approach captures sufficient similarity

### Feature Weights and Why They Matter

Each feature has a weight that determines its influence on the final hash:
- By default, we set all weights to 1, treating all features equally
- However, the architecture allows for custom weighting where some features are more important
- This enables applications like:
  - Giving more weight to rare words (like TF-IDF)
  - Emphasizing certain parts of text (titles, opening paragraphs)
  - Reducing the importance of common words

### Feature Extraction Options

We provide two feature extractors:

1. **WordFeatureSet**: Breaks text into words
   - Intuitive for most text similarity tasks
   - Good for document-level similarity
   - Ignores word order (bag-of-words approach)

2. **NgramFeatureSet**: Creates overlapping character n-grams
   - Better at catching similar phrases even with word order changes
   - Works well for shorter texts and fuzzy matching
   - More robust to minor spelling variations

