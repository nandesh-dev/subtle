package language

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type Tag struct {
	tag language.Tag
}

func (t Tag) String() string {
  return t.tag.String()
}

func Parse(rawTag string) (Tag, error) {
	tag, err := language.Parse(rawTag)
	if err != nil {
		return English, err
	}

	return Tag{tag: tag}, nil
}

var (
	Afrikaans            Tag = Tag{tag: language.Afrikaans}
	Amharic              Tag = Tag{tag: language.Amharic}
	Arabic               Tag = Tag{tag: language.Arabic}
	ModernStandardArabic Tag = Tag{tag: language.ModernStandardArabic}
	Azerbaijani          Tag = Tag{tag: language.Azerbaijani}
	Bulgarian            Tag = Tag{tag: language.Bulgarian}
	Bengali              Tag = Tag{tag: language.Bengali}
	Catalan              Tag = Tag{tag: language.Catalan}
	Czech                Tag = Tag{tag: language.Czech}
	Danish               Tag = Tag{tag: language.Danish}
	German               Tag = Tag{tag: language.German}
	Greek                Tag = Tag{tag: language.Greek}
	English              Tag = Tag{tag: language.English}
	AmericanEnglish      Tag = Tag{tag: language.AmericanEnglish}
	BritishEnglish       Tag = Tag{tag: language.BritishEnglish}
	Spanish              Tag = Tag{tag: language.Spanish}
	EuropeanSpanish      Tag = Tag{tag: language.EuropeanSpanish}
	LatinAmericanSpanish Tag = Tag{tag: language.LatinAmericanSpanish}
	Estonian             Tag = Tag{tag: language.Estonian}
	Persian              Tag = Tag{tag: language.Persian}
	Finnish              Tag = Tag{tag: language.Finnish}
	Filipino             Tag = Tag{tag: language.Filipino}
	French               Tag = Tag{tag: language.French}
	CanadianFrench       Tag = Tag{tag: language.CanadianFrench}
	Gujarati             Tag = Tag{tag: language.Gujarati}
	Hebrew               Tag = Tag{tag: language.Hebrew}
	Hindi                Tag = Tag{tag: language.Hindi}
	Croatian             Tag = Tag{tag: language.Croatian}
	Hungarian            Tag = Tag{tag: language.Hungarian}
	Armenian             Tag = Tag{tag: language.Armenian}
	Indonesian           Tag = Tag{tag: language.Indonesian}
	Icelandic            Tag = Tag{tag: language.Icelandic}
	Italian              Tag = Tag{tag: language.Italian}
	Japanese             Tag = Tag{tag: language.Japanese}
	Georgian             Tag = Tag{tag: language.Georgian}
	Kazakh               Tag = Tag{tag: language.Kazakh}
	Khmer                Tag = Tag{tag: language.Khmer}
	Kannada              Tag = Tag{tag: language.Kannada}
	Korean               Tag = Tag{tag: language.Korean}
	Kirghiz              Tag = Tag{tag: language.Kirghiz}
	Lao                  Tag = Tag{tag: language.Lao}
	Lithuanian           Tag = Tag{tag: language.Lithuanian}
	Latvian              Tag = Tag{tag: language.Latvian}
	Macedonian           Tag = Tag{tag: language.Macedonian}
	Malayalam            Tag = Tag{tag: language.Malayalam}
	Mongolian            Tag = Tag{tag: language.Mongolian}
	Marathi              Tag = Tag{tag: language.Marathi}
	Malay                Tag = Tag{tag: language.Malay}
	Burmese              Tag = Tag{tag: language.Burmese}
	Nepali               Tag = Tag{tag: language.Nepali}
	Dutch                Tag = Tag{tag: language.Dutch}
	Norwegian            Tag = Tag{tag: language.Norwegian}
	Punjabi              Tag = Tag{tag: language.Punjabi}
	Polish               Tag = Tag{tag: language.Polish}
	Portuguese           Tag = Tag{tag: language.Portuguese}
	BrazilianPortuguese  Tag = Tag{tag: language.BrazilianPortuguese}
	EuropeanPortuguese   Tag = Tag{tag: language.EuropeanPortuguese}
	Romanian             Tag = Tag{tag: language.Romanian}
	Russian              Tag = Tag{tag: language.Russian}
	Sinhala              Tag = Tag{tag: language.Sinhala}
	Slovak               Tag = Tag{tag: language.Slovak}
	Slovenian            Tag = Tag{tag: language.Slovenian}
	Albanian             Tag = Tag{tag: language.Albanian}
	Serbian              Tag = Tag{tag: language.Serbian}
	SerbianLatin         Tag = Tag{tag: language.SerbianLatin}
	Swedish              Tag = Tag{tag: language.Swedish}
	Swahili              Tag = Tag{tag: language.Swahili}
	Tamil                Tag = Tag{tag: language.Tamil}
	Telugu               Tag = Tag{tag: language.Telugu}
	Thai                 Tag = Tag{tag: language.Thai}
	Turkish              Tag = Tag{tag: language.Turkish}
	Ukrainian            Tag = Tag{tag: language.Ukrainian}
	Urdu                 Tag = Tag{tag: language.Urdu}
	Uzbek                Tag = Tag{tag: language.Uzbek}
	Vietnamese           Tag = Tag{tag: language.Vietnamese}
	Chinese              Tag = Tag{tag: language.Chinese}
	SimplifiedChinese    Tag = Tag{tag: language.SimplifiedChinese}
	TraditionalChinese   Tag = Tag{tag: language.TraditionalChinese}
	Zulu                 Tag = Tag{tag: language.Zulu}
)

func (t Tag) MarshalYAML() (interface{}, error) {
	return strings.ToUpper(t.tag.String()), nil
}

func (t *Tag) UnmarshalYAML(value *yaml.Node) error {
	var rawTag string
	if err := value.Decode(&rawTag); err != nil {
		return err
	}

	tag, err := language.Parse(strings.ToLower(rawTag))
	if err != nil {
		return err
	}

	t.tag = tag

	return nil
}

func (t Tag) Value() (driver.Value, error) {
	return strings.ToUpper(t.tag.String()), nil
}

func (t *Tag) Scan(src any) error {
	switch v := src.(type) {
	case string:
		tag, err := language.Parse(strings.ToLower(v))
		if err != nil {
			return err
		}

		t.tag = tag
	default:
		return fmt.Errorf("unsupposed type: %t", src)
	}

	return nil
}
