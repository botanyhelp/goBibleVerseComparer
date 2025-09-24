package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"bufio"
	"strings"
	"os"
	"regexp"
	"strconv"
	"flag"
	"slices"
	"maps"
	"sort"
)


// RopeSegment represents a segment of the rope.
// In a real rope, this would likely be a more complex struct
// possibly containing the actual string data and length.
type RopeSegment struct {
	Content string
	Length  int
}

// Rope represents the rope data structure.
// This example uses a simplified nested map for demonstration.
type Rope struct {
	// segments: map[segmentID]map[startIndex]map[endIndex]segmentContent
	// This is a highly simplified representation for demonstration purposes.
	// A real rope would use a balanced tree structure.
	Segments map[string]map[int]map[int]string
}

// NewRope creates a new empty Rope.
func NewRope() *Rope {
	return &Rope{
		Segments: make(map[string]map[int]map[int]string),
	}
}

// AddSegment adds a segment to the rope.
// In a real rope, this would involve tree manipulation.
func (r *Rope) AddSegment(segmentID string, startIndex, endIndex int, content string) {
	if _, ok := r.Segments[segmentID]; !ok {
		r.Segments[segmentID] = make(map[int]map[int]string)
	}
	if _, ok := r.Segments[segmentID][startIndex]; !ok {
		r.Segments[segmentID][startIndex] = make(map[int]string)
	}
	r.Segments[segmentID][startIndex][endIndex] = content
}

// GetSegmentContent retrieves the content of a specific segment.
func (r *Rope) GetSegmentContent(segmentID string, startIndex, endIndex int) (string, bool) {
	if segs, ok := r.Segments[segmentID]; ok {
		if startMap, ok := segs[startIndex]; ok {
			if content, ok := startMap[endIndex]; ok {
				return content, true
			}
		}
	}
	return "", false
}

// fetchBibleUrls retrieves the http url argument with http.Get,
// gets entire body with ioutil.ReadAll and returns a map[string][string]
//that holds the title-of-bible mapped to the URL where it can be retrieved
// we want something like this map:
//
// bibleTextTitles := map[string]string{
// 	"Berean Standard Bible":  "https://bereanbible.com/bsb.txt",
// 	"Catholic Public Domain Version": "https://bereanbible.com/cpdv.txt",
// }
//
func fetchBibleUrls(url string) map[string]string {
	//url := "https://openbible.com/textfiles/bibles.txt"
	var debug bool = false

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making HTTP request: %v", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Convert the byte slice to a string and print it
	//if debug { fmt.Println(string(body))}
	var fileContent string = string(body)
	dataMap := make(map[string]string)
    	scanner := bufio.NewScanner(strings.NewReader(fileContent))
    	for scanner.Scan() {
    		line := scanner.Text()
    		parts := strings.SplitN(line, "=", 2) // Split only on the first '='
    		if len(parts) == 2 {
    			key := strings.TrimSpace(parts[0])
    			value := strings.TrimSpace(parts[1])
    			dataMap[key] = value
    		}
    	}

    	if err := scanner.Err(); err != nil {
    		log.Fatalf("Error scanning file content: %v", err)
    	}

    	if debug { fmt.Println("Parsed Map:", dataMap) }
	return dataMap
}

// fetchBibleTextFromUrl retrieves the http url argument with http.Get,
// gets entire body with ioutil.ReadAll and returns the bible text as a string
func fetchBibleTextFromUrl(url string) string {
	//url := "https://openbible.com/textfiles/bsb.txt"
        var debug bool = true
	if debug { fmt.Printf("About to fetch %s\n", url) }
	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making HTTP request: %v", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Convert the byte slice to a string and print it
	//if debug { fmt.Println(string(body))}
	return string(body)
}


// fetchBibleTextFromFile opens the file at the filePath argument
// and returns the contents of the file as a string
func fetchBibleTextFromFile(filePath string) string{
	//filePath := "bsb.txt"
	// Create a dummy file for demonstration
	//err := os.WriteFile(filePath, []byte("Hello World!\nThis is a test file."), 0644)
	//if err != nil {
		//fmt.Errorf("Error creating file: %v\n", err)
		//return 
	//}

	// Read the file content into a byte slice
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Errorf("Error reading file: %v\n", err)
	}

	// Convert the byte slice to a string and return it
	return string(contentBytes)
}

