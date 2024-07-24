package gdt

import (
	"strconv"

	"github.com/hnhuaxi/ads"
	v2 "github.com/hnhuaxi/ads/gdt/v2"
	v3 "github.com/hnhuaxi/ads/gdt/v3"
	"github.com/stretchr/objx"
	"go.uber.org/zap"
)

type Config struct {
	OnlyAdcreatives bool
}

type GdtAdcreatives struct {
	AccountID       int64
	v2              *v2.GdtAPI
	v3              *v3.GdtV3API
	Config          Config
	adcreateivesFns []ads.AdcreativeMatchFunc
	log             *zap.SugaredLogger
}

func NewAdcreatives(accountId string, accessToken string, debug bool) (*GdtAdcreatives, error) {
	id, err := strconv.ParseInt(accountId, 10, 64)
	if err != nil {
		return nil, err
	}
	advs := &GdtAdcreatives{
		AccountID: id,
		v2:        v2.NewGdtAPI(accountId, accessToken, debug),
		v3:        v3.NewGdtAPI(accountId, accessToken, debug),
		log:       zap.S(),
	}

	advs.SetAdcreativesFunc(func(adcreative ads.Map) (process bool) {
		return !adcreative.Get("is_deleted").Bool()
	})

	return advs, nil
}

// Assets ...
func (g *GdtAdcreatives) Assets() (assets []*ads.Asset, err error) {
	r, total, err := g.v2.Adcreatives(1, 100)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		goto V3
	}
	{
		log := g.log.With("version", "v2")
		adcreatives := g.processAdcreatives(r)
		g.printJson(log, "adcreatives", adcreatives)

		var (
			needPages  bool
			pageTypes  = make(ads.Set[string])
			pageIds    = make(ads.Set[string])
			videoIds   = make(ads.Set[string])
			imageIds   = make(ads.Set[string])
			brandIds   = make(ads.Set[string])
			needVideos bool
			needBrands bool
			needImages bool
			// videosIds
			needXJPages  bool
			needProfiles bool
		)

		_ = needVideos
		_ = needXJPages
		_ = needProfiles
		_ = needBrands
		_ = videoIds

		for _, adcr := range adcreatives {
			pageType := adcr.Get("page_type").String()
			// switch pageType {
			// case "PAGE_TYPE_DEFAULT":
			pageSpec := adcr.Get("page_spec").ObjxMap()
			if pageSpec.Get("page_url").String() != "" {
				assets = append(assets, &ads.Asset{
					AccountID: strconv.FormatInt(g.AccountID, 10),
					Name:      adcr.Get("adcreative_name").Str(),
					AssetID:   adcr.Get("adcreative_id").Str(),
					PageType:  ads.PTPageUrl,
					SubType:   pageType,
					Url:       pageSpec.Get("page_url").String(),
				})
			} else if pageSpec.Get("page_id").Int() > 0 {
				pageTypes.Add("DEFAULT_PAGES")
				needPages = true
				pageIds.Add(itoa(pageSpec.Get("page_id").Int()))
			} else {
				pageTypes.Add("DEFAULT_PAGES")
				needPages = true
			}

			adcreativeElements := adcr.Get("adcreative_elements").ObjxMap()
			if !adcreativeElements.Value().IsNil() {
				if !adcreativeElements.Get("brand_component_options").IsNil() {
					adcreativeElements.Get("brand_component_options").EachObjxMap(func(i int, m objx.Map) bool {
						imageId := m.Get("value.brand_img.image_id").String()
						brandIds.Add(imageId)
						needBrands = true
						return true
					})
				}

				if !adcreativeElements.Get("image_component_options").IsNil() {
					adcreativeElements.Get("image_component_options").EachObjxMap(func(i int, m objx.Map) bool {
						imageId := m.Get("value.image_id").String()
						imageIds.Add(imageId)
						needImages = true
						return true
					})
				}

				if !adcreativeElements.Get("image3_component_options").IsNil() {
					adcreativeElements.Get("image3_component_options").EachObjxMap(func(i int, m objx.Map) bool {
						imageId := m.Get("value.image_id").String()
						imageIds.Add(imageId)
						needImages = true
						return true
					})
				}

				if !adcreativeElements.Get("video2_component_options").IsNil() {
					adcreativeElements.Get("video2_component_options").EachObjxMap(func(i int, m objx.Map) bool {
						videoId := m.Get("value.video_id").String()
						videoUrl := m.Get("value.video_url").String()
						coverImgId := m.Get("value.cover_image.image_id").String()
						if videoUrl != "" {
							assets = append(assets, &ads.Asset{
								AccountID: strconv.FormatInt(g.AccountID, 10),
								Name:      adcr.Get("adcreative_name").Str(),
								AssetID:   videoId,
								PageType:  ads.PTVideo,
								SubType:   "VIDEO2",
								Url:       videoUrl,
							})
						} else {
							videoIds.Add(videoId)
						}

						if coverImgId != "" {
							imageIds.Add(coverImgId)
							needImages = true
						}

						needVideos = true
						return true
					})
				}

			}
			// case "PAGE_TYPE_MINI_PROGRAM_WECHAT":
			// }
		}

		if needPages {
			var pages []ads.Map
			if g.Config.OnlyAdcreatives {
				pages, err = g.v2.AllPages(pageIds.Slice()...)
				if err != nil {
					return nil, err
				}
			} else {
				pages, err = g.v2.AllPages()
				if err != nil {
					return nil, err
				}
			}

			g.printJson(log, "pages", pages)

			for _, page := range pages {
				pageType := page.Get("page_type").String()
				switch pageType {
				case "PAGE_TYPE_DEFAULT":
					assets = append(assets, &ads.Asset{
						AccountID: strconv.FormatInt(g.AccountID, 10),
						Name:      page.Get("page_name").Str(),
						AssetID:   page.Get("page_id").Str(),
						PageType:  ads.PTPageUrl,
						SubType:   pageType,
						Url:       page.Get("page_url").String(),
					})
				}
			}
		}

		if needBrands {

		}

		if needImages {
			var images []ads.Map
			if g.Config.OnlyAdcreatives {
				images, err = g.v2.AllImages(imageIds.Slice()...)
				if err != nil {
					return nil, err
				}
			} else {
				images, err = g.v2.AllImages()
				if err != nil {
					return nil, err
				}
			}

			g.printJson(log, "images", images)

			for _, image := range images {
				assets = append(assets, &ads.Asset{
					AccountID: strconv.FormatInt(g.AccountID, 10),
					Name:      image.Get("image_name").Str(),
					AssetID:   image.Get("image_id").Str(),
					PageType:  ads.PTImage,
					Url:       image.Get("preview_url").String(),
					Signature: image.Get("signature").String(),
				})
			}
		}

		if needVideos {
			var videos []ads.Map
			if g.Config.OnlyAdcreatives {
				videos, err = g.v2.AllVideos(videoIds.Slice()...)
				if err != nil {
					return nil, err
				}
			} else {
				videos, err = g.v2.AllVideos()
				if err != nil {
					return nil, err
				}
			}
			g.printJson(log, "videos", videos)

			for _, video := range videos {
				videoType := video.Get("type").String()

				assets = append(assets, &ads.Asset{
					AccountID: strconv.FormatInt(g.AccountID, 10),
					Name:      video.Get("description").Str(),
					AssetID:   video.Get("video_id").Str(),
					PageType:  ads.PTVideo,
					SubType:   videoType,
					Url:       video.Get("preview_url").String(),
					Signature: video.Get("signature").String(),
				})

				assets = append(assets, &ads.Asset{
					AccountID: strconv.FormatInt(g.AccountID, 10),
					Name:      video.Get("description").Str(),
					AssetID:   video.Get("video_id").Str(),
					PageType:  ads.PTImage,
					Url:       video.Get("key_frame_image_url").String(),
				})
			}
		}
		return
	}
