package mongo

import "github.com/globalsign/mgo"

var session *mgo.Session
var mgoConfig = defaultConfig()

type config struct {
	IsInit bool
	URL    string
	*mgo.Credential
}

func defaultConfig() *config {
	return &config{
		IsInit: true,
		URL:    "localhost",
		Credential: &mgo.Credential{
			Username: "root",
			Password: "root",
		},
	}
}

func C(cname string) *mgo.Collection {
	if err := session.Ping(); err != nil {
		session.Close()
		Redial()
	}
	return session.DB("fate").C(cname)
}

func Redial() {
	s, err := mgo.Dial(mgoConfig.URL)
	if err != nil {
		panic(err)
	}

	s.Login(mgoConfig.Credential)
	// Optional. Switch the session to a monotonic behavior.
	s.SetMode(mgo.Monotonic, true)
	session = s
}

func Dial(url string, c *mgo.Credential) {
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	credential := mgoConfig.Credential
	if c != nil {
		credential = c
	}
	s.Login(credential)
	// Optional. Switch the session to a monotonic behavior.
	s.SetMode(mgo.Monotonic, true)
	if mgoConfig.IsInit {
		mgoConfig.URL = url
		mgoConfig.Credential = c
		mgoConfig.IsInit = false
	}
	session = s
}

func InsertIfNotExist(c *mgo.Collection, v interface{}) error {
	count, err := c.Find(v).Count()
	if err != nil || count != 0 {
		return err
	}
	err = c.Insert(v)
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	session.Close()
}
