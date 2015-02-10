package main

import (
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

var (
	CredFile string
	colorMap map[string]string
	gbot     *gobot.Gobot
	bot      *sphero.SpheroDriver
	reRGB    *regexp.Regexp
	reWord   *regexp.Regexp
)

type MyColor struct {
	R uint8
	G uint8
	B uint8
}

type TwitterCreds struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func init() {
	colorMap = map[string]string{
		"blue":  "0,0,255",
		"red":   "255,0,0",
		"green": "0,255,0",
		"white": "255,255,255",
		"black": "0,0,0",
	}

	flag.StringVar(&CredFile, "creds", "twittercreds.yaml", "path to credential file")

	reRGB = regexp.MustCompile("^@golangphilbot (\\d+) (\\d+) (\\d+)")
	reWord = regexp.MustCompile("@golangphilbot (.+)")

}

func initBot(queue chan *MyColor) {
	gbot = gobot.NewGobot()
	adaptor := sphero.NewSpheroAdaptor("Sphero", "/dev/tty.Sphero-OOR-AMP-SPP")
	bot = sphero.NewSpheroDriver(adaptor, "sphero")

	robot := gobot.NewRobot("sphero",
		[]gobot.Connection{adaptor},
		[]gobot.Device{bot},
		func() {
			for {
				fmt.Printf("bot listening on the channel...\n")
				c := <-queue
				fmt.Printf("robot setting color %v,%v,%v\n", c.R, c.G, c.B)
				SetColor(c.R, c.G, c.B)
			}
		},
	)

	gbot.AddRobot(robot)

	gbot.Start()

}

func SetColor(r, g, b uint8) {
	fmt.Printf("bot SetColor(%v,%v,%v)\n", r, g, b)
	bot.SetRGB(r, g, b)
}

func ReadCredentials(yamlPath string) (*TwitterCreds, error) {
	if yamlPath == "" {
		panic("No twittercreds.yaml file path set")
	}

	f, err := os.Open(yamlPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	creds := &TwitterCreds{}
	if err := yaml.Unmarshal(data, creds); err != nil {
		return nil, err
	}
	return creds, nil
}

// Find the first word after the bot's name
func FindFirstWord(tweet string) string {
	if !reWord.MatchString(tweet) {
		return ""
	}

	parts := strings.Split(tweet, " ")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// If the tweet begins with the bot name, then the second word is either a color name
// or rgb numbers
// @golangphilbot blue
// @golangphilbot 200 100 100
// @golangphilbot 200,100,100 // commas are optional
func ProcessTweet(q chan *MyColor, tweet string) {
	tweet = strings.Replace(tweet, ",", " ", -1)
	result := reRGB.FindStringSubmatch(tweet)
	if result != nil {
		fmt.Printf("result=%v,%v,%v\n", result[1], result[2], result[3])
		r, g, b := result[1], result[2], result[3]
		ru, _ := strconv.ParseUint(r, 10, 8)
		gu, _ := strconv.ParseUint(g, 10, 8)
		bu, _ := strconv.ParseUint(b, 10, 8)

		fmt.Printf("twitter setting component r,g,b=%v,%v,%v onto channel\n", r, g, b)

		q <- &MyColor{
			uint8(ru),
			uint8(gu),
			uint8(bu)}

	} else {
		// search by name and word
		word := FindFirstWord(tweet)
		if word != "" {
			fmt.Printf("twitter Using the word %v\n", word)
			rgb := colorMap[word]
			parts := strings.Split(rgb, ",")
			r, g, b := parts[0], parts[1], parts[2]
			ru, _ := strconv.ParseUint(r, 10, 8)
			gu, _ := strconv.ParseUint(g, 10, 8)
			bu, _ := strconv.ParseUint(b, 10, 8)

			fmt.Printf("twitter setting r,g,b=%v,%v,%v onto channel\n", r, g, b)

			q <- &MyColor{
				uint8(ru),
				uint8(gu),
				uint8(bu)}
		} else {
			fmt.Printf("ProcessTweet no word from %v\n", tweet)
		}
	}
}

func main() {
	var q = make(chan *MyColor)
	go initBot(q)
	fmt.Printf("Sleeping 10 seconds to init the robot\n")
	time.Sleep(time.Second * 10)
	fmt.Printf("awake!")

	creds, err := ReadCredentials("twittercreds.yaml")
	if err != nil {
		fmt.Printf("Error reading credential file %v\n", err)
		return
	}

	fmt.Println("Starting...")
	anaconda.SetConsumerKey(creds.ConsumerKey)
	anaconda.SetConsumerSecret(creds.ConsumerSecret)
	api := anaconda.NewTwitterApi(creds.AccessToken, creds.AccessTokenSecret)

	timeline := make(map[int64]string)

	for {
		fmt.Printf("Before api.GetMentionsTimeline...\n")
		values := url.Values{}
		values.Set("count", "3")
		tweets, err := api.GetMentionsTimeline(values)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return
		}
		fmt.Printf("got %v tweets\n", len(tweets))

		// here we want to range in reverse order, but there isn't one afaik
		// for _, tw := range tweets {
		for i := len(tweets) - 1; i >= 0; i-- {
			tw := tweets[i]
			fmt.Printf("Checking...\n")
			//fmt.Printf("tweet=%v\n", tw)
			if _, exists := timeline[tw.Id]; !exists {

				fmt.Printf("%v %v user=%v Text=%v\n", tw.Id,
					tw.CreatedAt, tw.User.ScreenName, tw.Text)
				timeline[tw.Id] = tw.Text
				ProcessTweet(q, tw.Text)
			}
		}

		time.Sleep(time.Second * 5)
	}
}
