package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func Slugify(title string) string {
	//convert title to lowercase
	title = norm.NFKD.String(title)

	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	//remove non alphanumeric/hypen char
	// re := regexp.MustCompile("[^a-z0-9-]+")
	// re := regexp.MustCompile("[^a-z0-9-]+")
	// s = re.ReplaceAllString(s, "")
	//removes duplicates hyphens
	redup := regexp.MustCompile("-+")
	s = redup.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "untitled"
	}
	return s
}

func Slugify2(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	re := regexp.MustCompile("[^a-z0-9-]+")
	s = re.ReplaceAllString(s, "")
	redup := regexp.MustCompile("-+")
	s = redup.ReplaceAllString(s, "-")
	return s
}

func SlugifyResearch(title string) string {
	s := norm.NFKD.String(title)
	s = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	s = re.ReplaceAllString(s, "")
	redup := regexp.MustCompile(`-+`)
	s = redup.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	if s == "" {
		s = "untitled"
	}

	return s

}
