package main

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/jefedavis/resume-pdf-converter/convert"
)

func main() {
	app := fiber.New()

	app.Get("/convert", func(c *fiber.Ctx) error {
		q := c.Query("url", os.Getenv("TARGET_URL"))
		file, err := ioutil.TempFile("", "resume-*.pdf")
		if err != nil {
			return err
		}

		defer os.Remove(file.Name())

		u, err := url.Parse(q)
		if err != nil {
			log.Fatal("unable to parse url: %w", err)
		}

		buf := convert.ConvertHTMLToPDF(u)

		if err := ioutil.WriteFile(file.Name(), buf, 0o644); err != nil {
			return err
		}

		return c.Download(file.Name())
	})

	app.Listen(":3000")
}