// parseVerse uses regexp library and a hardcoded regular expression 
// to extract and return a slice of four strings from the argument string
// it will operate on lines like this one 
// Genesis 1:1     In the beginning God created the heaven and the earth.
// to extract four strings that represent these entities:
// book chapterNumber verseNumber verse 
func parseVerse(line string) []string {
	pattern := `(.*) ([0-9][0-9]*):([0-9][0-9]*)\t(.*)`
	re := regexp.MustCompile(pattern)
	return re.FindStringSubmatch(line)
}

// readBibleIntoRope takes a string of an entire bible and returns 
// a pointer to our centerpiece data structure: Rope
func readBibleIntoRope(bibleOne string) (*Rope, error) {
	var debug bool = false
	// Create a Rope to hold bible verses
	myRope := NewRope()
	// Create a new reader from the string
	stringReader := strings.NewReader(bibleOne)
	// Create a new scanner from the string reader
	scanner := bufio.NewScanner(stringReader)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		if lineCount>2 {
			var mySliceOfVerseLine []string = parseVerse(line)
			if debug {
				fmt.Println(line)
				fmt.Println("now we print matches\n")
				for _,v := range(mySliceOfVerseLine) {
					fmt.Printf("|%v",v)
				}
				fmt.Println("done matches\n")
			}
			book := mySliceOfVerseLine[1]
			chapterNumber, err := strconv.Atoi(mySliceOfVerseLine[2])
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				return myRope, err
			}
			verseNumber, err := strconv.Atoi(mySliceOfVerseLine[3])
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				return myRope, err
			}
			verse := mySliceOfVerseLine[4]
			myRope.AddSegment(book, chapterNumber, verseNumber, verse)
		}
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading lines: %v\n", err)
	}
	fmt.Printf("We got %d lines\n", lineCount)

	return myRope, nil
}

