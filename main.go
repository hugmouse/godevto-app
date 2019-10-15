package main

import (
    "fmt"
    dev "github.com/hugmouse/godevto"
    "github.com/mitchellh/go-wordwrap"
    "image"
    "image/color"
    "log"
    "os/exec"
    "runtime"
    "strconv"
    "time"

    "github.com/aarzilli/nucular"
    "github.com/aarzilli/nucular/style"
)

// All empty while not loaded
var (
    response           = dev.Articles{}
    postCounter        = 0
    comments           = ""
    commentsAmount     = ""
    rating             = ""
    ratingAmount       = ""
    DatePublishedTitle = ""
    DatePublished      = ""
)

// post represents an article
var post = dev.Article{
    TypeOf:                 "",
    ID:                     0,
    Title:                  "Click on \"Get published articles\" button!",
    Description:            "Description label. Please see text above",
    CoverImage:             "",
    Published:              false,
    PublishedAt:            time.Time{},
    TagList:                "",
    Tags:                   nil,
    Slug:                   "",
    Path:                   "",
    URL:                    "",
    CanonicalURL:           "",
    CommentsCount:          0,
    PositiveReactionsCount: 0,
    PublishedTimestamp:     time.Time{},
    User:                   dev.User{},
    Organization:           dev.Organization{},
    FlareTag:               dev.FlareTag{},
}

// devtoThemeTable is just a color theme for app
var devtoThemeTable = style.ColorTable{
    ColorText:                  color.RGBA{253, 252, 252, 255},
    ColorWindow:                color.RGBA{0, 0, 0, 255},
    ColorHeader:                color.RGBA{0, 0, 0, 255},
    ColorHeaderFocused:         color.RGBA{0, 0, 0, 255},
    ColorBorder:                color.RGBA{255, 20, 147, 255},
    ColorButton:                color.RGBA{0, 0, 0, 255},
    ColorButtonHover:           color.RGBA{29, 1, 1, 255},
    ColorButtonActive:          color.RGBA{255, 102, 168, 255},
    ColorToggle:                color.RGBA{0, 0, 0, 255},
    ColorToggleHover:           color.RGBA{0, 0, 0, 255},
    ColorToggleCursor:          color.RGBA{0, 0, 0, 255},
    ColorSelect:                color.RGBA{0, 0, 0, 255},
    ColorSelectActive:          color.RGBA{0, 0, 0, 255},
    ColorSlider:                color.RGBA{0, 0, 0, 255},
    ColorSliderCursor:          color.RGBA{0, 0, 0, 255},
    ColorSliderCursorHover:     color.RGBA{0, 0, 0, 255},
    ColorSliderCursorActive:    color.RGBA{0, 0, 0, 255},
    ColorProperty:              color.RGBA{0, 0, 0, 255},
    ColorEdit:                  color.RGBA{0, 0, 0, 255},
    ColorEditCursor:            color.RGBA{0, 0, 0, 255},
    ColorCombo:                 color.RGBA{0, 0, 0, 255},
    ColorChart:                 color.RGBA{0, 0, 0, 255},
    ColorChartColor:            color.RGBA{0, 0, 0, 255},
    ColorChartColorHighlight:   color.RGBA{0, 0, 0, 255},
    ColorScrollbar:             color.RGBA{255, 20, 147, 255},
    ColorScrollbarCursor:       color.RGBA{0, 0, 0, 255},
    ColorScrollbarCursorHover:  color.RGBA{0, 0, 0, 255},
    ColorScrollbarCursorActive: color.RGBA{0, 0, 0, 255},
    ColorTabHeader:             color.RGBA{0, 0, 0, 255},
}

// OpenBrowser opens default browser
func OpenBrowser(url string) {
    var err error

    switch runtime.GOOS {
    case "linux":
        err = exec.Command("xdg-open", url).Start()
    case "windows":
        err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    default:
        err = fmt.Errorf("unsupported platform")
    }
    if err != nil {
        log.Fatal(err)
    }

}

