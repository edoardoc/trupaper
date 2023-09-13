package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"github.com/joho/godotenv"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type PaperLine struct {
	Position1 fyne.Position
	Position2 fyne.Position
}
type dbLine struct {
	documentId string
	lines      []PaperLine
}

// paper is a widget that can be drawn on, it is a container to detect mouse events
type paper struct {
	widget.BaseWidget
	mainContainer *fyne.Container // this holds everything

	// TODO: these two containers can now become one since there is the dblines array
	committedLinesContainer *fyne.Container // this holds the drawing that is in the db
	linesContainer          *fyne.Container // this holds whatever we are drawing at the moment

	current   []PaperLine
	dblines   []dbLine
	deltaDraw int // this is 0 when the paper is showing the current drawing, negative values to look into the "past"

	isDrawing  bool
	lastPos    fyne.Position
	myRWapiKey string
}

func NewPaper() *paper {
	p := &paper{}
	p.myRWapiKey = os.Getenv("API_KEY")
	p.ExtendBaseWidget(p)
	p.committedLinesContainer = container.NewWithoutLayout()
	p.linesContainer = container.NewWithoutLayout() // linesContainer is empty upon start

	err := p.loadAllLinesJSON() // load the committedLines from the remote db
	if err != nil {
		log.Printf("loading the committedLines from the remote db error %v ", err)
	}

	p.mainContainer = container.NewWithoutLayout()
	p.mainContainer.Add(p.committedLinesContainer)
	p.mainContainer.Add(p.linesContainer)

	return p
}

// this repaints the dblines up to the end + deltaDraw
func (p *paper) repaintCommittedLines() {
	p.committedLinesContainer.RemoveAll()
	for i := 0; i < len(p.dblines)+p.deltaDraw; i++ {
		p.drawCommittedData(p.dblines[i].lines)
	}
	p.committedLinesContainer.Refresh()
}

func (p *paper) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.mainContainer)
}

func (p *paper) commitCurrentLines() {
	if len(p.current) == 0 {
		return
	}
	err := p.sendLinesJSON()
	if err != nil {
		log.Println("error = ", err)
		return
	}
	//  then add last documentId lines to committedLinesContainer
	p.drawCommittedData(p.dblines[len(p.dblines)-1].lines)

	// I need to empty the linesContainer
	p.linesContainer.RemoveAll()
	p.current = []PaperLine{}
	p.mainContainer.Refresh()
}

func (p *paper) loadAllLinesJSON() error {
	page := 1
	perPage := 100
	var responseBody []byte
	var responseObject struct {
		Page      int `json:"page"`
		PerPage   int `json:"perPage"`
		Revisions []struct {
			Document struct {
				ID        string `json:"_id"`
				VaultMeta struct {
					Creator string `json:"creator"`
					Ts      int    `json:"ts"`
				} `json:"_vault_md"`
				Polyline []PaperLine `json:"polyline"`
			} `json:"document"`
			Revision      string `json:"revision"`
			TransactionID string `json:"transactionId"`
		} `json:"revisions"`
	}

	for {
		requestBody, err := json.Marshal(map[string]int{
			"page":    page,
			"perPage": perPage,
		})
		if err != nil {
			return err
		}

		request, err := http.NewRequest("POST", "https://vault.immudb.io/ics/api/v1/ledger/default/collection/default/documents/search", bytes.NewBuffer(requestBody))
		if err != nil {
			return err
		}

		request.Header.Set("accept", "application/json")
		request.Header.Set("X-API-Key", p.myRWapiKey)
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(responseBody, &responseObject)
		if err != nil {
			return err
		}

		fmt.Println("fetching page ", page, " Revisions: ", len(responseObject.Revisions)) // 16
		if len(responseObject.Revisions) == 0 {
			break
		}

		for _, r := range responseObject.Revisions {
			p.dblines = append(p.dblines, dbLine{documentId: r.Document.ID, lines: r.Document.Polyline})
		}
		page++
	}
	p.repaintCommittedLines()

	return nil
}

func (p *paper) Left() {
	p.deltaDraw--
	if p.deltaDraw < -len(p.dblines) {
		p.deltaDraw = -len(p.dblines)
	}
	fmt.Println("going left, deltaDraw is ", p.deltaDraw)
	p.repaintCommittedLines()
}

func (p *paper) Right() {
	p.deltaDraw++
	if p.deltaDraw > 0 {
		p.deltaDraw = 0
	}
	p.repaintCommittedLines()
	fmt.Println("going right, deltaDraw is ", p.deltaDraw)
}

func (p *paper) drawCommittedData(lines []PaperLine) {
	fmt.Printf("drawing %d lines\n", len(lines))
	for _, l := range lines {
		line := &canvas.Line{
			StrokeColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			StrokeWidth: 5,
			Position1:   l.Position1,
			Position2:   l.Position2,
		}

		p.committedLinesContainer.Add(line)
	}
}

// sendLinesJSON sends the current lines to the remote db
// and stores the current lines in the dblines structure
func (p *paper) sendLinesJSON() error {
	url := "https://vault.immudb.io/ics/api/v1/ledger/default/collection/default/document"

	jsonBytes, err := json.Marshal(map[string]interface{}{
		"polyline": p.current,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-API-Key", p.myRWapiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// handle response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	documentId := data["documentId"].(string)
	p.dblines = append(p.dblines, dbLine{documentId: documentId, lines: p.current})

	return nil
}

func (p *paper) MouseUp(w *desktop.MouseEvent) {
	p.isDrawing = false
	p.lastPos = fyne.Position{}
	p.commitCurrentLines()
}

func (p *paper) MouseDown(w *desktop.MouseEvent) {
	p.isDrawing = true
	p.deltaDraw = 0
	fmt.Println("Reset, deltaDraw is ", p.deltaDraw)
	p.repaintCommittedLines()
}

func (p *paper) MouseIn(_ *desktop.MouseEvent) {}
func (p *paper) MouseOut()                     {}
func (p *paper) MouseMoved(e *desktop.MouseEvent) {
	if p.isDrawing {
		if p.lastPos != (fyne.Position{}) {
			line := canvas.NewLine(color.Black)
			line.StrokeWidth = 5
			line.Position1 = p.lastPos
			line.Position2 = e.Position

			p.current = append(p.current, PaperLine{Position1: p.lastPos, Position2: e.Position})
			p.linesContainer.Add(line)
			p.linesContainer.Refresh()
		}
		p.lastPos = e.Position
	} else {
		p.lastPos = fyne.Position{}
	}
}
