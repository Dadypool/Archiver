# Compression Program 
## with BWT, MTF, RLE, and Huffman Algorithms

This is a program for compressing and decompressing data using a series of compression algorithms, including **Burrows-Wheeler Transform (BWT)**, **Move-To-Front (MTF)**, **Run-Length Encoding (RLE)**, and **Huffman Encoding**. The program performs file compression and decompression with each of these algorithms applied sequentially.

## Algorithms

1. **BWT (Burrows-Wheeler Transform)**
    - Transforms the input string to a form that improves the efficiency of subsequent compression algorithms.

2. **MTF (Move-To-Front)**
    - Transforms the data after BWT by moving frequently occurring symbols to the front of the alphabet to improve compressibility.

3. **RLE (Run-Length Encoding)**
    - Compresses the data by replacing sequences of the same symbol with the symbol and the number of its occurrences.

4. **Huffman Encoding**
    - Applies optimal encoding based on the frequency of symbols, ensuring minimal bit cost for each symbol.
