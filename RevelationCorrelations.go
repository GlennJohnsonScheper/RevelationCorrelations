// RevelationCorrelations.go

// This GO app will extract together the OT verses that correlate with Book of Revelation verses,
// According to "C:\a\t\The Use of the Old Testament in the Book of Revelation.txt", gotten from
// https://www.pre-trib.org/articles/dr-arnold-fruchtenbaum/message/the-use-of-the-old-testament-in-the-book-of-revelation/read

// Metric: 7:30AM - 5:30PM = 10 hours to write v1 and done.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var kjv = []string{}

var kjvTitles = []string{}

var kjvStarts = []int{}

var bookRefs = make(map[string]int)

var nextRefs = make(map[int]struct{})

var regexpNonDigits = regexp.MustCompile(`[^0-9]+`)

var regexpNewline = regexp.MustCompile(`(\r\n|\r|\n)`)

var regexpNonUSAscii = regexp.MustCompile(`[^ -~]+`)

var regexpOneToThreeDigits = regexp.MustCompile(`^[1-9][0-9]{0,2}$`)

// Input data from Dr. Fructenbaum, manually edited:

// 1. Hand-made into code.

// 2. Added an s on Psalm.

// 3. Changed I II Kings Chronicles...
// "I " -> "First .* "
// "II " -> "Second .* "

// 4. Divided one mis-punct line:
// 	"Isaiah 30:33: Daniel 7:11",

// 5. Split one crossing chapters:
// 	"Deuteronomy 31:30-32:44",

// 6. Two more mis-puncts to fix (both s/b "; "):
// not 1 colon in [ 8 17-18;10 9, 10, 12, 15, 19] //  8:17-18;10:9, 10, 12, 15, 19
// not 1 colon in [ 1 10, 10 14] //  1:10, 10:14

