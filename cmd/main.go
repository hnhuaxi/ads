package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/fatih/structs"
	"github.com/hnhuaxi/ads"
	"github.com/hysios/x/utils"

	_ "github.com/hnhuaxi/ads/gdt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "accounts list"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	accessToken     = flag.String("access_token", "", "access token")
	debug           = flag.Bool("debug", false, "debug mode")
	provider        = flag.String("provider", "GDT", "provider")
	verbose         = flag.Bool("verbose", false, "verbose")
	onlyAdcreatives = flag.Bool("only_adcreatives", false, "only adcreatives")
	csvFile         = flag.String("output", "", "output to csv file")
)

var accounts arrayFlags

func main() {
	flag.Var(&accounts, "account", "ad account id.")
	flag.Parse()

	var (
		token  = utils.Default(*accessToken, os.Getenv("GDT_ACCESS_TOKEN"))
		output *csv.Writer
	)

	log := setLogger(*verbose)
	if len(accounts) == 0 {
		accounts = append(accounts, os.Getenv("GDT_ACCOUNT_ID"))
	}
	log.Infow("list accounts", "accounts", accounts)

	if *csvFile != "" {
		file, err := os.Create(*csvFile)
		if err != nil {
			log.Fatalf("create csv file error: %v", err)
		}
		defer file.Close()
		output = csv.NewWriter(file)
		writeHeader(output, &ads.Asset{})
		defer output.Flush()
	}

	for _, accId := range accounts {
		get, err := ads.Open(*provider, accId, token, *debug)
		if err != nil {
			log.Fatalf("open provider error: %v", err)
			return
		}

		if *onlyAdcreatives {
			get.OnlyAdcreatives(true)
		}

		assets, err := get.Assets()
		if err != nil {
			log.Fatalf("get assets error: %v", err)
			return
		}

		for _, asset := range assets {
			log.With("asset", asset).Info("asset")
		}

		if *csvFile != "" {
			writeTo(output, assets)
		}
	}
}

func setLogger(verbose bool) *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()
	logger.Info("info logging enabled")

	if verbose {
		atom.SetLevel(zap.DebugLevel)
	}
	zap.ReplaceGlobals(logger)
	return logger.Sugar()
}

// writeHeader
func writeHeader(w *csv.Writer, asset *ads.Asset) error {
	if w == nil {
		return errors.New("csv writer is nil")
	}

	s := structs.New(asset)
	rows := make([]string, 0, len(s.Names()))
	for _, name := range s.Names() {
		rows = append(rows, name)
	}

	w.Write(rows)
	return nil

}

// writeTo
func writeTo(w *csv.Writer, assets []*ads.Asset) error {
	if w == nil {
		return errors.New("csv writer is nil")
	}

	for _, asset := range assets {
		s := structs.New(asset)
		rows := make([]string, 0, len(s.Values()))
		for _, value := range s.Values() {
			rows = append(rows, fmt.Sprint(value))
		}

		w.Write(rows)
	}
	return nil
}

func init() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}
