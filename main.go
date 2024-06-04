package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
)

func GetRaffles(page playwright.Page) []string {
	raffle_boxes, err := page.Locator(".panel-raffle:not(.raffle-entered)").All()
	if err != nil {
		log.Fatal("Could not locate raffle boxes")
	}

	var raffle_refs []string
	for _, rbox := range raffle_boxes {
		raffle := rbox.Locator(".raffle-name a")
		href, err := raffle.GetAttribute("href")

		if err != nil {
			title, _ := raffle.InnerText()
			log.Printf("Failed to get href attribute for %s", title)
		} else {
			raffle_refs = append(raffle_refs, href)
		}
	}

	return raffle_refs
}

func EnterRaffle(page playwright.Page, link string) {
	page.Goto(link)

	err := page.Locator(".enter-raffle-btns").Locator("button:not(#raffle-enter)").Click()
	if err != nil {
		log.Println("Couldn't find 'Enter raffle' button")
	} else {
		text, _ := page.Locator(".subtitle").TextContent()
		log.Printf("Entered to %s", text)
	}
}

func ScrollToEnd(page playwright.Page) {
	for i := 0; i < 3; i++ {
		page.Mouse().Wheel(0, 15000)
		time.Sleep(1 * time.Second)
	}
}

const MainUrl = "https://scrap.tf"

func main() {
	err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
	})
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(false)})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	log.Println("Loading scraptf website")
	if _, err = page.Goto(MainUrl); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	fmt.Println("Press enter when ready :)")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	if _, err = page.Goto(MainUrl + "/raffles"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	ScrollToEnd(page)

	raffles := GetRaffles(page)
	log.Printf("Found %d raffles!", len(raffles))
	for _, raffle_link := range raffles {
		EnterRaffle(page, MainUrl+raffle_link)
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
