package kee

import(
    "strings"
    "runtime"
    "path"
    "os"
    "log"
    "bufio"
    "errors"
)

type jumAdjectives struct {
    files []string
    words [][]string
}
type jumNouns struct {    
    files []string
    words [][]string
}
type jumVerbs struct {
    files []string
    words [][]string
}
type jumAdverbs struct {
    files []string
    words [][]string
}
type jumWord interface {
    readWords()
    getWords(syl int) []string
    randomWord(syl int) string
}
type jumCtrl struct {
    phrase []jumWord
    syls []int
}

func (j *jumCtrl) babble(syls []int) (string, uint64) {
    var (
        res string = ""
        tmp string = ""
        space uint64 = 1
    )

    for k, w := range j.phrase {
        tmp = w.randomWord(syls[k])
        space = space * uint64(len(w.getWords(syls[k])))
        if len(tmp) > 1 { tmp = strings.ToUpper(tmp[:1]) + tmp[1:] }
        res += tmp
    }

    return res, space
}

// Generates a random phrase and returns jumble "object". 
// Takes number of syllables for adjective, noun, verb, adverb respectively. 
// Pass 0 for syllable count to skip word.
func (j *jumCtrl) New(sylAdj, sylNoun, sylVerb, sylAdv int) (jumble, error) {
    syls := []int{sylAdj, sylNoun, sylVerb, sylAdv}
    for _, s := range syls {
        if s < 0 || s > 4 { return jumble{}, errors.New("Bad syllable count") }
    }
    phrase, space := j.babble(syls)
    return jumble{phrase, space}, nil
}

type jumble struct {
    phrase string
    space uint64
}

func (m jumble) String() string {
    return m.phrase
}

func (m jumble) SampleSpace() uint64 {
    return m.space
}

func (adj *jumAdjectives) readWords() {
    adj.files = []string{
        "",
        "1syllableadjectives.txt",
        "2syllableadjectives.txt",
        "3syllableadjectives.txt",
        "4syllableadjectives.txt", 
    }
    adj.words = [][]string{
        []string{""},   // 1
        []string{},     // 689
        []string{},     // 5187
        []string{},     // 6924
        []string{},     // 5301
    }
    for i := 1; i < 5; i++ {
        adj.words[i] = jumReadFile("adjectives", adj.files[i])
    }
}

func (noun *jumNouns) readWords() {
    noun.files = []string{
        "",
        "1syllablenouns.txt",
        "2syllablenouns.txt",
        "3syllablenouns.txt",
        "4syllablenouns.txt",
    }
    noun.words = [][]string{
        []string{""},   // 1
        []string{},     // 5865
        []string{},     // 22110
        []string{},     // 20602
        []string{},     // 12247
    }
    for i := 1; i < 5; i++ {
        noun.words[i] = jumReadFile("nouns", noun.files[i])
    }
}

func (verb *jumVerbs) readWords() {
    verb.files = []string{
        "",
        "1syllableverbs.txt",
        "2syllableverbs.txt",
        "3syllableverbs.txt",
        "4syllableverbs.txt",
    }
    verb.words = [][]string{
        []string{""},   // 1
        []string{},     // 3719
        []string{},     // 8568
        []string{},     // 6365
        []string{},     // 3986
    } 
    for i := 1; i < 5; i++ {
        verb.words[i] = jumReadFile("verbs", verb.files[i])
    }
}

func (adv *jumAdverbs) readWords() {
    adv.files = []string{
        "",
        "1syllableadverbs.txt",
        "2syllableadverbs.txt",
        "3syllableadverbs.txt",
        "4syllableadverbs.txt",
    }
    adv.words = [][]string{
        []string{""},   // 1
        []string{},     // 168
        []string{},     // 769
        []string{},     // 1545
        []string{},     // 1428
    }
    for i := 1; i < 5; i++ {
        adv.words[i] = jumReadFile("adverbs", adv.files[i])
    }
}

func (w jumAdjectives) getWords(syl int) []string {
    if w.words == nil { return nil }
    return w.words[syl]
}

func (w jumNouns) getWords(syl int) []string {
    if w.words == nil { return nil }
    return w.words[syl]
}

func (w jumVerbs) getWords(syl int) []string {
    if w.words == nil { return nil }
    return w.words[syl]
}

func (w jumAdverbs) getWords(syl int) []string {
    if w.words == nil { return nil }
    return w.words[syl]
}

func (w *jumAdjectives) randomWord(syl int) string {
    return jumRandomWord(w, syl)
}

func (w *jumNouns) randomWord(syl int) string {
    return jumRandomWord(w, syl)
}

func (w *jumVerbs) randomWord(syl int) string {
    return jumRandomWord(w, syl)
}

func (w *jumAdverbs) randomWord(syl int) string {
    res := jumRandomWord(w, syl)
    return res
}


// -- Helpers -- 

func jumRandomWord(w jumWord, syl int) string {
    if w.getWords(syl) == nil { 
        w.readWords() 
    }
    dict := w.getWords(syl)
    lim := len(dict)
    idx := randIntn(lim)
    res := string(dict[idx])
    return res
}

func jumReadFile(subdir, fn string) []string {
    _, cwd, _, _ := runtime.Caller(1)
    dir := path.Join(path.Dir(cwd), "words", subdir, fn)
    file, err := os.Open(dir)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    words := []string{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        words = append(words, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return words
}