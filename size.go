package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var tr = &http.Transport{
	MaxIdleConnsPerHost: 12,
}

type Code string

func One() (Code, error) {
	const url = "http://bgmstore.net/random"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := tr.RoundTrip(req)
	if err != nil {
		return "", err
	}

	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	loc := resp.Header.Get("Location")
	code, err := unwrapLoc(loc)

	return Code(code), err
}

// e.g /view/rxBQI/random
// 위 Location 값에서 가운데에 위치한 rxBQI 값을 추출한다.
func unwrapLoc(loc string) (string, error) {
	i := len("/view/")
	if len(loc) < i {
		return "", errors.New("too busy")
	}

	loc = loc[i:]

	// rxBQI 에 해당하는 문자의 길이가 항상 5 글자인 것 같으나 보장되지 않았으므로
	// 안전하게 "/"" 이전 까지를 반환한다.
	i = strings.IndexRune(loc, '/')

	return loc[:i], nil
}

type Some []Code

// Meet 는 Some 을 One() 의 반환값으로 채운다.
//
// n 개의 고루틴을 사용하여 병렬적으로 Some 을 채워나간다. n 이 1 이라면 Some 이
// 순차적으로 채워지겠지만, 그렇지 않다면 순차적으로 채워지지 않는다.
func (some Some) Meet(n int) error {
	errc, returnc := make(chan error), make(chan bool)
	q := make(chan int)

	for i := 0; i < n; i++ {
		go func() {
			var err error
			for {
				select {
				case i := <-q:
					some[i], err = One()
					errc <- err
				case <-returnc:
					return
				}
			}
		}()
	}

	go func() {
		for i := 0; i < len(some); i++ {
			select {
			case q <- i:
			case <-returnc:
				return
			}
		}
	}()

	defer close(returnc)

	for i := 0; i < len(some); i++ {
		if err := <-errc; err != nil {
			return err
		}
	}

	return nil
}

func TestMeet(size, n int) {
	fmt.Printf("size = %4d, n = %2d: ", size, n)
	some := make(Some, size)
	w := time.Now()
	err := some.Meet(n)
	if err == nil {
		fmt.Println(time.Since(w))
	} else {
		fmt.Println(err)
	}
}

func ExitHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "success")
	PrintCounter(w)
	Exit()
}

func StatusHandler(w http.ResponseWriter, req *http.Request) {
	if err := PrintCounter(w); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintln(w, "mode =", os.Args[1])
}

func SizeHandler(w http.ResponseWriter, req *http.Request) {
	if err := PrintCounter(w); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	datasize, _, err := ReadSize()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	counter.RLock()
	L := float64(datasize)
	R := float64(counter.err)
	m := float64(counter.add + counter.err)
	counter.RUnlock()

	if R == 0 {
		io.WriteString(w, "현재 상태에서는 추정이 불가능합니다. 조금 기다려주세요.")
		return
	}

	for _, pair := range NDT {
		a, b := Quess(pair[0], m, R, L)
		fmt.Fprintf(w, "k = %01.2f, p = %01.4f, [%15.6f, %15.6f]\n",
			pair[0], pair[1], a, b)
	}
}

func fetch(codec chan Code) error {
	adder, err := database.Prepare("INSERT INTO data (code) VALUES (?)")
	if err != nil {
		return err
	}

	if err := MoveStay(); err != nil {
		return err
	}

	for {
		select {
		case code := <-codec:
			databaseMu.Lock()
			_, err := adder.Exec(string(code))
			databaseMu.Unlock()

			if err == nil {
				counter.Lock()
				counter.add += 1
				counter.Unlock()
			} else {
				counter.Lock()
				counter.err += 1
				counter.Unlock()
			}

		case <-exitc:
			return nil
		}
	}
}

func size(codec chan Code) error {
	adder, err := database.Prepare("INSERT INTO stay (code) VALUES (?)")
	if err != nil {
		return err
	}

	selecter, err := database.Prepare("SELECT coalesce((SELECT 1 FROM data WHERE code = ?), 0)")
	if err != nil {
		return err
	}

	for {
		select {
		case code := <-codec:
			var in bool
			databaseMu.Lock()
			err := selecter.QueryRow(string(code)).Scan(&in)
			databaseMu.Unlock()
			if err == nil {
				counter.Lock()
				if in {
					counter.err += 1
				} else {
					databaseMu.Lock()
					adder.Exec(string(code))
					databaseMu.Unlock()
					counter.add += 1
				}
				counter.Unlock()
			}

		case <-exitc:
			// MoveStay()
			return nil
		}
	}
}