var inputs = []string{
	"Revelation 1:1",
	"Daniel 2:28-29",
	"Revelation 1:4",
	"Isaiah 11:2",
	"Revelation 1:5",
	"Genesis 49:11",
	"Psalms 89:27",
	"Revelation 1:6",
	"Exodus 19:6",
	"Isaiah 61:6",
	"Revelation 1:7",
	"Daniel 7:13",
	"Zechariah 12:10-14",
	"Revelation 1:8",
	"Isaiah 41:4",
	"Revelation 1:12",
	"Exodus 25:37; 37:23",
	"Revelation 1:13",
	"Daniel 7:13; 10:5, 16",
	"Revelation 1:14",
	"Daniel 7:9; 10:6",
	"Revelation 1:15",
	"Ezekiel 1:7, 24; 43:2",
	"Daniel 10:6",
	"Revelation 1:16",
	"Judges 5:31",
	"Isaiah 49:2",
	"Revelation 1:17",
	"Isaiah 41:4; 44:6; 48:12",
	"Daniel 8:17-18; 10:9, 10, 12, 15, 19",
	"Revelation 1:18",
	"Job 3:17",
	"Hosea 13:14",
	"Revelation 2:4",
	"Jeremiah 2:2",
	"Revelation 2:7",
	"Genesis 2:9; 3:22-24",
	"Proverbs 11:30; 13:12",
	"Ezekiel 31:8",
	"Revelation 2:12",
	"Isaiah 49:2",
	"Revelation 2:14",
	"Numbers 25:1-3",
	"Revelation 2:17",
	"Exodus 16:33-34",
	"Isaiah 62:2; 65:15",
	"Revelation 2:18",
	"Daniel 10:6",
	"Revelation 2:20",
	"First .* Kings 16:31-32",
	"Second .* Kings 9:7, 22",
	"Revelation 2:23",
	"Psalms 7:9; 26:2; 28:4",
	"Jeremiah 11:20; 17:10",
	"Revelation 2:27",
	"Psalms 2:7-9",
	"Isaiah 30:14",
	"Jeremiah 19:11",
	"Revelation 3:4",
	"Ecclesiastes 9:8",
	"Revelation 3:5",
	"Exodus 32:32-33",
	"Revelation 3:7",
	"Isaiah 22:22",
	"Revelation 3:9",
	"Isaiah 43:4; 49:23; 60:14",
	"Revelation 3:12",
	"Isaiah 62:2",
	"Ezekiel 48:35",
	"Revelation 3:14",
	"Genesis 49:3",
	"Deuteronomy 21:17",
	"Revelation 3:18",
	"Isaiah 55:1",
	"Revelation 3:19",
	"Proverbs 3:12",
	"Revelation 4:1",
	"Ezekiel 1:1",
	"Revelation 4:2",
	"Isaiah 6:1",
	"Ezekiel 1:26-28",
	"Daniel 7:9",
	"Revelation 4:3",
	"Ezekiel 1:26, 28; 10:1",
	"Revelation 4:5",
	"Exodus 19:16; 25:37",
	"Isaiah 11:2",
	"Ezekiel 1:13",
	"Revelation 4:6",
	"Ezekiel 1:5, 18, 22, 26; 10:1, 12",
	"Revelation 4:7",
	"Ezekiel 1:10; 10:14",
	"Revelation 4:8",
	"Isaiah 6:2-3",
	"Ezekiel 1:18; 10:12",
	"Revelation 4:9",
	"Deuteronomy 32:40",
	"Daniel 4:34; 6:26; 12:7",
	"Revelation 4:11",
	"Genesis 1:1",
	"Revelation 5:1",
	"Ezekiel 2:9-10",
	"Daniel 12:4",
	"Revelation 5:5",
	"Genesis 49:9-10",
	"Isaiah 11:1, 10",
	"Revelation 5:6",
	"Isaiah 11:2",
	"Zechariah 3:8-9; 4:10",
	"Revelation 5:8",
	"Psalms 111:2",
	"Revelation 5:9",
	"Psalms 40:3; 98:1; 144:9; 149:1",
	"Isaiah 42:10",
	"Daniel 5:19",
	"Revelation 5:10",
	"Exodus 19:6",
	"Isaiah 61:6",
	"Revelation 5:11",
	"Daniel 7:10",
	"Revelation 6:2",
	"Zechariah 1:8; 6:3",
	"Revelation 6:4",
	"Zechariah 1:8; 6:2",
	"Revelation 6:5",
	"Zechariah 6:2",
	"Revelation 6:8",
	"Jeremiah 15:2-3; 24:10; 29:17",
	"Ezekiel 14:21",
	"Hosea 13:14",
	"Zechariah 6:3",
	"Revelation 6:12",
	"Isaiah 50:3",
	"Joel 2:10",
	"Revelation 6:13",
	"Isaiah 34:4",
	"Revelation 6:14",
	"Isaiah 34:4",
	"Nahum 1:5",
	"Revelation 6:15",
	"Psalms 48:4-6",
	"Isaiah 2:10-12, 19",
	"Revelation 6:16",
	"Hosea 10:8",
	"Revelation 6:17",
	"Psalms 76:7",
	"Jeremiah 30:7",
	"Nahum 1:6",
	"Zephaniah 1:14-18",
	"Malachi 3:2",
	"Revelation 7:1",
	"Isaiah 11:2",
	"Jeremiah 49:36",
	"Ezekiel 7:2; 37:9",
	"Daniel 7:2",
	"Zechariah 6:5",
	"Revelation 7:3",
	"Ezekiel 9:4-6",
	"Revelation 7:4",
	"Genesis 49:1-28",
	"Revelation 7:9",
	"Leviticus 23:40",
	"Revelation 7:10",
	"Psalms 3:8",
	"Revelation 7:14",
	"Genesis 49:11",
	"Revelation 7:15",
	"Leviticus 26:11",
	"Revelation 7:16",
	"Psalms 121:5-6",
	"Isaiah 49:10",
	"Revelation 7:17",
	"Psalms 23:1-2",
	"Ezekiel 34:23",
	"Revelation 8:3",
	"Psalms 141:2",
	"Revelation 8:4",
	"Psalms 141:2",
	"Revelation 8:5",
	"Ezekiel 10:2",
	"Revelation 8:5-6",
	"Exodus 19:16",
	"Revelation 8:7",
	"Exodus 9:23-24",
	"Psalms 18:13",
	"Isaiah 28:2",
	"Revelation 8:8",
	"Exodus 7:17-19",
	"Revelation 8:10",
	"Isaiah 14:12",
	"Revelation 8:11",
	"Jeremiah 9:15; 23:15",
	"Revelation 8:12",
	"Isaiah 13:10",
	"Revelation 9:1",
	"Isaiah 14:12-14",
	"Revelation 9:2",
	"Genesis 19:28",
	"Exodus 19:8",
	"Revelation 9:3",
	"Exodus 10:12-15",
	"Revelation 9:4",
	"Ezekiel 9:4",
	"Revelation 9:6",
	"Job 3:21",
	"Revelation 9:8",
	"Joel 1:6",
	"Revelation 9:9",
	"Joel 2:5",
	"Revelation 9:11",
	"Job 26:6; 28:22; 31:12",
	"Psalms 88:11",
	"Proverbs 15:11",
	"Revelation 9:14",
	"Genesis 15:18",
	"Deuteronomy 1:7",
	"Joshua 1:4",
	"Revelation 10:1",
	"Ezekiel 1:26-28",
	"Revelation 10:4",
	"Daniel 8:26; 12:4-9",
	"Revelation 10:5",
	"Deuteronomy 32:40",
	"Daniel 12:7",
	"Revelation 10:6",
	"Genesis 1:1",
	"Deuteronomy 32:40",
	"Nehemiah 9:6",
	// "Daniel 12:17", // T.B.D #1 --- Daniel only has 12:1–13
	// On the website, anchor text 12:17 displays verse 12:13:
	"Daniel 12:13",
	"Revelation 10:7",
	"Amos 3:7",
	"Revelation 10:9",
	"Jeremiah 15:16",
	// "Ezekiel 2:8-33", // TBD #2 -- EZ 2 only has 10 verses
	// On the website, anchor text 2:8-33 fetches crazy 2:8-33:33
	// Studying, it should be until 3:3. Split this item in two:
	"Ezekiel 2:8-10",
	"Ezekiel 3:1-3",
	"Revelation 10:11",
	"Ezekiel 37:4, 9",
	"Revelation 11:1",
	"Ezekiel 40:3-4",
	"Zechariah 2:1-2",
	"Revelation 11:2",
	"Ezekiel 40:17-20",
	"Revelation 11:4",
	"Zechariah 4:1-3, 11-14",
	"Revelation 11:5",
	"Numbers 16:35",
	"Second .* Kings 1:10-12",
	"Revelation 11:6",
	"Exodus 7:19-25",
	"First .* Kings 17:1",
	"Revelation 11:7",
	"Exodus 7:3, 7, 8, 21",
	"Revelation 11:8",
	"Isaiah 1:9-10; 3:9",
	"Jeremiah 23:14",
	"Ezekiel 16:49",
	"Ezekiel 23:3, 8, 19, 27",
	"Revelation 11:9",
	"Psalms 79:2-3",
	"Revelation 11:11",
	"Ezekiel 37:9-10",
	"Revelation 11:15",
	"Exodus 15:18",
	"Daniel 2:44-45; 7:13-14, 27",
	"Revelation 11:18",
	"Psalms 2:1-3; 46:6; 115:13",
	"Revelation 12:1",
	"Genesis 37:9-11",
	"Revelation 12:2",
	"Isaiah 26:17; 66:7",
	"Micah 4:9-10",
	"Revelation 12:3",
	"Isaiah 27:1",
	"Daniel 7:7, 20, 24",
	"Revelation 12:4",
	"Daniel 8:10",
	"Revelation 12:5",
	"Psalms 2:8-9",
	"Isaiah 66:7",
	"Revelation 12:7",
	"Daniel 10:13, 21; 12:1",
	"Revelation 12:9",
	"Genesis 3:1",
	"Job 1:6; 2:1",
	"Zechariah 3:1",
	"Revelation 12:10",
	"Job 1:9-11; 2:4-5",
	"Zechariah 3:1",
	"Revelation 12:14",
	"Exodus 19:4",
	"Deuteronomy 32:11",
	"Isaiah 40:31",
	"Daniel 7:25; 12:7",
	"Hosea 2:14-15",
	"Revelation 12:15", // serpent cast out of his mouth water as a flood
	// "Hosea 15:10", // T.B.D. #3 -- wrong chapter...
	"Hosea 5:10", // therefore I will pour out my wrath upon them like water
	"Revelation 12:17",
	"Genesis 3:15",
	"Revelation 13:1",
	"Daniel 7:3, 7, 8",
	"Revelation 13:2",
	"Daniel 7:4-6, 8",
	"Revelation 13:3",
	"Daniel 7:8",
	"Revelation 13:4",
	"Daniel 8:24",
	"Revelation 13:5",
	"Daniel 7:8, 11, 20, 25; 11:36",
	"Revelation 13:7",
	"Daniel 7:21",
	"Revelation 13:8",
	"Daniel 12:1",
	"Revelation 13:10",
	"Jeremiah 15:2; 43:11",
	"Revelation 13:11",
	"Daniel 8:3",
	"Revelation 13:13",
	"First .* Kings 1:9-12",
	"Revelation 14:1",
	"Psalms 2:6",
	"Ezekiel 9:4",
	"Revelation 14:2",
	"Ezekiel 1:24; 43:2",
	"Revelation 14:3",
	"Psalms 144:9",
	"Revelation 14:7",
	"Exodus 20:11",
	"Revelation 14:8",
	"Isaiah 21:9",
	"Jeremiah 51:7-8",
	"Revelation 14:10",
	"Genesis 19:24",
	"Psalms 75:8",
	"Isaiah 51:17",
	"Revelation 14:11",
	"Isaiah 34:10; 66:24",
	"Revelation 14:14",
	"Daniel 7:13",
	"Revelation 14:18",
	"Joel 3:13",
	"Revelation 14:19",
	"Isaiah 63:1-6",
	"Revelation 14:20",
	"Joel 3:13",
	"Revelation 15:1",
	"Leviticus 26:21",
	"Revelation 15:3",
	"Exodus 15:1-18",
	"Deuteronomy 31:30", // split from next line
	"Deuteronomy 32:1-44",
	"Psalms 92:5",
	"Psalms 111:2; 139:14",
	"Revelation 15:4",
	"Psalms 86:9",
	"Isaiah 66:23",
	"Jeremiah 10:7",
	"Revelation 15:5",
	"Exodus 38:21",
	"Revelation 15:6",
	"Leviticus 26:21",
	"Revelation 15:7",
	"Jeremiah 25:15",
	"Revelation 15:8",
	"Exodus 40:34-35",
	"Leviticus 26:21",
	"First .* Kings 8:10-11",
	"Second .* Chronicles 5:13-14",
	"Isaiah 6:1-4",
	"Revelation 16:1",
	"Psalms 79:6",
	"Jeremiah 10:25",
	"Ezekiel 22:31",
	"Revelation 16:2",
	"Exodus 9:9-11",
	"Deuteronomy 28:35",
	"Revelation 16:3",
	"Exodus 7:17-25",
	"Revelation 16:4",
	"Exodus 7:17-21",
	"Psalms 78:44",
	"Revelation 16:5",
	"Psalms 145:17",
	"Revelation 16:6",
	"Isaiah 49:26",
	"Revelation 16:7",
	"Psalms 19:9; 145:17",
	"Revelation 16:10",
	"Exodus 10:21-23",
	"Revelation 16:12",
	"Isaiah 11:15-16; 41:2, 25; 46:11",
	"Jeremiah 51:36",
	"Revelation 16:13",
	"Exodus 8:6",
	"Revelation 16:14",
	"First .* Kings 22:21-23",
	"Revelation 16:16",
	"Judges 5:19",
	"Second .* Kings 23:29-30",
	"Second .* Chronicles 35:22",
	"Zechariah 12:11",
	"Revelation 16:19",
	"Jeremiah 25:15",
	"Revelation 16:21",
	"Exodus 9:18-25",
	"Revelation 17:1",
	"Jeremiah 51:13",
	"Nahum 3:4",
	"Revelation 17:2",
	"Isaiah 23:17",
	"Revelation 17:3",
	"Daniel 7:7",
	"Revelation 17:4",
	"Jeremiah 51:7",
	"Ezekiel 28:13",
	"Revelation 17:8",
	"Exodus 32:32-33",
	"Daniel 12:1",
	"Revelation 17:12",
	"Daniel 7:24-25",
	"Revelation 17:16",
	"Leviticus 21:9",
	"Revelation 18:1",
	"Ezekiel 43:2",
	"Revelation 18:2",
	"Isaiah 21:9; 34:13-15",
	"Jeremiah 50:30; 51:37",
	"Revelation 18:3",
	"Jeremiah 51:7",
	"Revelation 18:4",
	"Isaiah 52:11",
	"Jeremiah 50:8; 51:6, 45",
	"Revelation 18:5",
	"Jeremiah 41:9",
	"Revelation 18:6",
	"Psalms 137:8",
	"Jeremiah 50:15, 29",
	"Revelation 18:7",
	"Isaiah 47:7-8",
	"Zephaniah 2:15",
	"Revelation 18:8",
	"Isaiah 47:9",
	"Jeremiah 50:31-32",
	"Revelation 18:9-19",
	"Ezekiel 26:16-18; 27:26-31",
	"Revelation 18:9",
	"Jeremiah 50:46",
	"Revelation 18:10",
	"Isaiah 13:1",
	"Revelation 18:12",
	"Ezekiel 27:12-25",
	"Revelation 18:20",
	"Jeremiah 51:48",
	"Revelation 18:21",
	"Jeremiah 51:63-64",
	"Revelation 18:22",
	"Isaiah 24:8",
	"Jeremiah 25:10",
	"Ezekiel 26:13",
	"Revelation 18:23",
	"Jeremiah 7:34; 16:9; 25:10",
	"Nahum 3:4",
	"Revelation 19:2",
	"Deuteronomy 32:43",
	"Psalms 119:137",
	"Jeremiah 51:48",
	"Revelation 19:3",
	"Isaiah 34:9-10",
	"Jeremiah 51:48",
	"Revelation 19:5",
	"Psalms 22:23; 134:1; 135:1",
	"Revelation 19:6",
	"Psalms 93:1; 97:1",
	"Ezekiel 1:24; 43:2",
	"Daniel 10:6",
	"Revelation 19:11",
	"Psalms 18:10; 45:3-4",
	"Isaiah 11:4-5",
	"Ezekiel 1:1",
	"Revelation 19:13",
	"Isaiah 63:3",
	"Revelation 19:15",
	"Psalms 2:8-9",
	"Isaiah 11:4; 63:3-6",
	"Revelation 19:16",
	"Deuteronomy 10:17",
	"Revelation 19:17",
	"Isaiah 34:6-7",
	"Ezekiel 39:17",
	"Revelation 19:18",
	"Isaiah 34:6-7",
	"Ezekiel 39:18",
	"Revelation 19:19",
	"Psalms 2:2",
	"Joel 3:9-11",
	"Revelation 19:20",
	"Isaiah 30:33",
	"Daniel 7:11",
	"Revelation 19:21",
	"Ezekiel 39:19-20",
	"Revelation 20:2",
	"Genesis 3:1, 13-14",
	"Isaiah 24:21-22",
	"Revelation 20:4",
	"Daniel 7:9, 22, 27; 12:2",
	"Revelation 20:5",
	"Isaiah 26:14",
	"Revelation 20:6",
	"Exodus 19:6",
	"Isaiah 26:19",
	"Revelation 20:8",
	"Ezekiel 38:2; 39:1, 6",
	"Revelation 20:9",
	"Deuteronomy 23:14",
	"Second .* Kings 1:9-12",
	"Ezekiel 38:22; 39:6",
	"Revelation 20:11",
	"Daniel 2:35",
	"Revelation 20:12",
	"Exodus 32:32-33",
	"Psalms 62:12; 69:28",
	"Daniel 7:10",
	"Revelation 20:15",
	"Exodus 32:32-33",
	"Daniel 12:1",
	"Revelation 21:1",
	"Isaiah 65:17; 66:22",
	"Revelation 21:3",
	"Leviticus 26:11-12",
	"Ezekiel 37:27",
	"Revelation 21:4",
	"Isaiah 25:8; 35:10; 51:11; 65:19",
	"Revelation 21:9",
	"Leviticus 26:21",
	"Revelation 21:10",
	"Ezekiel 40:2",
	"Revelation 21:11",
	"Isaiah 60:1-2",
	"Ezekiel 43:2",
	"Revelation 21:12-13",
	"Ezekiel 48:31-34",
	"Revelation 21:15",
	"Ezekiel 40:3, 5",
	"Revelation 21:19-20",
	"Exodus 28:17-20",
	"Isaiah 54:11-12",
	"Revelation 21:23",
	"Isaiah 60:19-20",
	"Revelation 21:24",
	"Isaiah 60:3-5, 16",
	"Revelation 21:25",
	"Isaiah 60:11",
	"Zechariah 14:7",
	"Revelation 21:26",
	"Isaiah 60:5, 16",
	"Revelation 21:27",
	"Isaiah 52:1",
	"Ezekiel 44:9",
	"Zechariah 14:21",
	"Revelation 22:1",
	"Psalms 46:4",
	"Ezekiel 47:1",
	"Zechariah 14:8",
	"Revelation 22:2",
	"Genesis 2:9; 3:22-24",
	"Ezekiel 47:12",
	"Revelation 22:3",
	"Genesis 3:17-19",
	"Zechariah 14:11",
	"Revelation 22:4",
	"Psalms 17:15",
	"Ezekiel 9:4",
	"Revelation 22:5",
	"Isaiah 60:19",
	"Daniel 7:18, 22, 27",
	"Zechariah 14:7",
	"Revelation 22:10",
	"Daniel 8:26; 12:4, 9",
	"Revelation 22:11",
	"Ezekiel 3:27",
	"Daniel 12:10",
	"Revelation 22:12",
	"Psalms 62:12",
	"Isaiah 40:10; 62:11",
	"Revelation 22:13",
	"Isaiah 44:6",
	"Revelation 22:14",
	"Genesis 2:9; 3:22-24",
	"Proverbs 11:30",
	"Revelation 22:15",
	"Deuteronomy 23:18",
	"Revelation 22:18-19",
	"Deuteronomy 4:2; 12:32",
	"Revelation 22:19",
	"Deuteronomy 29:19-20",
}

