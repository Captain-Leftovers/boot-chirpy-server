package helpers

import "strings"

func CensorProfanity(s string) string {

	result := []string{}

	censor := "****"

	profanityWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	toCensor := strings.Split(s, " ")

CENSOR_LOOP:
	for _, to := range toCensor {
		for _, prof := range profanityWords {
			if strings.ToLower(to) == strings.ToLower(prof) {
				result = append(result, censor)
				continue CENSOR_LOOP
			}
		}
		result = append(result, to)
	}
	return strings.Join(result, " ")

}

// make loops named maybe
