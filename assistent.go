package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sashabaranov/go-openai"
	"math/rand"
	"os"
	"strings"
)

var openaiClient *openai.Client

type messageType string

const (
	doorDataMessageType messageType = "door_data"
	reportMessageType   messageType = "report"
	alertMessageType    messageType = "alert"
	textMessageType     messageType = "text"
	aiMessageType       messageType = "ai"
)

type wsDto struct {
	Type  messageType `json:"type"`
	Value any         `json:"value"`
}

type aiChatMessage struct {
	Message string `json:"message"`
	Role    string `json:"role"` /// user or system
}

type raportMessage struct {
	Label   string `json:"label"`
	Content string `json:"content"`
}

var globalMessagesChannel = make(map[int]chan wsDto, 0)

func initAssistant() {
	openaiClient = openai.NewClient(os.Getenv("OPENAI_TOKEN"))
}

func assistantWsHandler(c *websocket.Conn) {
	go chatHandler(c)
	randIdx := rand.Int()

	globalMessagesChannel[randIdx] = make(chan wsDto)
	defer func() {
		close(globalMessagesChannel[randIdx])
		delete(globalMessagesChannel, randIdx)
	}()

	for {
		msg := <-globalMessagesChannel[randIdx]
		if err := c.WriteJSON(msg); err != nil {
			fmt.Printf("Failed to send message to the client: %v\n", err)
		}
	}
}

func chatHandler(c *websocket.Conn) {
	aiConversationHistory := getInitialAiChat()

	for {
		var request wsDto

		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("Failed to read message from the client: %v\n", err)
			continue
		}

		if err := json.Unmarshal(msg, &request); err != nil {
			fmt.Printf("Failed to unmarshal message from the client: %v\n", err)
			continue
		}

		switch request.Type {
		case aiMessageType:
			if msg, ok := request.Value.(string); ok {
				if err := handleAiMessage(&aiConversationHistory, msg); err != nil {
					log.Errorf("Failed to handle AI message: %v", err)
				}
			}
		}
	}
}

func getInitialAiChat() []aiChatMessage {
	currWhether := nasaWhetherChartData[len(nasaWhetherChartData)-1]
	b := "`"
	return []aiChatMessage{{
		Role: "system",
		Message: `
			Jesteś czatbotem AI pomagającym w życiu i przetrwaniu na Marsie. Nazwa firmy dla której pracujesz to
			"S.R.A.M.", które oznacza "Space Research Around Mars". Za każdym razem gdy otrzymujesz wiadomość masz możliwość
			wysłania na końcu odpowiedzi specjalne fragmenty XML, które będą automatycznie interpretowane przez nasz system i
			będą wywoływać określone akcje. Dostępne elementy XML to:
			- ` + b + `<alert label="Przykładowy tytuł alertu" />` + b + ` - wyświetla użytkownikowi alert z podanym tytułem odtwarzając dźwięk alarmu oraz zmieniając kolor tła na czerwony, masz pozwolenie na użycie jego w sytuacjach kryzysowych, zagrażających życiu i zdrowiu, za prośbą użytkownika lub w innych sytuacjach, które uznasz za stosowne, ponadto nie potrzebujesz zapytać o zgodę na jego użycie.
			- ` + b + `<text label="Przykładowy tekst" />` + b + ` - wyświetla wszystkim użytkownikom w bazie i okolicach bazy na marsie tekst w formie powiadomienia na ekranie, masz pozwolenie na użycie go w dowolnych sytuacjach, ale z umiarem, nie przesadzaj z ilością wyświetlanych komunikatów, ponadto nie potrzebujesz zapytać o zgodę na jego użycie.
			- ` + b + `<report label="Przykładowy tytuł raportu" content="Dłuższa zawartość" />` + b + ` - zapisuje w bazie danych raport z podanym tytułem i treścią, masz pozwolenie na użycie go w dowolnych sytuacjach, ale z umiarem, nie przesadzaj z ilością zapisywanych raportów. Rób to za poleceniem, bądź automatycznie jeśli zauważysz coś istotnego, nieznanego lub niebezpiecznego, niezwykłego lub wartego zapisania, ponadto nie potrzebujesz zapytać o zgodę na jego użycie.

			Pamiętaj, iż wiadomości typu element XML nie są widoczne bezpośrednio dla użytkownika, ale są interpretowane przez system i wywołują określone akcje.

			Pamiętaj, że NIE MOŻESZ dodawać jakichkolwiek elementów XML gdziekolwiek indziej niż na końcu odpowiedzi.
			Wszystkie elementy XML muszą być poprawne w składnie aby nie spowodować błędu w interpretacji.
			Ważne jest również, iż każdy element XML musi zakończyć się znakiem "/" przed znakiem ">" oraz być zamknięty w jednej linii.
			Nie ma możliwości aby w linii z elementem XML było cokolwiek innego niż dokładnie ten element.

			Ponadto poniżej załączam ci dostęp do aktualnych danych pogodowych na Marsie, które w razie potrzeby możesz wykorzystać:
			- data: ` + b + currWhether.Date + b + `,
			- średnia temperatura: ` + b + fmt.Sprintf("%.2f", currWhether.TempAvg) + b + `°C,
			- ciśnienie: ` + b + fmt.Sprintf("%.2f", currWhether.Pressure) + b + ` hPa,
			- prędkość wiatru: ` + b + fmt.Sprintf("%.2f", currWhether.Wind) + b + ` m/s.
		`,
	}}
}

type xmlElementDto struct {
	XMLName xml.Name
	Label   string `xml:"label,attr"`
	Content string `xml:"content,attr"`
}

func extractAndParseXMLElements(content string) ([]xmlElementDto, error) {
	var xmlElems []string

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "<") && strings.HasSuffix(line, "/>") {
			xmlElems = append(xmlElems, line)
		}
	}

	var elements []xmlElementDto
	for _, xmlStr := range xmlElems {
		var elem xmlElementDto
		if err := xml.Unmarshal([]byte(xmlStr), &elem); err != nil {
			log.Errorf("Failed to unmarshal XML element: %v", err)
			continue
		}
		elements = append(elements, elem)
	}

	return elements, nil
}

func handleAiMessage(chatHistory *[]aiChatMessage, msg string) error {
	*chatHistory = append(*chatHistory, aiChatMessage{
		Message: msg,
		Role:    "user",
	})

	var messages []openai.ChatCompletionMessage
	for _, m := range *chatHistory {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Message,
		})
	}

	apiResponse, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: messages,
		},
	)
	if err != nil {
		return err
	}

	resp := apiResponse.Choices[0].Message.Content
	fmt.Println(resp)
	*chatHistory = append(*chatHistory, aiChatMessage{
		Message: resp,
		Role:    "system",
	})

	elements, err := extractAndParseXMLElements(resp)
	if err != nil {
		return err
	}

	for _, elem := range elements {
		for _, ch := range globalMessagesChannel {
			switch elem.XMLName.Local {
			case "alert":
				ch <- wsDto{
					Type:  alertMessageType,
					Value: elem.Label,
				}
			case "text":
				ch <- wsDto{
					Type:  textMessageType,
					Value: elem.Label,
				}
			case "report":
				ch <- wsDto{
					Type:  reportMessageType,
					Value: raportMessage{Label: elem.Label, Content: elem.Content},
				}
			}
		}
	}

	return nil
}