func main() {
	prepareKjv()
	prepareRefs1()
	prepareRefs2()
	// It's getting there:
	// fmt.Println(bookRefs)
	// map[Amos:73936 Daniel:71610 Deuteronomy:15783 Ecclesiastes:56073 Exodus:4874 Ezekiel:67023 First .* Kings:29152 Genesis:103 Hosea:72995 Isaiah:57138 Jeremiah:61612 Job:43414 Joel:73682 Joshua:18974 Judges:21135 Leviticus:8929 Malachi:76490 Micah:74666 Nahum:75047 Nehemiah:41480 Numbers:11811 Proverbs:53599 Psalms:45999 Revelation:98489 Second .* Chronicles:37440 Second .* Kings:32048 Zechariah:75738 Zephaniah:75410]

	fmt.Println("Extracting KJV verses referenced in 'The Use of the Old Testament in the Book of Revelation' gotten from")
	fmt.Println("<https://www.pre-trib.org/articles/dr-arnold-fruchtenbaum/message/the-use-of-the-old-testament-in-the-book-of-revelation/read>.")

	prepareRefs3() // which will finally act, and make output
}

func prepareKjv() {
	file, err := os.Open("ProjectGutenbergKjvBible.txt")
	// file fetched 2022-07-16 from https://www.gutenberg.org/files/10/10-0.txt
	if err != nil {
		panic(err)
	}
	sb, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	// fmt.Println(len(sb))

	// Kjv is utf-8, just in case, show me why.
	// Yeah, I see some 4000 apostraphes using ... E2 80 99 ...
	// One more other char found in the text title, don't care.
	sb = []byte(strings.Replace(string(sb), "\u2019", "'", -1))
	// fmt.Println(len(sb))
	kjv = regexpNewline.Split(string(sb), -1)
	// fmt.Println(len(kjv))

	// for _, s := range kjv {
	// 	if regexpNonUSAscii.FindStringIndex(s) != nil {
	// 		fmt.Println("NonAscii", s)
	// 	}
	// }

	// lines between about 25..95 tell me exact book title strings.
	for i := 23; i < 98; i++ {
		if kjv[i] == "" {
			continue
		}
		if strings.Index(kjv[i], "Testament") != -1 {
			continue
		}
		kjvTitles = append(kjvTitles, kjv[i])
	}
	// fmt.Println(len(kjvTitles))
	j := 100
	for i := 0; i < len(kjvTitles); i++ {
		found := false
		for ; j < len(kjv); j++ {
			if kjv[j] == "" {
				continue
			}
			if unicode.IsDigit(rune(kjv[j][0])) {
				continue
			}
			// List atop has single space, later a double space!
			lite := strings.Replace(kjv[j], "  ", " ", -1)
			if lite == kjvTitles[i] {
				kjvStarts = append(kjvStarts, j)
				found = true
				break
			}
		}
		if !found {
			panic(kjvTitles[i])
		}
	}
	// fmt.Println(len(kjvStarts))
	// Good! 66 book titles, all found.
}

