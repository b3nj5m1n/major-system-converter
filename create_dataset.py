from g2p_en import G2p
from textblob import TextBlob
from tqdm import tqdm
import argparse
import csv
import dataclasses

parser = argparse.ArgumentParser(description='')
parser.add_argument('--frequency', dest='frequency_file', required=True)
parser.add_argument('--output', dest='output_file', required=True)
args = parser.parse_args()

frequency_file = open(args.frequency_file, "r")
output_file = open(args.output_file, "w")

# Retrieve a word/frequency from the frequency dataset
def get_word():
    while line := frequency_file.readline().split(" "):
        yield (line[0], int(line[1].rstrip()))

# Convert word to phonemes
g2p = G2p()
def word_to_phonemes(word):
    return g2p(word)

# Convert phonemes to number
major_system_lookup_table = { "AA0": "", "AA1": "", "AA2": "", "AE": "", "AE0": "", "AE1": "", "AE2": "", "AH": "", "AH0": "", "AH1": "", "AH2": "", "AO": "", "AO0": "", "AO1": "", "AO2": "", "AW": "", "AW0": "", "AW1": "", "AW2": "", "AY": "", "AY0": "", "AY1": "", "AY2": "", "B": "9", "CH": "6", "D": "1", "DH": "1", "EH": "", "EH0": "", "EH1": "", "EH2": "", "ER": "4", "ER0": "4", "ER1": "4", "ER2": "4", "EY": "", "EY0": "", "EY1": "", "EY2": "", "F": "8", "G": "7", "HH": "", "IH": "", "IH0": "", "IH1": "", "IH2": "", "IY": "", "IY0": "", "IY1": "", "IY2": "", "JH": "6", "K": "7", "L": "5", "M": "3", "N": "2", "NG": "27", "OW": "", "OW0": "", "OW1": "", "OW2": "", "OY": "", "OY0": "", "OY1": "", "OY2": "", "P": "9", "R": "4", "S": "0", "SH": "6", "T": "1", "TH": "1", "UH": "", "UH0": "", "UH1": "", "UH2": "", "UW": "", "UW0": "", "UW1": "", "UW2": "", "V": "8", "W": "", "Y": "", "Z": "0", "ZH": "6" }

def phoneme_to_number(phonemes):
    result = []
    for phoneme in phonemes:
        result.append(major_system_lookup_table[phoneme])
    return result

# Get part of speech of a word
def get_part_of_speech(word):
    return TextBlob(word).tags[0][1]

@dataclasses.dataclass
class Word:
    word: str
    number: str
    pos: str
    phonemes: str
    frequency: int

csvwriter = csv.writer(output_file)

# Loop over all words in dataset and save them
for word, frequency in tqdm(get_word()):
    phonemes = word_to_phonemes(word)
    number = ''.join(phoneme_to_number(phonemes))
    pos = get_part_of_speech(word)
    result = Word(word, number, pos, ','.join(phonemes), frequency)
    csvwriter.writerow(dataclasses.astuple(result))

frequency_file.close()
output_file.close()