func main() {
	// debug true will print too much information (got love if you want it -Bob Dylan)
	var debug bool = false
	if debug { fmt.Printf("Mr. Rogers loves you\n")}
	
	var bibleOne string
	//if bibleByFile is true, then you must have the real and hardcoded kjv.txt file in the current directory
	var bibleByFile bool = false
	//if bibleByUrl is true, then you must pick from one of the 14 or so bible URLs in the bibles array shown in a comment somewhere
	var bibleByUrl bool = true
	var bibleTexts []string
	var bibleTitles []string
	
	//var bibleUrls [2]string = [2]string{"https://openbible.com/textfiles/asv.txt","https://openbible.com/textfiles/cpdv.txt"}

        bibleUrls := fetchBibleUrls("http://pennstatehousing.s3-website.us-east-2.amazonaws.com/bibles/bibles.txt")
        if debug { fmt.Printf("bibleUrls: %v\n", bibleUrls)}

	if bibleByUrl {
               for bibleName, bibleURL := range(bibleUrls) {
                       bibleOne = fetchBibleTextFromUrl(bibleURL)
                       bibleTexts = append(bibleTexts, bibleOne)
                       bibleTitles = append(bibleTitles, bibleName)
               }
	}

	var bibleTextFilePaths []string = []string{"testdata/kjv.txt", "testdata/web.txt"}
	if bibleByFile {
		//kjv.txt is entire bible but with first 2 lines are not verses
		//kjv10.txt is first ten verses of bible
		//var myFilePath string = "kjv10.txt"
		for _, myFilePath := range(bibleTextFilePaths) {
			//var myFilePath string = "testdata/kjv.txt"
			bibleOne = fetchBibleTextFromFile(myFilePath)
			bibleTexts = append(bibleTexts, bibleOne)
		}
	}

	var bibleRopes []*Rope
	var err error
	for _, bibleOne := range(bibleTexts) {
		myRope, _ := readBibleIntoRope(bibleOne)
		bibleRopes = append(bibleRopes, myRope)
	}

	var validBooks []string = []string{"1 Chronicles","1 Corinthians","1 John","1 Kings","1 Peter","1 Samuel","1 Thessalonians","1 Timothy","2 Chronicles","2 Corinthians","2 John","2 Kings","2 Peter","2 Samuel","2 Thessalonians","2 Timothy","3 John","Acts","Amos","Colossians","Daniel","Deuteronomy","Ecclesiastes","Ephesians","Esther","Exodus","Ezekiel","Ezra","Galatians","Genesis","Habakkuk","Haggai","Hebrews","Hosea","Isaiah","James","Jeremiah","Job","Joel","John","Jonah","Joshua","Jude","Judges","Lamentations","Leviticus","Luke","Malachi","Mark","Matthew","Micah","Nahum","Nehemiah","Numbers","Obadiah","Philemon","Philippians","Proverbs","Psalm","Revelation","Romans","Ruth","Song of Solomon","Titus","Zechariah","Zephaniah",}
	if debug { fmt.Printf("validBooks are:\n%v\n", validBooks)}
	var book string
	flag.StringVar(&book, "book", "Mark", "the name of the book, Genesis, Mark, Luke, capitalized")
	var chapterNumber int
	flag.IntVar(&chapterNumber, "chapterNumber", 1, "the number of the chapter, like 3 in John 3:16")
	var verseNumber int
	flag.IntVar(&verseNumber, "verseNumber", 1, "the number of the verse, like 16 in John 3:16")

	flag.Parse() // Parse command-line flags
	if debug { fmt.Println("Based on your command-line flags we will look for %s:%d:%d\n", book, chapterNumber, verseNumber) }

	if debug { fmt.Println("Otherwise, enter some text (press Ctrl+D or Ctrl+Z and Enter to finish):") }

	reader := bufio.NewReader(os.Stdin)

	var goodBookYet bool = false
	for !goodBookYet {
		// Prompt for and read the first value
		fmt.Print("Enter the book, like 'Genesis' or '2 Corinthians': ")
		book, _ = reader.ReadString('\n')
		book = strings.TrimSpace(book)
		// Check if the book provided by user i" is in the slice
		if slices.Contains(validBooks, book) {
			goodBookYet = true
			if debug { fmt.Printf("'%s' is in the validBooks slice.\n", book) }
		} else {
			fmt.Printf("%s is NOT in the list of valid books:\n%v\n\n", book, validBooks)
		}
	}


	var chapterNumberString string
	var verseNumberString string

	var goodChapterNumberYet bool = false

	// chapterNumberInt will hold chapter number which we need later to get set of verseNumber in that chapter
	var chapterNumberInt int
	for !goodChapterNumberYet {
		// Prompt for and read the second value
		fmt.Print("Enter the chapter number: ")
		chapterNumberString, _ = reader.ReadString('\n')
		chapterNumberString = strings.TrimSpace(chapterNumberString)
		var firstRope *Rope = bibleRopes[0]
		//var currentBookInRope = 
		// Create a new set of int 
	        chapterSet := make(map[int]bool)
		for outerKey, innerMap := range firstRope.Segments[book] {
			// Add elements to the set
			chapterSet[outerKey] = true
			if debug {
				for innerKey, innerValue := range innerMap {
					fmt.Printf("Outer key: %s, Inner key: %d, Value: %s\n", outerKey, innerKey, innerValue)
				}
			}
		}

		var chapterSetKeys []int = slices.Collect(maps.Keys(chapterSet))
		if debug { fmt.Printf("chapterSetKeys: %v\n", chapterSetKeys) }
		sort.Ints(chapterSetKeys)
		if debug {fmt.Printf("sortedChapterSetKeys: %v\n", chapterSetKeys)}
		//fmt.Printf("sortedChapterSetKeys: %v\n", sort.Ints(slices.Collect(maps.Keys(chapterSet))))
		// Check if the book provided by user is in the set of chapters for that book
		chapterNumberInt, _ = strconv.Atoi(chapterNumberString)
		// Check for membership
		_, ok := chapterSet[chapterNumberInt]
		if ok {
			if debug { fmt.Printf("%v is in the chapterSet %v\n", chapterNumberInt, chapterSetKeys) }
			goodChapterNumberYet = true
		} else {
			//slices.Collect(maps.Keys(someMap))
			//WORKS, but not sorted
			//fmt.Printf("%s is NOT in the list of valid chapters of %s:\n%v\n\n", chapterNumberString,book,slices.Collect(maps.Keys(chapterSet)))
			fmt.Printf("%s is NOT in the list of valid chapters of %s:\n%v\n\n", chapterNumberString,book,chapterSetKeys)
		}
		
	}

	var goodVerseNumberYet bool = false
	for !goodVerseNumberYet {
		// Prompt for and read the third value
		fmt.Print("Enter the verse number: ")
		verseNumberString, _ = reader.ReadString('\n')
		verseNumberString = strings.TrimSpace(verseNumberString)
		var firstRope *Rope = bibleRopes[0]
		//var currentBookInRope = 
		// Create a new set of int 
	        verseSet := make(map[int]bool)
		for outerKey, innerMap := range firstRope.Segments[book] {
			for innerKey, innerValue := range innerMap {
				if outerKey == chapterNumberInt {
					// Add elements to the set
					verseSet[innerKey] = true
					if debug { fmt.Printf("Outer key: %s, Inner key: %d, Value: %s\n", outerKey, innerKey, innerValue)}
				}
			}
		}
		var verseSetKeys []int = slices.Collect(maps.Keys(verseSet))
		if debug { fmt.Printf("verseSetKeys: %v\n", verseSetKeys) }
		sort.Ints(verseSetKeys)
		if debug {fmt.Printf("sortedVerseSetKeys: %v\n", verseSetKeys)}
		//fmt.Printf("sortedVerseSetKeys: %v\n", verseSetKeys)
		//fmt.Printf("sortedChapterSetKeys: %v\n", sort.Ints(slices.Collect(maps.Keys(verseSet))))
		// Check if the book provided by user is in the set of verses for that book
		verseNumberInt, _ := strconv.Atoi(verseNumberString)
		// Check for membership
		_, ok := verseSet[verseNumberInt]
		if ok {
			if debug { fmt.Printf("%v is in the verseSet: %v", verseNumberInt, verseSetKeys) }
			goodVerseNumberYet = true
		} else {
			//slices.Collect(maps.Keys(someMap))
			//WORKS, but not sorted
			//fmt.Printf("%s is NOT in the list of valid verses of %s:\n%v\n\n", verseNumberString,book,slices.Collect(maps.Keys(verseSet)))
			fmt.Printf("%s is NOT in the list of valid verse numbers of %s:%d, and so please enter a verse number from this list:\n%v\n\n", verseNumberString,book,chapterNumberInt,verseSetKeys)
		}
	}

	//var err error
	chapterNumber, err = strconv.Atoi(chapterNumberString)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}
	verseNumber, err = strconv.Atoi(verseNumberString)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}

	// Print the collected values
	fmt.Printf("%s %d:%d\n", book, chapterNumber, verseNumber)
	for bibleIndex, myRope := range(bibleRopes) {
		content, found := myRope.GetSegmentContent(book, chapterNumber, verseNumber)
		if found {
			//fmt.Printf("%s: %s\n", bibleTextFilePaths[bibleIndex], content)
			fmt.Printf("%s:    %s\n", content, bibleTitles[bibleIndex])
		}
	}
}

	// Adding some "segments" to our conceptual rope
	//myRope.AddSegment("2 Samuel", 13, 28, "Now Absalom had commanded his servants, saying, Mark ye now when Amnon’s heart is merry with wine, and when I say unto you, Smite Amnon; then kill him, fear not: have not I commanded you? be courageous, and be valiant.")
	//myRope.AddSegment("1 Kings", 20, 7, "Then the king of Israel called all the elders of the land, and said, Mark, I pray you, and see how this [man] seeketh mischief: for he sent unto me for my wives, and for my children, and for my silver, and for //my gold; and I denied him not.")
	//myRope.AddSegment("Job", 21, 5, "Mark me, and be astonished, and lay [your] hand upon [your] mouth.")
	//myRope.AddSegment("Job", 33, 31, "Mark well, O Job, hearken unto me: hold thy peace, and I will speak.")
	//myRope.AddSegment("Psalm", 37, 37, "Mark the perfect [man], and behold the upright: for the end of [that] man [is] peace.")
	//myRope.AddSegment("Psalm", 48, 13, "Mark ye well her bulwarks, consider her palaces; that ye may tell [it] to the generation following.")
	//myRope.AddSegment("Mark", 1, 1, "The beginning of the gospel of Jesus Christ, the Son of God;")
	//myRope.AddSegment("Mark", 1, 2, "As it is written in the prophets, Behold, I send //my messenger before thy face, which shall prepare thy way before thee.")
	//myRope.AddSegment("Mark", 1, 3, "The voice of one crying in the wilderness, Prepare ye the way of the Lord, make his paths straight.")
	//myRope.AddSegment("Mark", 1, 4, "John did baptize in the wilderness, and preach the baptism of repentance for the remission of sins.")

	// Retrieving segment content
	//content1, found1 := myRope.GetSegmentContent("2 Samuel", 13, 28)
	//if found1 {
	//	fmt.Printf("Segment '2 Samuel': %s\n", content1)
	//}

	//content2, found2 := myRope.GetSegmentContent("Psalm", 48, 13)
	//if found2 {
	//	fmt.Printf("Segment 'Psalm': %s\n", content2)
	//}

	//content3, found3 := myRope.GetSegmentContent("Mark", 1, 3)
	//if found3 {
	//	fmt.Printf("Segment 'Mark': %s\n", content3)
	//}

	//content4, found4 := myRope.GetSegmentContent("Genesis", 1, 10)
	//if found4 {
	//	fmt.Printf("Segment 'Genesis': %s\n", content4)
	//}