func prepareRefs1() {
	for _, s := range inputs {
		lh := regexpNonDigits.FindString(s)
		lh = strings.TrimSpace(lh)
		// rh := s[len(lh):]
		// _ = rh
		// fmt.Println(lh, "//", rh)
		// N.B. this big golang newbie risk:
		// Just increment value as if there.
		bookRefs[lh]++
	}

	// show me
	// pq := []string{}
	// for b, n := range bookRefs {
	// 	pq = append(pq, fmt.Sprintf("%3v %v", n, b))
	// }
	// sort.Strings(pq)
	// for _, nb := range pq {
	// 	fmt.Println(nb)
	// }

	//   1 Amos
	//   1 Ecclesiastes
	//   1 Joshua
	//   1 Malachi
	//   1 Micah
	//   1 Nehemiah
	//   2 II Chronicles
	//   2 Judges
	//   2 Numbers
	//   2 Zephaniah
	//   4 II Kings
	//   4 Nahum
	//   4 Proverbs
	//   5 Hosea
	//   5 I Kings
	//   5 Job
	//   6 Joel
	//   8 Leviticus
	//  14 Deuteronomy
	//  17 Zechariah
	//  18 Genesis
	//  30 Exodus
	//  31 Jeremiah
	//  38 Psalm
	//  41 Daniel
	//  50 Ezekiel
	//  71 Isaiah
	// 226 Revelation

	// how naively can I match these names to kjvTitles?
	// Obviously, I and II won't fly.
	// Next func, please...
}

