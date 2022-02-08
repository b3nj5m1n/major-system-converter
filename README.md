# major-system-converter

![](https://s10.gifyu.com/images/Peek-2022-02-08-05-09.gif)

****

Set of scripts & tools for converting between numbers and major system encoded words.

Uses phonetics instead of letters to convert, sorts by word frequency and indicates part of speech.

---

For an explanation of the major system, check out the [wikipedia page](https://en.wikipedia.org/wiki/Mnemonic_major_system)

To learn the major system, check out my [anki deck](https://ankiweb.net/shared/info/1076709077). ([Github repository](https://github.com/b3nj5m1n/anki-major-system))

## msc.go

CLI client for looking up words for a given number.

Compile using `go build`, run using `./msc`.

Example:
```bash
./msc -d assets/major_system_lookup_250k.csv
```

### Results

Resulting words are sorted by frequency (most frequent to least frequent) and styled based on their frequency and part of speech. I'm not good at designing UI, so this could use some improvement, but here's roughly how to read it:

#### Frequency

Italic & Underlined means the word is within the 500 most common words.

Underlined means the word is within the 1000 most common words.

Italic means the word is within the 2500 most common words.

Dimmed colors mean the word is **NOT** in the 10000 most common words.

#### Part of Speech

Adjectives are blue tones, nouns are magenta, verbs are yellow.

The most desirable ones have that as their background color, these will be singular nouns and the base form of verbs.

The ones where this is the foreground color will be plurals, other tenses of verbs, etc.

## create_dataset.py

Script for creating a major system dataset. (this contains a word, the number that word decodes to using the major system, the part of speech of that word, the individual phonemes of the word, and frequency information for that word)

Takes in a wikipedia frequency dataset, see [IlyaSemenov/wikipedia-word-frequency](https://github.com/IlyaSemenov/wikipedia-word-frequency).

Example:
```bash
python create_dataset.py --frequency assets/enwiki-20210820-words-frequency.txt --output assets/major_system_lookup.csv
```

This uses g2p to get the phonemes for the words (this relies on cmudict), and textblob for getting information about the part of speech. Both of these may be inaccurate in some cases.

Running the script on the whole wikipedia dump takes about 9h on my machine, so maybe use one of the provided datasets.

## major-system-converter.py

Experimental python CLI I quickly hacked together to test the dataset.

Example:
```bash
python major_system_converter.py --dataset assets/major_system_lookup_250k.csv
```

## assets/

Contains the latest wikipedia word frequency dataset I could find, as well as precomputed major system datasets created using create_dataset.py.
