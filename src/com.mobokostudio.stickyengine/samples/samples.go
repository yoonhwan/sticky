package samples

import (
	"fmt"
	"time"
	"appengine"
	"appengine/user"
	"appengine/memcache"
	
	"net/http"
)

func WebSample() {
    http.HandleFunc("/", userHandler)
    http.HandleFunc("/memcache", memcachedHandler)
    
    http.ListenAndServe(":4000",nil)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
    
    c := appengine.NewContext(r)
    u := user.Current(c)
    if u == nil {
        url, err := user.LoginURL(c, r.URL.String())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Location", url)
        w.WriteHeader(http.StatusFound)
        return
    }
    fmt.Fprintf(w, "Hello, %v!", u)
}

func memcachedHandler(w http.ResponseWriter, r *http.Request) {
	
	c := appengine.NewContext(r)
    // Create an Item
	item := &memcache.Item{
		Key:   "lyric",
		Value: []byte("Oh, give me a home"),
	}
	// Add the item to the memcache, if the key does not already exist
	if err := memcache.Add(c, item); err == memcache.ErrNotStored {
		c.Infof("item with key %q already exists", item.Key)
	} else if err != nil {
		c.Errorf("error adding item: %v", err)
	}

	// Change the Value of the item
	item.Value = []byte("Where the buffalo roam")
	// Set the item, unconditionally
	if err := memcache.Set(c, item); err != nil {
		c.Errorf("error setting item: %v", err)
	}

	// Get the item from the memcache
	if item, err := memcache.Get(c, "lyric"); err == memcache.ErrCacheMiss {
		c.Infof("item not in the cache")
	} else if err != nil {
		c.Errorf("error getting item: %v", err)
	} else {
		c.Infof("the lyric is %q", item.Value)
	}
}

func GoroutineSample() {
	defer func() {
		if str := recover(); str != nil {
			fmt.Println(str)
		}
	}()

	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		for {
			c1 <- "from 1"
			time.Sleep(time.Second * 2)
		}
	}()
	go func() {
		for {
			c2 <- "from 2"
			time.Sleep(time.Second * 3)
		}
	}()
	go func() {
		for {
			select {
			case msg1 := <-c1:
				fmt.Println(msg1)
			case msg2 := <-c2:
				fmt.Println(msg2)
			case xxx, ok := <-time.After(time.Second):
				fmt.Println(xxx, ok, " ", "timeout")

			}
		}
	}()
	var input string
	fmt.Scanln(&input)
}