func prepareRefs2() {
	// uses bookRefs from prepareRefs1
	for br, _ := range bookRefs {
		// Does this uniquely agree with just one kjvTitle?
		// 0 II Kings
		// 2 Jeremiah
		// The Book of the Prophet Jeremiah -- take the first one
		// The Lamentations of Jeremiah
		// 0 I Kings
		// 0 II Chronicles
		// 0 Psalm - tweak that in inputData
		// The Book of Psalms
		// Fixed all but Jeremiah, having 2, which is okay.

		regexpName := regexp.MustCompile(`\b` + br + `\b`)

		// show me
		// n := 0
		// for i := 0; i < len(kjvTitles); i++ {
		// 	if regexpName.FindStringIndex(kjvTitles[i]) != nil {
		// 		n++
		// 	}
		// }
		// if n != 1 {
		// 	fmt.Println(n, br)
		// }

		// done counting, act on first:
		for i := 0; i < len(kjvTitles); i++ {
			if regexpName.FindStringIndex(kjvTitles[i]) != nil {
				// MODIFY THE PURPOSE of bookRefs' int:
				bookRefs[br] = kjvStarts[i] // Now holds offset into kjv
				// I also want to stop some loop before any next book
				nextRefs[kjvStarts[i]] = struct{}{}
				break
			}
		}
	}
}

