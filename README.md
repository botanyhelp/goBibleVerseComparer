# goBibleVerseComparer

* goBibleVerseComparer is a CLI tool that fetches a list of well-formatted bibles and chooses two randomly for comparison. 
* the user enters three items to operate:
    * Book
    * Chapter Number
    * Verse Number
* ..possibly like this:
    * Genesis
    * 5
    * 5
* ..and the program will return the text of that verse from the two bibles being used
* users can also type these two commands at any time:
    * help
    * quit

## Installation

```
git clone https://github.com/botanyhelp/goBibleVerseComparer.git
cd goBibleVerseComparer
go run main.go
```

* goBibleVerseComparer uses only standard libraries and so the go.mod will be mostly empty
* if you have trouble, then you might edit **go.mod** to change its version to whatever version of go is installed on your system:

```
cat go.mod
go version
vi go.mod
```

## Usage

* users can type **help** or **quit** at any time
* below is a short example of a possible interaction

```
go run main.go 

Type 'quit' or 'help' anytime.
Enter the book, like 'Genesis' or '2 Corinthians': Genesis
Enter the chapter number: 5
Enter the verse number: 5
Genesis 5:5
And all the days that Adam lived were nine hundred and thirty years: and he died.:    American Standard Version
And all the time that passed while Adam lived was nine hundred and thirty years, and then he died.:    Catholic Public Domain Version

Type 'quit' or 'help' anytime.
Enter the book, like 'Genesis' or '2 Corinthians': help

At any prompt you can type anything.  If your entry is unusable, there will be help provided.  For example if you misspell a book, like 'Jon', you will get a list of all the valid book names that you can choose from.  Likewise, if you choose a chapter number is not in the book you chose, or a verse number is not in the chapter, valid numbers will be presented.  You can always type 'quit' or 'help'.



Type 'quit' or 'help' anytime.
Enter the book, like 'Genesis' or '2 Corinthians': John 
Enter the chapter number: 3
Enter the verse number: 16
John 3:16
For God so loved the world, that he gave his only begotten Son, that whosoever believeth on him should not perish, but have eternal life.:    American Standard Version
For God so loved the world that he gave his only-begotten Son, so that all who believe in him may not perish, but may have eternal life.:    Catholic Public Domain Version

Type 'quit' or 'help' anytime.
Enter the book, like 'Genesis' or '2 Corinthians': 1 John
Enter the chapter number: 4
Enter the verse number: 8
1 John 4:8
He that loveth not knoweth not God; for God is love.:    American Standard Version
Whoever does not love, does not know God. For God is love.:    Catholic Public Domain Version

Type 'quit' or 'help' anytime.
Enter the book, like 'Genesis' or '2 Corinthians': Gaga
Gaga is NOT in the list of valid books, which are shown here:
[1 Chronicles 1 Corinthians 1 John 1 Kings 1 Peter 1 Samuel 1 Thessalonians 1 Timothy 2 Chronicles 2 Corinthians 2 John 2 Kings 2 Peter 2 Samuel 2 Thessalonians 2 Timothy 3 John Acts Amos Colossians Daniel Deuteronomy Ecclesiastes Ephesians Esther Exodus Ezekiel Ezra Galatians Genesis Habakkuk Haggai Hebrews Hosea Isaiah James Jeremiah Job Joel John Jonah Joshua Jude Judges Lamentations Leviticus Luke Malachi Mark Matthew Micah Nahum Nehemiah Numbers Obadiah Philemon Philippians Proverbs Psalm Revelation Romans Ruth Song of Solomon Titus Zechariah Zephaniah]


Type 'quit' or 'help' anytime.
Enter the book, like 'Genesis' or '2 Corinthians': quit
God loves you! Goodbye! Terminating program.
```

