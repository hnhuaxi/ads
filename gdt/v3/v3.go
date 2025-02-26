package v3

import (
	"context"
	"strconv"

	"github.com/antihax/optional"
	"github.com/hnhuaxi/ads"
	"github.com/hysios/x/utils/ptr"
	adsv3 "github.com/tencentad/marketing-api-go-sdk/pkg/ads/v3"
	apiv3 "github.com/tencentad/marketing-api-go-sdk/pkg/api/v3"
	config "github.com/tencentad/marketing-api-go-sdk/pkg/config/v3"
)

// adsv3 "github.com/tencentad/marketing-api-go-sdk/pkg/ads/v3"

type GdtV3API struct {
	AccountID int64
	*adsv3.SDKClient
}

func NewGdtAPI(accountId string, accessToken string, debug bool) *GdtV3API {
	id, _ := strconv.ParseInt(accountId, 10, 64)
	return &GdtV3API{
		AccountID: id,
		SDKClient: adsv3.Init(&config.SDKConfig{
			AccessToken: accessToken,
			IsDebug:     debug,
		}),
	}
}

var AdvertisersFields = []string{
	"dynamic_creative_id",
	"dynamic_creative_name",
	"dynamic_creative_type",
	"adgroup_id",
	"delivery_mode",
	"configured_status",
	"adgroup.campaign_type",
	"created_time",
	"last_modified_time",
	"creative_components",
	"is_deleted",
}

// Adcreatives
func (g *GdtV3API) Adcreatives(page, pageSize int) (objs []ads.Map, total int64, err error) {
	var ctx = context.TODO()

	resp, _, err := g.SDKClient.DynamicCreatives().Get(ctx, g.AccountID, &apiv3.DynamicCreativesGetOpts{
		Page:     optional.NewInt64(int64(page)),
		PageSize: optional.NewInt64(int64(pageSize)),
		Fields:   optional.NewInterface(AdvertisersFields),
	})

	if err != nil {
		return nil, 0, err
	}

	if resp.List == nil {
		return nil, 0, nil
	}

	if resp.PageInfo == nil {
		return nil, 0, nil
	}

	objs, err = ads.ToMapSlice(*resp.List)
	if err != nil {
		return nil, 0, err
	}

	return objs, ptr.Type(resp.PageInfo.TotalNumber), nil
}

// AllAdcreatives
func (g *GdtV3API) AllAdcreatives() (objs []ads.Map, err error) {
	for page := 1; ; page++ {
		resp, total, err := g.Adcreatives(page, 100)
		if err != nil {
			objs = append(objs, resp...)
			return objs, nil
		}

		objs = append(objs, resp...)
		if int64(len(objs)) >= total {
			break
		}
	}

	return
}

func (g *GdtV3API) Pages(pageType string, page, pageSize int) (objs []ads.Map, total int64, err error) {
	var ctx = context.TODO()

	resp, _, err := g.SDKClient.Pages().Get(ctx, g.AccountID, &apiv3.PagesGetOpts{
		Filtering: optional.NewInterface([]interface{}{
			map[string]interface{}{
				"field":    "page_type",
				"operator": "EQUALS",
				"value": []string{
					pageType,
				},
			},
		}),
		Page:     optional.NewInt64(int64(page)),
		PageSize: optional.NewInt64(int64(pageSize)),
	})

	if err != nil {
		return nil, 0, err
	}

	if resp.List == nil {
		return nil, 0, nil
	}

	if resp.PageInfo == nil {
		return nil, 0, nil
	}

	objs, err = ads.ToMapSlice(*resp.List)
	if err != nil {
		return nil, 0, err
	}

	return objs, ptr.Type(resp.PageInfo.TotalNumber), nil
}

var WechatPagesFields = []string{
	"page_id",
	"page_name",
	"adcreative_template_id",
	"marketing_goal",
	"marketing_sub_goal",
	"marketing_target_type",
	"marketing_carrier_type",
	"page_type",
	"marketing_carrier_id",
	"canvas_type",
	"page_status",
	"site_set",
	"marketing_scene",
	"source_type",
	"live_video_mode",
	"live_video_sub_mode",
	"live_notice_id",
	"product_catalog_id",
	"product_source",
	"raw_adcreative_template_id",
	"buying_type",
	"product_mode",
}

