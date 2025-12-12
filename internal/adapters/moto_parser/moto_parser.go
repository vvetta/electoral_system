package motoparser

import (
	"fmt"
	"log"
	"unicode"
	"net/http"
	"strconv"
	"strings"

	"github.com/vvetta/electoral_system/internal/domain"
	"github.com/vvetta/electoral_system/internal/usecase"
	
	"golang.org/x/net/html"
)

type motoParser struct {
	url string
	prodCardClassName string
	maxPageCount int
}

func NewMotoParser(
	url string,
	prodCardClassName string,
	maxPageCount int,
) usecase.MotoParser {
	return &motoParser{
		url: url,
		prodCardClassName: prodCardClassName,
		maxPageCount: maxPageCount,	
	}
}

/*
Парсер настроен на получение мотоциклов с сайта: mr-moto.ru
На сайте нет никакой защиты от ботов и простой парсинг html

Пагинация на сайте реализована простым увеличением страницы
через кнопку "Показать еще". Для того чтобы спарсить все
мотоциклы было решено использвать offset.
*/
func (p *motoParser) GetAllMoto() ([]domain.Moto, error) {
	log.Print("MotoParser: Start!")	
	var motos []domain.Moto

	var offset int = 0
	for i := 1; i <= p.maxPageCount; i++ {
		url := p.url + "?nav-catalog=page-" + strconv.Itoa(i)
		log.Print("MotoParser: parsing page: ", url)

		motosFromPage, err := p.getMotosFromPage(url, offset); 
		if err != nil {
			log.Print("MotoParser: get motos from page error: ", err)
			return nil, fmt.Errorf("%w: Ошибка получения мотоциклов со страницы: %s", domain.ParseMotoError, url)
		}

		if len(motosFromPage) == 0 {
			log.Print("MotoParser: find 0 moto from page. Stop parsing...")
			break
		}
		log.Printf("MotoParser: find %d moto from page", len(motosFromPage))

		motos = append(motos, motosFromPage...)
		offset += len(motosFromPage)
		log.Print("MotoParser: motoFromPage succsessfully appended to motos")	
	}

	log.Print("MotoParser: End!")
	return motos, nil
}

func (p *motoParser) getMotosFromPage(url string, offset int) ([]domain.Moto, error) {
	log.Print("MotoParser-getMotosFromPage: Start!")	

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("MotoParser-getMotosFromPage: create request error: ", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print("MotoParser-getMotosFromPage: client.do error: ", err)
		return nil, err	
	}

	if resp.StatusCode >= 400 {
		log.Print("MotoParser-getMotosFromPage: resp status code >= 400")
		return nil, domain.ParseMotoError
	}

	defer resp.Body.Close()

	pageNode, err := html.Parse(resp.Body)
	if err != nil {
		log.Print("MotoParser-getMotosFromPage: html.parse error: ", err)
		return nil, err	
	}

	motosFromPage, err := p.getMotosFromHtml(pageNode, offset)
	if err != nil {
		log.Print("MotoParser-getMotosFromPage: get motos from html error: ", err)
		return nil, err
	}

	log.Print("MotoParser-getMotosFromPage: End!")
	return motosFromPage, nil
}

func (p *motoParser) getMotosFromHtml(pageNode *html.Node, offset int) ([]domain.Moto, error) {
	log.Print("MotoParser-getMotosFromHtml: Start!")
	
	var cards []*html.Node

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n == nil {
			return
		}

		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					classes := strings.Fields(attr.Val)
					for _, value := range classes {
						if value == p.prodCardClassName {
							cards = append(cards, n)
							break
						}
					}
					break
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
    	walk(child)
    }
	}

	walk(pageNode)

	if len(cards) == 0 {
		log.Print("MotoParser-getMotosFromHtml: no cards found")
		return nil, nil // или domain.ParseMotoError — по твоему дизайну
	}

	if offset >= len(cards) {
		log.Print("MotoParser-getMotosFromHtml: offset >= len(cards), no new motos")
		return nil, nil
	}

	start := offset
	var motos []domain.Moto
	
	for i := start; i < len(cards); i++ {
		card := cards[i]

		moto, err := parseMotoCard(card)
		if err != nil {
			log.Printf("MotoParser-getMotosFromHtml: parse card #%d error: %v", i, err)
			continue
		}

		motos = append(motos, moto)
	}

	log.Printf("MotoParser-getMotosFromHtml: parsed %d motos", len(motos))
	log.Print("MotoParser-getMotosFromHtml: End!")
	return motos, nil
}

func parseMotoCard(card *html.Node) (domain.Moto, error) {
	var m domain.Moto

	if card == nil {
		return m, fmt.Errorf("card is nil")
	}

	// 1. Находим div.slider-card__title
  titleNodes := getNodesByClass(card, "slider-card__title")
  if len(titleNodes) == 0 {
      return m, fmt.Errorf("slider-card__title not found")
  }

    // 2. Достаём текст (он лежит внутри <a>…, но можно просто собрать весь текст внутри div)
  motoTitle := strings.TrimSpace(getText(titleNodes[0]))
  m.Name = motoTitle

	// Блок характеристик
	infoBlocks := getNodesByClass(card, "slider-card__info")
	if len(infoBlocks) == 0 {
		return m, fmt.Errorf("slider-card__info not found")
	}
	info := infoBlocks[0]

	rows := getNodesByClass(info, "slider-card__row")
	for _, row := range rows {
		nameNodes := getNodesByClass(row, "slider-card__info-name")
		valueNodes := getNodesByClass(row, "slider-card__info-text")
		if len(nameNodes) == 0 || len(valueNodes) == 0 {
			continue
		}

		name := getText(nameNodes[0])
		value := getText(valueNodes[0])

		switch name {
		case "Год":
			if v, err := parseIntFromString(value); err == nil {
				m.Year = v
			}
		case "Пробег ТС":
			if v, err := parseIntFromString(value); err == nil {
				m.Mileage = v
			}
		case "Объем Д":
			if v, err := parseIntFromString(value); err == nil {
				m.EngineSize = v
			}
		case "Класс мототехники":
			m.MotoType = value
		case "Мотосалон":
			m.Location = value
		}
	}

	// Цена
	priceNodes := getNodesByClass(card, "slider-card__price-title")
	if len(priceNodes) > 0 {
		priceText := getText(priceNodes[0]) // типа "1 190 000 р."
		if v, err := parseIntFromString(priceText); err == nil {
			m.Price = int64(v)
		}
	}

	return m, nil
}

func getNodesByClass(root *html.Node, className string) []*html.Node {
	var result []*html.Node

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n == nil {
			return
		}

		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					classes := strings.Fields(attr.Val)
					for _, c := range classes {
						if c == className {
							result = append(result, n)
							break
						}
					}
					break
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(root)
	return result
}

func parseIntFromString(s string) (int, error) {
	var digits []rune
	for _, r := range s {
		if unicode.IsDigit(r) {
			digits = append(digits, r)
		}
	}
	if len(digits) == 0 {
		return 0, fmt.Errorf("no digits in %q", s)
	}

	return strconv.Atoi(string(digits))
}

func getText(n *html.Node) string {
	if n == nil {
		return ""
	}

	var b strings.Builder

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			b.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(n)
	return strings.TrimSpace(b.String())
}
