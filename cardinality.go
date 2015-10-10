package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/klauspost/shutdown"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var cpu = flag.Int("cpu", runtime.NumCPU(), "cpu")

var parallel = flag.Int("parallel", 25, "parallel")

type Code string

var database = new(struct {
	*sql.DB
	sync.RWMutex
})

func newDB(name string) (err error) {
	database.DB, err = sql.Open("sqlite3", name)
	if err != nil {
		return
	}

	defs := `
		/* 추정에 이용될 비교 대상 데이터 집합 */
		CREATE TABLE IF NOT EXISTS comparison ( Code TEXT UNIQUE );

		/* 부수적으로 얻은 새로운 데이터 집합 */
		CREATE TABLE IF NOT EXISTS new ( Code TEXT UNIQUE );

		CREATE TABLE IF NOT EXISTS card (
			Time DATETIME DEFAULT CURRENT_TIMESTAMP,

			/* 비교 대상 데이터 집합 크기 */
			a INTEGER UNIQUE,

			/* 중복 데이터 갯수 */
			r INTEGER,

			/* 표본 데이터 갯수 */
			N INTEGER
		);

		CREATE TRIGGER IF NOT EXISTS trigger_exist_code BEFORE INSERT ON new
		WHEN new.Code IN comparison
		BEGIN
			INSERT OR REPLACE INTO card SELECT
				coalesce(card.Time, CURRENT_TIMESTAMP),
				comp.a,
				coalesce(card.r, 0) + 1,
				coalesce(card.N, 0) + 1
			FROM
				(SELECT COUNT(*) AS a FROM comparison) AS comp
				LEFT JOIN card ON card.a = comp.a;

			SELECT RAISE(IGNORE);
		END;

		CREATE TRIGGER IF NOT EXISTS trigger_non_exist_code INSERT ON new
		WHEN new.Code NOT IN comparison
		BEGIN
			INSERT OR REPLACE INTO card
			SELECT
				coalesce(card.Time, CURRENT_TIMESTAMP),
				comp.a,
				coalesce(card.r, 0),
				coalesce(card.N, 0) + 1
			FROM
				(SELECT COUNT(*) AS a FROM comparison) AS comp
				LEFT JOIN card ON card.a = comp.a;
		END;
	`

	if _, err = database.Exec(defs); err != nil {
		database.DB.Close()
		return
	}

	return
}

func Contains(set []string, str string) bool {
	for _, el := range set {
		if el == str {
			return true
		}
	}

	return false
}

func CodePool() chan Code {
	codec := make(chan Code)

	for i := 0; i < *parallel; i++ {
		go func() {
			for {
				code, err := One()
				if err == nil {
					codec <- code
				}
			}
		}()
	}

	return codec
}

type engine func(http.ResponseWriter, *http.Request) error

func (fn engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if shutdown.Lock() {
		if err := fn(w, req); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
		shutdown.Unlock()
		return
	}

	w.WriteHeader(http.StatusServiceUnavailable)
}

func AllowOnlyTop(w http.ResponseWriter, req *http.Request) bool {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return false
	}
	return true
}

func UpdateIndexHandler(w http.ResponseWriter, req *http.Request) error {
	if !AllowOnlyTop(w, req) {
		return nil
	}

	var rowcount int

	database.RLock()
	row := database.QueryRow("SELECT COUNT(*) FROM comparison")
	err := row.Scan(&rowcount)
	database.RUnlock()
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "All %10d", rowcount)

	return nil
}

func Update() {
	if _, err := database.Exec(`
		INSERT OR IGNORE INTO comparison (code) SELECT code FROM new;
		DELETE FROM new;
		VACUUM;
	`); err != nil {
		fmt.Println("데이터를 정리할 수 없습니다.")
		fmt.Println(err)
		return
	}

	stmt, _ := database.Prepare("INSERT INTO comparison (code) VALUES (?)")

	http.Handle("/", engine(UpdateIndexHandler))

	codec := CodePool()
	notifier := shutdown.First()

	for {
		select {
		case code := <-codec:
			database.Lock()
			stmt.Exec(string(code))
			database.Unlock()
		case n := <-notifier:
			stmt.Close()
			close(n)
			return
		}
	}
}

func BloatHandler(w http.ResponseWriter, req *http.Request) error {
	if !AllowOnlyTop(w, req) {
		return nil
	}

	var early, newcount int

	database.RLock()
	row := database.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM comparison),
			(SELECT COUNT(*) FROM new)
		;`,
	)
	err := row.Scan(&early, &newcount)
	database.RUnlock()
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "(Bloat) early = %10d new = %10d", early, newcount)

	return nil
}

func Bloat() {
	stmt, err := database.Prepare("INSERT INTO new (code) VALUES (?)")
	if err != nil {
		fmt.Println(err)
		return
	}

	http.Handle("/", engine(BloatHandler))
	http.Handle("/card", engine(CardHandler))

	codec := CodePool()
	notifier := shutdown.First()

	for {
		select {
		case code := <-codec:
			database.Lock()
			stmt.Exec(string(code))
			database.Unlock()
		case n := <-notifier:
			stmt.Close()
			close(n)
			return
		}
	}
}

func CardHandler(w http.ResponseWriter, req *http.Request) error {
	database.RLock()
	rows, err := database.Query("SELECT * FROM card ORDER BY Time DESC")
	database.RUnlock()
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var when time.Time
		var a, r, N int

		if err := rows.Scan(&when, &a, &r, &N); err != nil {
			return nil
		}
		
		fmt.Fprintf(w, "%s a = %10d, r = %10d, N = %10d\n", when, a, r, N)
	}

	return nil
}

func Card() {
	http.Handle("/", engine(BloatHandler))
	http.Handle("/card", engine(CardHandler))

	select{}
}

func usage() {
	fmt.Fprintln(
		os.Stderr,
		os.Args[0], "update | bloat | card [-cpu] [-parallel]",
	)
	flag.PrintDefaults()
}

func main() {
	flag.CommandLine.Usage = usage
	flag.Parse()

	mode := flag.Arg(0)
	if !Contains([]string{
		"update",
		"bloat",
		"card",
	}, mode) {
		usage()
		os.Exit(2)
	}

	if err := newDB("database"); err != nil {
		log.Print(err)
		return
	}

	runtime.GOMAXPROCS(*cpu)
	shutdown.OnSignal(0, os.Interrupt, os.Kill, syscall.SIGTERM)

	shutdown.ThirdFunc(func(interface{}) {
		database.Lock()
		database.Close()
		database.Unlock()
	}, nil)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	switch mode {
	case "update":
		Update()
	case "bloat":
		Bloat()
	case "card":
		Card()
	}

	if !shutdown.Started() {
		shutdown.Shutdown()
	}
}