// WechatPages
func (g *GdtV3API) WechatPages(page, pageSize int) (objs []ads.Map, total int64, err error) {
	var ctx = context.TODO()

	resp, _, err := g.SDKClient.WechatPages().Get(ctx, g.AccountID, &apiv3.WechatPagesGetOpts{
		Page:     optional.NewInt64(int64(page)),
		PageSize: optional.NewInt64(int64(pageSize)),
		Fields:   optional.NewInterface(WechatPagesFields),
	})

	if err != nil {
		return nil, 0, err
	}

	if resp.List == nil {
		return nil, 0, nil
	}

	if resp.PageInfo == nil {
		return nil, 0, nil
	}

	objs, err = ads.ToMapSlice(*resp.List)
	if err != nil {
		return nil, 0, err
	}

	return objs, ptr.Type(resp.PageInfo.TotalNumber), nil
}

var XJPagesFields = []string{}

// XJPages
func (g *GdtV3API) XJPages(pageType string, page, pageSize int) (objs []ads.Map, total int64, err error) {
	var ctx = context.TODO()

	resp, _, err := g.SDKClient.XijingPageList().Get(ctx, g.AccountID, &apiv3.XijingPageListGetOpts{
		PageIndex: optional.NewInt64(int64(page)),
		PageSize:  optional.NewInt64(int64(pageSize)),
		PageType:  optional.NewInterface(pageType),
	})

	if err != nil {
		return nil, 0, err
	}

	if resp.List == nil {
		return nil, 0, nil
	}

	if resp.PageInfo == nil {
		return nil, 0, nil
	}

	objs, err = ads.ToMapSlice(*resp.List)
	if err != nil {
		return nil, 0, err
	}

	return objs, ptr.Type(resp.PageInfo.TotalNumber), nil
}

// AllXJPages
func (g *GdtV3API) AllXJPages(pageType string) (objs []ads.Map, err error) {
	for page := 1; ; page++ {
		resp, total, err := g.XJPages(pageType, page, 100)
		if err != nil {
			objs = append(objs, resp...)
			return objs, nil
		}
		objs = append(objs, resp...)
		if int64(len(objs)) >= total {
			break
		}
	}

	return
}

var XJPages_TYPES = []string{
	"XJ_DEFAULT_H5",
	"XJ_ANDROID_APP_H5",
	"XJ_IOS_APP_H5",
	"XJ_WEBSITE_H5",
	"XJ_ANDROID_APP_NATIVE",
	"XJ_IOS_APP_NATIVE",
	"XJ_WEBSITE_NATIVE",
	"XJ_FENGLING_LBS",
}

// AllPages
func (g *GdtV3API) AllPages(pageType string, subTypes ...string) (objs []ads.Map, err error) {
	switch pageType {
	case "XJ_PAGES":
		allPages := make([]ads.Map, 0)
		if len(subTypes) == 0 {
			subTypes = XJPages_TYPES
		}

		for _, subType := range subTypes {
			pages, err := g.AllXJPages(subType)
			if err != nil {
				allPages = append(allPages, pages...)
				continue
				// return nil, err
			}
			allPages = append(allPages, pages...)
		}
		return allPages, nil
	case "WECHAT_PAGES":
		return g.AllWechatPages()
	default:
		for page := 1; ; page++ {
			resp, total, err := g.Pages(pageType, page, 100)
			if err != nil {
				objs = append(objs, resp...)
				return objs, nil
			}
			objs = append(objs, resp...)
			if int64(len(objs)) >= total {
				break
			}
		}
	}

	return
}

// AllWechatPages
func (g *GdtV3API) AllWechatPages() (objs []ads.Map, err error) {
	for page := 1; ; page++ {
		resp, total, err := g.WechatPages(page, 100)
		if err != nil {
			objs = append(objs, resp...)
			return objs, nil
		}
		objs = append(objs, resp...)
		if int64(len(objs)) >= total {
			break
		}
	}

	return
}

var VideoFields = []string{
	"media_signature",
	"media_id",
	"media_width",
	"media_height",
	"created_time",
	"last_modified_time",
	"source_type",
	"product_catalog_id",
	"product_outer_id",
	"owner_account_id",
	"status",
	"media_description",
	"sample_aspect_ratio",
}