func main() {
    wnd := nucular.NewMasterWindowSize(nucular.WindowNoScrollbar, "DEV.TO api testing with GUI", image.Point{640, 346}, updatefn)
    wnd.SetStyle(style.FromTable(devtoThemeTable, 1.0))
    wnd.Main()
}

// updatefn is main window
func updatefn(w *nucular.Window) {
    w.Row(50).Dynamic(1)
    if w.ButtonText(fmt.Sprint("Get published articles")) {
        // Get 29 articles with sorting by score (1 day)
        response, _ = dev.GetPublishedArticles(dev.QueryArticle{
            Page:     0,
            Tag:      "",
            Username: "",
            State:    "",
            Top:      1,
        })
        DatePublishedTitle = "Published at:"
        comments = "comments:"
        rating = "rating:"

        // If API returns no description in article
        post.Title = wordwrap.WrapString(response[0].Title, 80)
        if response[0].Description == "" {
            post.Description = "No description given"
        } else {
            post.Description = wordwrap.WrapString(response[0].Description, 80)
        }

        ratingAmount   = strconv.Itoa(response[0].PositiveReactionsCount)
        commentsAmount = strconv.Itoa(response[0].CommentsCount)
        DatePublished  = response[0].PublishedTimestamp.Format(time.RFC822)
        post.URL       = response[0].URL
    }

    // Title
    w.Row(50).Dynamic(1)
    w.LabelColored(post.Title, "CC", color.RGBA{0, 118, 255, 255})

    // Date
    w.Row(1).Dynamic(2)
    w.Label(DatePublishedTitle, "RC")
    w.Label(DatePublished, "LC")

    // Description
    w.Row(100).Dynamic(1)
    w.Label(post.Description, "CC")

    // Rating
    w.Row(50).Dynamic(4)
    w.Label(comments, "RC")
    w.LabelColored(commentsAmount, "LC", color.RGBA{255, 20, 147, 255})

    // Comments
    w.Label(rating, "RC")
    w.LabelColored(ratingAmount, "LC", color.RGBA{255, 20, 147, 255})

    // URL
    w.Row(30).Dynamic(1)
    if w.ButtonText("Open this article in browser") {
        if post.URL != "" {
            OpenBrowser(post.URL)
        }
    }

    // Previous and Next buttons
    w.Row(30).Dynamic(2)
    if w.ButtonText("Previous") {
        if len(response) < 1 {
            post.Title = "Please click on button that says: \"Get published articles\"."
            post.Description = "Dude."
        } else {
            if postCounter == 0 || postCounter == 29 {
                postCounter = 29
            }
            postCounter--

            post.Title = wordwrap.WrapString(response[postCounter].Title, 80)
            if response[postCounter].Description == "" {
                post.Description = "No description given"
            } else {
                post.Description = wordwrap.WrapString(response[postCounter].Description, 80)
            }
            ratingAmount   = strconv.Itoa(response[postCounter].PositiveReactionsCount)
            commentsAmount = strconv.Itoa(response[postCounter].CommentsCount)
            DatePublished  = response[postCounter].PublishedTimestamp.Format(time.RFC822)
            post.URL       = response[postCounter].URL
        }
    }
    if w.ButtonText("Next") {
        if len(response) < 1 {
            post.Title = "Please click on button that says: \"Get published articles\", dude!"
            post.Description = "Dude!"
        } else {
            if postCounter == 0 || postCounter == 29 {
                postCounter = 0
            }

            postCounter++
            post.Title = wordwrap.WrapString(response[postCounter].Title, 80)
            if response[postCounter].Description == "" {
                post.Description = "No description given"
            } else {
                post.Description = wordwrap.WrapString(response[postCounter].Description, 80)
            }
            ratingAmount   = strconv.Itoa(response[postCounter].PositiveReactionsCount)
            commentsAmount = strconv.Itoa(response[postCounter].CommentsCount)
            DatePublished  = response[postCounter].PublishedTimestamp.Format(time.RFC822)
            post.URL       = response[postCounter].URL
        }
    }
}
