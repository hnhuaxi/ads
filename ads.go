package ads

import (
	"fmt"

	"github.com/hysios/x/providers"
)

//go:generate stringer -type=PageType
type PageType int

const (
	PTUnknown PageType = iota
	PTPageUrl
	PTVideo
	PTImage
)

type Asset struct {
	AccountID   string
	AccountName string
	AssetID     string
	Name        string
	PageType    PageType
	SubType     string
	Url         string
	Signature   string
}

type AdcreativeMatchFunc func(adcreative Map) (process bool)

type GetAdcreatives interface {
	Assets() ([]*Asset, error)
	SetAdcreativesFunc(match AdcreativeMatchFunc)
	OnlyAdcreatives(on bool)
}

func Open(provider string, accountId string, accessToken string, debug bool) (GetAdcreatives, error) {
	ctor, ok := advProviders.Lookup(provider)
	if !ok {
		return nil, fmt.Errorf("provider %s not found", provider)
	}

	return ctor(accountId, accessToken, debug)
}

var advProviders providers.Provider[string, func(string, string, bool) (GetAdcreatives, error)]

func RegisterProvider(name string, f func(string, string, bool) (GetAdcreatives, error)) {
	advProviders.Register(name, f)
}