V3:
	{
		log := g.log.With("version", "v3")
		r, err = g.v3.AllAdcreatives()
		if err != nil {
			return nil, err
		}

		adcreatives := g.processAdcreatives(r)
		_ = adcreatives

		var (
			needPages  bool
			pageTypes  = make(ads.Set[string])
			needVideos bool
			videosIds  = make(ads.Set[string])
			imagesIds  = make(ads.Set[string])
			pagesIds   = make(ads.Set[string])
			// videosIds
			needImages   bool
			needXJPages  bool
			needProfiles bool
		)

		_ = needVideos
		_ = needXJPages
		_ = needProfiles

		for _, adcr := range adcreatives {
			creativeType := adcr.Get("dynamic_creative_type").String()
			log.Debugf("creativeType: %s", creativeType)

			jumpInfo := adcr.Get("creative_components.main_jump_info")      // 跳转信息
			brand := adcr.Get("creative_components.brand")                  // 品牌信息
			video := adcr.Get("creative_components.video")                  // 视频
			image := adcr.Get("creative_components.image")                  // 图片
			wechatChannel := adcr.Get("creative_components.wechat_channel") // 微信视频号

			if !jumpInfo.IsNil() {
				jumpInfo.EachObjxMap(func(i int, m objx.Map) bool {
					log.Debugf("page_type: %s", m.Get("value.page_type").String())
					pageType := m.Get("value.page_type").String()
					pageSpec := m.Get("value.page_spec")

					if !pageSpec.IsNil() {
						if pageSpec.ObjxMap().Get("wechat_canvas_spec.page_id").String() != "" {
							pagesIds.Add(pageSpec.ObjxMap().Get("wechat_canvas_spec.page_id").String())
							needPages = true
						} else if pageSpec.ObjxMap().Get("h5_spec.page_url").String() != "" {
							assets = append(assets, &ads.Asset{
								AccountID: strconv.FormatInt(g.AccountID, 10),
								Name:      adcr.Get("adcreative_name").Str(),
								AssetID:   adcr.Get("adcreative_id").Str(),
								PageType:  ads.PTPageUrl,
								SubType:   pageType,
								Url:       pageSpec.ObjxMap().Get("h5_spec.page_url").String(),
							})
						}
					}
					switch pageType {
					case "PAGE_TYPE_WECHAT_CANVAS":
						pageTypes.Add("WECHAT_PAGES")
						needPages = true
					}
					return true
				})
			}

			if !brand.IsNil() {

				brand.EachObjxMap(func(i int, m objx.Map) bool {
					log.Debugf("page_type: %s", m.Get("value.page_type").String())
					pageType := m.Get("value.page_type").String()
					switch pageType {
					case "PAGE_TYPE_WECHAT_CHANNELS_PROFILE":
						// addPage("WECHAT_PAGES")
						needProfiles = true
					}
					return true
				})
			}

			if !video.IsNil() {
				video.EachObjxMap(func(i int, m objx.Map) bool {
					videoId := m.Get("value.video_id").String()
					log.Debugf("video_id: %s", videoId)
					needVideos = true
					videosIds.Add(videoId)
					return true
				})
			}

			if !image.IsNil() {
				image.EachObjxMap(func(i int, m objx.Map) bool {
					imageId := m.Get("value.image_id").String()
					log.Debugf("image_id: %s", imageId)
					needImages = true
					imagesIds.Add(imageId)
					return true
				})
			}

			_ = wechatChannel

		}
		g.printJson(log, "adcreatives", adcreatives)

		if needPages {
			for _, pageType := range pageTypes.Slice() {
				pages, err := g.v3.AllPages(pageType)
				if err != nil {
					return nil, err
				}

				for _, page := range pages {
					pageType := page.Get("page_type").String()
					assets = append(assets, &ads.Asset{
						AccountID: strconv.FormatInt(g.AccountID, 10),
						Name:      page.Get("page_name").Str(),
						AssetID:   page.Get("page_id").Str(),
						PageType:  ads.PTPageUrl,
						SubType:   pageType,
						Url:       page.Get("preview_url").String(),
					})
				}
				g.printJson(log, "pages", pages)
			}
		}

		if needVideos {
			var videos []ads.Map
			if g.Config.OnlyAdcreatives {
				videos, err = g.v3.AllVideos(videosIds.Slice()...)
				if err != nil {
					return nil, err
				}
			} else {
				videos, err = g.v3.AllVideos()
				if err != nil {
					return nil, err
				}
			}

			g.printJson(log, "videos", videos)
			for _, video := range videos {
				videoType := video.Get("type").String()

				assets = append(assets, &ads.Asset{
					AccountID: strconv.FormatInt(g.AccountID, 10),
					Name:      video.Get("description").Str(),
					AssetID:   video.Get("video_id").Str(),
					PageType:  ads.PTVideo,
					SubType:   videoType,
					Url:       video.Get("preview_url").String(),
				})

				assets = append(assets, &ads.Asset{
					AccountID: strconv.FormatInt(g.AccountID, 10),
					Name:      video.Get("description").Str(),
					AssetID:   video.Get("video_id").Str(),
					PageType:  ads.PTImage,
					Url:       video.Get("key_frame_image_url").String(),
				})
			}
		}

		if needImages {
			var images []ads.Map
			if g.Config.OnlyAdcreatives {
				images, err = g.v3.AllImages(imagesIds.Slice()...)
				if err != nil {
					return nil, err
				}
			} else {
				images, err = g.v3.AllImages()
				if err != nil {
					return nil, err
				}
			}

			g.printJson(log, "images", images)
			for _, image := range images {
				assets = append(assets, &ads.Asset{
					AccountID: strconv.FormatInt(g.AccountID, 10),
					Name:      image.Get("image_name").Str(),
					AssetID:   image.Get("image_id").Str(),
					PageType:  ads.PTImage,
					Url:       image.Get("preview_url").String(),
				})
			}
		}
		// printJson(ads)
		return
	}

}

// SetAdcreativesFunc ...
func (g *GdtAdcreatives) SetAdcreativesFunc(match ads.AdcreativeMatchFunc) {
	g.adcreateivesFns = append(g.adcreateivesFns, match)
}

// OnlyAdcreatives ...
func (g *GdtAdcreatives) OnlyAdcreatives(on bool) {
	g.Config.OnlyAdcreatives = on
}

// processAdcreatives ...
func (g *GdtAdcreatives) processAdcreatives(adcreatives []ads.Map) (processes []ads.Map) {
	for _, adcreative := range adcreatives {
		for _, fn := range g.adcreateivesFns {
			if fn(adcreative) {
				processes = append(processes, adcreative)
			}
		}
	}

	return processes
}

// printJson
func (g *GdtAdcreatives) printJson(log *zap.SugaredLogger, key string, v interface{}) {
	log.With(key, v).Debug("json")
}

var _ ads.GetAdcreatives = (*GdtAdcreatives)(nil)

func init() {
	ads.RegisterProvider("GDT", func(accountId string, accessToken string, debug bool) (ads.GetAdcreatives, error) {
		return NewAdcreatives(accountId, accessToken, debug)
	})
}

func itoa(i int) string {
	return strconv.Itoa(i)
}
