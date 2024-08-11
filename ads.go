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
	PTText
)

type Asset struct {
	AccountID   string
	AccountName string
	AssetID     string
	Name        string
	PageType    PageType
	SubType     string
	Texts       []string
	SubAssets   []*SubAsset
	Signature   string
	Version     string
}

type SubAssetType int

const (
	SATUnknown SubAssetType = iota
	SATImage
	SATVideo
	SATPageUrl
)

type SubAsset struct {
	Type SubAssetType
	Url  string
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
