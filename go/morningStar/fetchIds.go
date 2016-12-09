package morningStar

import (
	"log"
	"regexp"
)

const urlIds = `https://elasticsearch.vibioh.fr/funds/morningStarId/_search?size=8000`

var idRegex = regexp.MustCompile(`"_id":"(.*?)"`)

func fetchIds() [][]byte {
	body, err := getBody(urlIds)
	if err != nil {
		log.Print(err)
		return nil
	}

	idsMatch := idRegex.FindAllSubmatch(body, -1)

	ids := make([][]byte, 0, len(idsMatch))
	for _, match := range idsMatch {
		ids = append(ids, match[1])
	}

	return ids
}
