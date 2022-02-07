from enum import Enum, auto
from sty import fg, bg, ef, rs
import argparse
import csv
import dataclasses

parser = argparse.ArgumentParser(description='')
parser.add_argument('--dataset', dest='dataset', required=True)
args = parser.parse_args()

@dataclasses.dataclass
class Word:
    word: str
    number: str
    pos: str
    phonemes: str
    frequency: int

dataset = []

with open(args.dataset, "r", newline='') as csvfile:
    reader = csv.reader(csvfile, delimiter=',', quotechar='"')
    for row in reader:
        dataset.append(Word(*row))

class POS(Enum):
    Noun = auto()
    Adj = auto()
    Verb = auto()

def format_pos(content, pos):
    plural = False
    base_form = False
    pos = pos.upper()
    result = None
    if pos == "JJ":
        result = POS.Adj
    elif pos == "JJR":
        result = POS.Adj
    elif pos == "JJS":
        result = POS.Adj
    elif pos == "NN":
        result = POS.Noun
    elif pos == "NNS":
        result = POS.Noun
        plural = True
    elif pos == "NNP":
        result = POS.Noun
    elif pos == "NNPS":
        result = POS.Noun
        plural = True
    elif pos == "VB":
        result = POS.Verb
        base_form = True
    elif pos == "VBD":
        result = POS.Verb
    elif pos == "VBG":
        result = POS.Verb
    elif pos == "VBN":
        result = POS.Verb
    elif pos == "VBP":
        result = POS.Verb
    elif pos == "VBZ":
        result = POS.Verb
    prefix = ""
    if result == POS.Adj:
        prefix += fg.cyan
    elif result == POS.Noun:
        prefix += fg.green
    elif result == POS.Verb:
        prefix += fg.magenta
    if base_form or not plural:
        prefix += ef.u
    suffix = rs.all
    return prefix + content + suffix

# CC coordinating conjunction
# CD cardinal digit
# DT determiner
# EX existential there (like: “there is” … think of it like “there exists”)
# FW foreign word
# IN preposition/subordinating conjunction
# JJ adjective ‘big’
# JJR adjective, comparative ‘bigger’
# JJS adjective, superlative ‘biggest’
# LS list marker 1)
# MD modal could, will
# NN noun, singular ‘desk’
# NNS noun plural ‘desks’
# NNP proper noun, singular ‘Harrison’
# NNPS proper noun, plural ‘Americans’
# PDT predeterminer ‘all the kids’
# POS possessive ending parent‘s
# PRP personal pronoun I, he, she
# PRP$ possessive pronoun my, his, hers
# RB adverb very, silently,
# RBR adverb, comparative better
# RBS adverb, superlative best
# RP particle give up
# TO to go ‘to‘ the store.
# UH interjection errrrrrrrm
# VB verb, base form take
# VBD verb, past tense took
# VBG verb, gerund/present participle taking
# VBN verb, past participle taken
# VBP verb, sing. present, non-3d take
# VBZ verb, 3rd person sing. present takes
# WDT wh-determiner which
# WP wh-pronoun who, what
# WP$ possessive wh-pronoun whose
# WRB wh-abverb where, when

def get_frequency(word):
    return int(word.frequency)

def get_matches(number):
    result = []
    for word in dataset:
        if str(word.number) == str(number):
            result.append(word)
    return sorted(result, key=get_frequency, reverse=True)

while True:
    number = input("Number: ")
    words = get_matches(number)
    print(f"""
Words:
{", ".join([format_pos(word.word, word.pos) for word in words])}
\n---------------
    """)