func prepareRefs3() {
	// now study the rh sides
	// After applying changes 1-6,
	// every piece of rh has 1 colon. (ch:verses...)
	for _, s := range inputs {
		lh := regexpNonDigits.FindString(s)
		lh = strings.TrimSpace(lh)
		rh := s[len(lh):]
		// fmt.Println(rh)
		ss := strings.Split(rh, "; ")
		for _, piece := range ss {
			cv := strings.Split(piece, ":")
			if len(cv) != 2 {
				fmt.Println("not 1 colon in", cv, "//", rh)
			}

			// Check all chapter formatting
			ch := strings.TrimSpace(cv[0])
			if !regexpOneToThreeDigits.MatchString(ch) {
				fmt.Println("not 1-3 ch", ch)
			}

			// Check all verse... formatting
			// Now all parts are 1-3 digits,
			// mixed with some "-" and ", ".
			vs := cv[1]
			// fmt.Println(vs)

			// Prepare to loop over such
			// split optional commas out
			oco := strings.Split(vs, ", ")
			for _, nc := range oco {
				// split hyphen ranges
				sv := strings.Split(nc, "-")
				// What's left must parse:
				// Good, all are ints now.
				// for _, v := range sv {
				// 	n, err := strconv.Atoi(v)
				// 	if err != nil {
				// 		panic(err)
				// 	}
				// 	_ = n
				// }
				switch len(sv) {

				// 11-th hour insight:
				// Instead of looping over any range,
				// pass the whole range. Even better!

				case 1:
					// v, err := strconv.Atoi(sv[0])
					// if err != nil {
					// 	panic("sv[]")
					// }
					// actUpon(lh, ch, strconv.Itoa(v))

					actUpon(lh, ch, sv[0], sv[0]) // yes, both are 0.

					break
				case 2:
					// range
					// v0, err := strconv.Atoi(sv[0])
					// if err != nil {
					// 	panic("sv[0]")
					// }
					// v1, err := strconv.Atoi(sv[1])
					// if err != nil {
					// 	panic("sv[1]")
					// }
					// for v := v0; v < v1; v++ {
					// 	actUpon(lh, ch, strconv.Itoa(v))
					// }

					actUpon(lh, ch, sv[0], sv[1])

					break
				default:
					panic("len(sv)")
					break
				}

			}
		}
	}
}

