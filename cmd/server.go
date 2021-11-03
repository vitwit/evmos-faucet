package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gsk967/cosmos-faucet/config"
	"github.com/gsk967/cosmos-faucet/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"html/template"
	"io"
	"math/big"
	"net/http"
)

func StartServer(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the ethermint faucet server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			return startServer(clientCtx, cfg, cmd.Flags())
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	err := cmd.MarkFlagRequired(flags.FlagFrom)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	err := t.templates.ExecuteTemplate(w, name, data)
	return err
}

func startServer(ctx client.Context, cfg *config.Config, flagSet *pflag.FlagSet) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = renderer
	tokens := sdk.NewInt(cfg.Faucet.Amount).Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(cfg.Faucet.Decimals)), nil)))
	maxTokens := sdk.NewInt(cfg.Faucet.MaxTokens).Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(cfg.Faucet.Decimals)), nil)))

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"tokens":    fmt.Sprintf("%s%s", tokens, cfg.Faucet.Denom),
			"maxTokens": fmt.Sprintf("%s%s", maxTokens, cfg.Faucet.Denom),
		})
	})

	e.POST("/", func(c echo.Context) error {
		val := c.FormValue("accountAddress")
		toAddress, err := utils.ValidateAccountAddress(val)
		if err != nil {
			return c.Render(http.StatusOK, "index.html", map[string]interface{}{
				"errorMessage": fmt.Sprintf("Error while account checking... \n %s", err.Error()),
				"tokens":       fmt.Sprintf("%s%s", tokens, cfg.Faucet.Denom),
				"maxTokens":    fmt.Sprintf("%s%s", maxTokens, cfg.Faucet.Denom),
			})
		}
		e.Logger.Printf(fmt.Sprintf("toAddress %s", toAddress))
		err = utils.GetTokens(ctx, cfg, flagSet, toAddress)
		if err != nil {
			return c.Render(http.StatusOK, "index.html", map[string]interface{}{
				"errorMessage": fmt.Sprintf("Error while submit the tx... \n %s", err.Error()),
				"tokens":       fmt.Sprintf("%s%s", tokens, cfg.Faucet.Denom),
				"maxTokens":    fmt.Sprintf("%s%s", maxTokens, cfg.Faucet.Denom),
			})
		} else {
			return c.Render(http.StatusOK, "index.html", map[string]interface{}{
				"successMessage": fmt.Sprintf("%s%s tokens send to your account %s", tokens, cfg.Faucet.Denom, val),
				"tokens":         fmt.Sprintf("%s%s", tokens, cfg.Faucet.Denom),
				"maxTokens":      fmt.Sprintf("%s%s", maxTokens, cfg.Faucet.Denom),
			})
		}
	})

	return e.Start(fmt.Sprintf("0.0.0.0:%d", cfg.UI.Port))
}
