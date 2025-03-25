# Go Web Scraper

A **command-line tool** written in Go for **scraping** (downloading and parsing) HTML pages, extracting key information (like headlines and links), and exporting the results as structured data (JSON).

## Features

- **Configurable URL**: Scrape any website by passing the `-url` flag.
- **Custom Output**: Save the parsed data to a specified file with `-out`.
- **Automatic JSON Export**: Data is marshaled into a JSON array of objects with `title` and `url` fields.
- **Minimal Dependencies**: Uses only the Go standard library plus [goquery](https://github.com/PuerkitoBio/goquery) for HTML parsing.
- **Easy to Extend**: The `scrape` function can be customized to extract other elements (images, dates, etc.).

## Installation

1. **Clone** this repository:
   ```bash
   git clone https://github.com/LLR1/go-web-scraper.git
   cd go-web-scraper