//"https://openbible.com/textfiles/bsb.txt",
//"https://openbible.com/textfiles/brb.txt",
//"https://openbible.com/textfiles/asv.txt",
//"https://openbible.com/textfiles/akjv.txt",
//"https://openbible.com/textfiles/cpdv.txt",
//"https://openbible.com/textfiles/dbt.txt",
//"https://openbible.com/textfiles/drb.txt",
//"https://openbible.com/textfiles/erv.txt",
//"https://openbible.com/textfiles/jps.txt",
//"https://openbible.com/textfiles/kjv.txt",
//"https://openbible.com/textfiles/slt.txt",
//"https://openbible.com/textfiles/wvt.txt",
//"https://openbible.com/textfiles/web.txt",
//"https://openbible.com/textfiles/ylt.txt",
//"https://archive.org/download/cuv_20220420/CUV_txt.tar.gz",

//	// these were good URLs in Sept 2025
//	// the first 14 are all similar in this way
//	// 1. they have the first 2 lines of non-verse garbage to discard
//	// 2. they all have one verse line in the same format like this:
//	// #(.*) ([0-9][0-9]*):([0-9][0-9]*)\t(.*)#
//	// ..which might work with golang regexp package
//	// the last one is tar gzip but has good and uniform chinese with 13 lines of non-verse at the top of the file
//	var bibles [15]string = [15]string{"https://openbible.com/textfiles/bsb.txt","https://openbible.com/textfiles/brb.txt","https://openbible.com/textfiles/asv.txt","https://openbible.com/textfiles/akjv.txt","https://openbible.com/textfiles/cpdv.txt","https://openbible.com/textfiles/dbt.txt","https://openbible.com/textfiles/drb.txt","https://openbible.com/textfiles/erv.txt","https://openbible.com/textfiles/jps.txt","https://openbible.com/textfiles/kjv.txt","https://openbible.com/textfiles/slt.txt","https://openbible.com/textfiles/wvt.txt","https://openbible.com/textfiles/web.txt","https://openbible.com/textfiles/ylt.txt","https://archive.org/download/cuv_20220420/CUV_txt.tar.gz"}
//	//LDS
//	//https://github.com/beandog/lds-scriptures/archive/2020.12.08.zip
//	//cp /var/tmp/lds-scriptures-2020.12.08/text/kjv-scriptures.txt ~/goStuff/bibleone/
//	//https://scriptures.nephi.org/mysql
//	//SELECT vol.volume_title, b.book_title, c.chapter_number, v.scripture_text FROM volumes vol JOIN books b on b.volume_id = vol.id JOIN chapters c ON c.book_id = b.id JOIN verses v ON v.chapter_id = c.id WHERE b.book_title = 'John' AND c.chapter_number = 3 AND v.verse_number = 16;
//	//SELECT scripture_text FROM scriptures WHERE verse_title = 'John 3:16';
//	//
//	//https://scriptures.nephi.org/postgresql
//	//SELECT vol.volume_title, b.book_title, c.chapter_number, v.scripture_text FROM volumes vol JOIN books b on b.volume_id = vol.id JOIN chapters c ON c.book_id = b.id JOIN verses v ON v.chapter_id = c.id WHERE b.book_title = 'John' AND c.chapter_number = 3 AND v.verse_number = 16;
//	//SELECT scripture_text FROM scriptures WHERE verse_title = 'John 3:16';
//	for _,bible := range(bibles) {
//		if debug {
//			fmt.Printf("%s\n",bible)
//		}
//	}

// fetchBibleUrls retrieves the http url argument with http.Get,
// gets entire body with ioutil.ReadAll and returns a map[string][string]
//that holds the title-of-bible mapped to the URL where it can be retrieved
// we want something like this map:
//
// bibleTextTitles := map[string]string{
// 	"Berean Standard Bible":  "https://bereanbible.com/bsb.txt",
// 	"Catholic Public Domain Version": "https://bereanbible.com/cpdv.txt",
// }
//
