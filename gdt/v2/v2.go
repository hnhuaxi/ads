package v2

import (
	"context"
	"strconv"

	"github.com/antihax/optional"
	"github.com/hnhuaxi/ads"
	"github.com/hysios/x/utils/ptr"
	gdtads "github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	"github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"go.uber.org/zap"
)

type GdtAPI struct {
	AccountID int64
	*gdtads.SDKClient
	log *zap.SugaredLogger
}

func NewGdtAPI(accountId string, accessToken string, debug bool) *GdtAPI {
	id, _ := strconv.ParseInt(accountId, 10, 64)
	tads := gdtads.Init(&config.SDKConfig{
		AccessToken: accessToken,
		IsDebug:     debug,
	})
	tads.UseProduction()

	return &GdtAPI{
		AccountID: id,
		SDKClient: tads,
		log:       zap.S(),
	}
}

var AdcreativesFields = []string{
	"adcreative_id",
	"adcreative_name",
	"campaign_id",
	"promoted_object_type",
	"adcreative_template_id",
	"adcreative_elements",
	"created_time",
	"last_modified_time",
	"page_type",
	"page_spec",
}

// Adcreatives
func (g *GdtAPI) Adcreatives(page, pageSize int) (objs []ads.Map, total int64, err error) {
	var ctx = context.TODO()

	resp, _, err := g.SDKClient.Adcreatives().Get(ctx, g.AccountID, &api.AdcreativesGetOpts{
		Page:     optional.NewInt64(int64(page)),
		PageSize: optional.NewInt64(int64(pageSize)),
		Fields:   optional.NewInterface(AdcreativesFields),
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

// Pages
func (g *GdtAPI) Pages(page, pageSize int, ids ...string) (objs []ads.Map, total int64, err error) {
	var ctx = context.TODO()

	resp, _, err := g.SDKClient.Pages().Get(ctx, g.AccountID, &api.PagesGetOpts{
		Page:      optional.NewInt64(int64(page)),
		PageSize:  optional.NewInt64(int64(pageSize)),
		Filtering: filterIds("page_id", ids),
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

// PagesLoop
func (g *GdtAPI) PagesLoop(ids ...string) (objs []ads.Map, err error) {
	for page := 1; ; page++ {
		resp, total, err := g.Pages(page, 100, ids...)
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

// PagesAll
func (g *GdtAPI) AllPages(ids ...string) (all []ads.Map, err error) {
	if len(ids) == 0 {
		return g.PagesLoop()
	}

	for len(ids) > 0 {
		l := min(100, len(ids))
		objs, _ := g.PagesLoop(ids[:l]...)
		all = append(all, objs...)
		ids = ids[l:]
	}

	return
}

var ImageFields = []string{
	"image_id",
	"width",
	"height",
	"signature",
	"preview_url",
	"created_time",
	"last_modified_time",
	"source_type",
	// "product_catalog_id",
	// "product_outer_id",
	// "owner_account_id",
	"status",
	// "image_description",
	// "sample_aspect_ratio",
}

// Images
func (g *GdtAPI) Images(page, pageSize int, ids ...string) (objs []ads.Map, total int64, err error) {
	var ctx = context.TODO()

	resp, _, err := g.SDKClient.Images().Get(ctx, g.AccountID, &api.ImagesGetOpts{
		Page:      optional.NewInt64(int64(page)),
		PageSize:  optional.NewInt64(int64(pageSize)),
		Filtering: filterIds("image_id", ids),
		Fields:    optional.NewInterface(ImageFields),
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

// ImagesLoop
func (g *GdtAPI) ImagesLoop(ids ...string) (objs []ads.Map, err error) {
	for page := 1; ; page++ {
		resp, total, err := g.Images(page, 100, ids...)
		if err != nil {
			objs = append(objs, resp...)
			g.log.Errorw("images loop error", "error", err)
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
func (g *GdtAPI) AllImages(ids ...string) (all []ads.Map, err error) {
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

var VideoFields = []string{
	"video_id",
	"signature",
	"width",
	"height",
	"created_time",
	"last_modified_time",
	"source_type",
	"preview_url",
	"key_frame_image_url",
	"type",
	"status",
}

// Videos
func (g *GdtAPI) Videos(page, pageSize int, ids ...string) (objs []ads.Map, total int64, err error) {
	var ctx = context.TODO()

	resp, _, err := g.SDKClient.Videos().Get(ctx, g.AccountID, &api.VideosGetOpts{
		Page:      optional.NewInt64(int64(page)),
		PageSize:  optional.NewInt64(int64(pageSize)),
		Fields:    optional.NewInterface(VideoFields),
		Filtering: filterIds("video_id", ids),
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

// VideosLoop
func (g *GdtAPI) VideosLoop(ids ...string) (objs []ads.Map, err error) {
	for page := 1; ; page++ {
		resp, total, err := g.Videos(page, 100, ids...)
		if err != nil {
			objs = append(objs, resp...)
			g.log.Errorw("videos loop error", "error", err)
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
func (g *GdtAPI) AllVideos(ids ...string) (all []ads.Map, err error) {
	if len(ids) == 0 {
		return g.VideosLoop()
	}

	for len(ids) > 0 {
		l := min(100, len(ids))
		objs, _ := g.VideosLoop(ids[:l]...)
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