func ReadSize() (int, int, error) {
	var datasize, staysize int
	databaseMu.RLock()
	err := database.QueryRow(
		"SELECT (SELECT COUNT(*) FROM data), (SELECT COUNT(*) FROM stay)",
	).Scan(
		&datasize,
		&staysize,
	)
	databaseMu.RUnlock()
	return datasize, staysize, err
}

func MoveStay() error {
	databaseMu.Lock()
	err := func() error {
		if _, err := database.Exec(
			"INSERT OR IGNORE INTO data (code) SELECT code FROM stay",
		); err != nil {
			return err
		}
		_, err := database.Exec("DELETE FROM stay")
		return err
	}()
	databaseMu.Unlock()

	return err
}

var exitc = make(chan bool)

var exitOnce sync.Once

func Exit() {
	exitOnce.Do(func() {
		close(exitc)
	})
}

var counter = new(struct {
	add uint64
	err uint64
	sync.RWMutex
})

func PrintCounter(w io.Writer) error {
	datasize, staysize, err := ReadSize()
	if err != nil {
		return err
	}

	counter.RLock()

	fmt.Fprintf(w, "add = %8d\n", counter.add)
	fmt.Fprintf(w, "err = %8d\n", counter.err)
	fmt.Fprintf(w, " +  = %8d\n", counter.add+counter.err)
	fmt.Fprintf(w, "datasize = %8d\n", datasize)
	fmt.Fprintf(w, "staysize = %8d\n", staysize)
	counter.RUnlock()

	return nil
}

var database *sql.DB

var databaseMu sync.RWMutex

func openDB() (err error) {
	database, err = sql.Open("sqlite3", "database")
	return
}

func usage() {
	fmt.Println(os.Args[0], " fetch | size [add err]")
	fmt.Println("    fetch")
	fmt.Println("      서버로부터 새로운 코드를 받아 저장소에 추가합니다.")
	fmt.Println("    size")
	fmt.Println("      서버가 가지고 있는 코드의 수를 추정합니다.")
	fmt.Println("      추정의 정확도는 저장소의 크기와 프로그램의 작동 시간에 비레합니다.")
	fmt.Println("        add")
	fmt.Println("          add 값을 설정합니다.")
	fmt.Println("        err")
	fmt.Println("          err 값을 설정합니다.")
	fmt.Println()
	// fmt.Println("  ")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := openDB(); err != nil {
		log.Print(err)
		return
	}

	defer database.Close()

	for _, query := range []string{
		"CREATE TABLE IF NOT EXISTS data ( code text unique );",
		"CREATE TABLE IF NOT EXISTS stay ( code text unique );",
		"CREATE INDEX IF NOT EXISTS index_data_code ON data(code);",
		"CREATE INDEX IF NOT EXISTS index_stay_code ON stay(code);",
	} {
		if _, err := database.Exec(query); err != nil {
			log.Print(err)
			return
		}
	}

	if len(os.Args[1:]) < 1 {
		usage()
		return
	}

	http.HandleFunc("/exit", ExitHandler)
	http.HandleFunc("/status", StatusHandler)
	http.HandleFunc("/size", SizeHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	codec := make(chan Code)

	for i := 0; i < 25; i++ {
		go func() {
			for {
				code, err := One()
				if err == nil {
					codec <- code
				}
			}
		}()
	}

	if err := func() error {
		if os.Args[1] == "fetch" {
			return fetch(codec)
		}

		if os.Args[1] == "size" {
			if len(os.Args[2:]) >= 2 {
				counter.Lock()
				err := func() (err error) {
					counter.add, err = strconv.ParseUint(os.Args[2], 10, 64)
					if err != nil {
						return
					}

					counter.err, err = strconv.ParseUint(os.Args[3], 10, 64)
					if err != nil {
						return
					}

					return
				}()
				counter.Unlock()
				if err != nil {
					return err
				}
			}
			return size(codec)
		}

		return nil
	}(); err != nil {
		fmt.Println(err)
	}
}