/* I forgot, does this KJV text have 1 newline per verse?

...NO....!

For one thing, they joined verses into paragraphs!
For another thing, they put newlines in paragraphs!
I do not quickly see any paragraphs without leading ch:vs.

==============
...help meet for him.

2:21 And the LORD God caused a deep sleep to fall upon Adam, and he
slept: and he took one of his ribs, and closed up the flesh instead
thereof; 2:22 And the rib, which the LORD God had taken from man, made
he a woman, and brought her unto the man.

2:23 And Adam said...
==============
Oh, That's ugly to realize late on.
Perhaps it suggests an opportunity.
Go back to passing in whole ranges.
*/

func actUpon(name string, ch string, v1 string, v2 string) {

	// My insight is to keep entire paragraphs.
	// So outer loop will go by blocks with "".

	inText := false
	atopParagraph := 0

	// this cannot match any sooner, at shorter ch, nor elsewhere:
	chapter := ch + ":"
	foundChapter := false

	verse1 := chapter + v1
	foundV1 := false
	atopV1 := 0

	verse2 := chapter + v2
	foundV2 := false
	pastV2 := 0

	for i := bookRefs[name] + 1; i < len(kjv); i++ {
		if kjv[i] == "" {
			// blank line separating paragraphs
			if inText && foundV2 {
				pastV2 = i
				break
			}
			inText = false
		} else {
			// non-blank line of text
			if !inText {
				atopParagraph = i
			}
			inText = true
			if !foundChapter {
				if strings.Index(kjv[i], chapter) != -1 {
					foundChapter = true
					// fmt.Println(kjv[i])
				}
			}
			if foundChapter && !foundV1 {
				if strings.Index(kjv[i], verse1) != -1 {
					foundV1 = true
					atopV1 = atopParagraph
				}
			}
			if foundV1 && !foundV2 {
				if strings.Index(kjv[i], verse2) != -1 {
					foundV2 = true
					// but past awaits blank line
				}
			}
		}

		// stop before any next book (except rev. stop atop loop)
		_, ok := nextRefs[i]
		if ok {
			break
		}
	}

	if v1 == v2 {
		fmt.Printf("\r\n%v %v:%v\r\n\r\n", name, ch, v1)
	} else {
		fmt.Printf("\r\n%v %v:%v-%v\r\n\r\n", name, ch, v1, v2)
	}

	if foundV2 {
		// print from [atop to past).
		// Windows \r\n
		for i := atopV1; i < pastV2; i++ {
			fmt.Printf("%v\r\n", kjv[i])
		}
	} else {
		panic("!foundV2")
	}
}