// Videos
func (g *GdtV3API) Videos(page, pageSize int, ids ...string) (objs []ads.Map, total int64, err error) {
	var (
		ctx  = context.TODO()
		opts = &apiv3.VideosGetOpts{
			AccountId: optional.NewInt64(g.AccountID),
			Page:      optional.NewInt64(int64(page)),
			PageSize:  optional.NewInt64(int64(pageSize)),
			Fields:    optional.NewInterface(VideoFields),
			Filtering: filterIds("media_id", ids),
		}
	)

	resp, _, err := g.SDKClient.Videos().Get(ctx, opts)

	if err != nil {
		return nil, 0, err
	}

	if resp.List == nil {
		return nil, 0, nil
	}

	if resp.PageInfo == nil {
		return nil, 0, nil
	}

	objs, err = ads.ToMapSlice(*resp.List)
	if err != nil {
		return nil, 0, err
	}

	return objs, ptr.Type(resp.PageInfo.TotalNumber), nil
}

// AllVideosLoop
func (g *GdtV3API) AllVideosLoop(ids ...string) (objs []ads.Map, err error) {
	for page := 1; ; page++ {
		resp, total, err := g.Videos(page, 100, ids...)
		if err != nil {
			objs = append(objs, resp...)
			return objs, nil
		}

		objs = append(objs, resp...)
		if int64(len(objs)) >= total {
			break
		}
	}

	return
}

// AllVideos
func (g *GdtV3API) AllVideos(ids ...string) (all []ads.Map, err error) {
	if len(ids) == 0 {
		return g.AllVideosLoop()
	}

	for len(ids) > 0 {
		l := min(100, len(ids))
		objs, _ := g.AllVideosLoop(ids[:l]...)
		all = append(all, objs...)
		ids = ids[l:]
	}
	return
}

var ImageFields = []string{
	"image_signature",
	"image_id",
	"image_width",
	"image_height",
	"created_time",
	"last_modified_time",
	"source_type",
	"product_catalog_id",
	"product_outer_id",
	"owner_account_id",
	"status",
	"image_description",
	"sample_aspect_ratio",
}

// Images
func (g *GdtV3API) Images(page, pageSize int, ids ...string) (objs []ads.Map, total int64, err error) {
	var (
		ctx  = context.TODO()
		opts = &apiv3.ImagesGetOpts{
			AccountId: optional.NewInt64(g.AccountID),
			Page:      optional.NewInt64(int64(page)),
			PageSize:  optional.NewInt64(int64(pageSize)),
			Fields:    optional.NewInterface(ImageFields),
		}
	)

	resp, _, err := g.SDKClient.Images().Get(ctx, opts)

	if err != nil {
		return nil, 0, err
	}

	if resp.List == nil {
		return nil, 0, nil
	}

	if resp.PageInfo == nil {
		return nil, 0, nil
	}

	objs, err = ads.ToMapSlice(*resp.List)
	if err != nil {
		return nil, 0, err
	}

	return objs, ptr.Type(resp.PageInfo.TotalNumber), nil
}

// ImagesLoop
func (g *GdtV3API) ImagesLoop(ids ...string) (objs []ads.Map, err error) {
	for page := 1; ; page++ {
		resp, total, err := g.Images(page, 100, ids...)
		if err != nil {
			objs = append(objs, resp...)
			return objs, nil
		}

		objs = append(objs, resp...)
		if int64(len(objs)) >= total {
			break
		}
	}

	return
}

// AllImages
func (g *GdtV3API) AllImages(ids ...string) (all []ads.Map, err error) {
	if len(ids) == 0 {
		return g.ImagesLoop()
	}

	for len(ids) > 0 {
		l := min(100, len(ids))
		objs, _ := g.ImagesLoop(ids[:l]...)
		all = append(all, objs...)
		ids = ids[l:]
	}
	return
}

func filterIds(key string, ids []string) optional.Interface {
	if len(ids) == 0 {
		return optional.EmptyInterface()
	}
	return optional.NewInterface([]interface{}{
		map[string]interface{}{
			"field":    key,
			"operator": "IN",
			"values":   ids,
		},
	})
}
