package cmd

import (
	"dagger.io/go/cmd/dagger/logger"
	"dagger.io/go/dagger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new route",
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Fix Viper bug for duplicate flags:
		// https://github.com/spf13/viper/issues/233
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			panic(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		lg := logger.New()
		ctx := lg.WithContext(cmd.Context())
		store := dagger.DefaultStore()

		upRouteFlag, err := cmd.Flags().GetBool("up")
		if err != nil {
			lg.Fatal().Err(err).Str("flag", "up").Msg("unable to resolve flag")
		}

		routeName := getRouteName(lg, cmd)

		// TODO: Implement options: --layout-*, --setup
		route, err := store.CreateRoute(ctx, routeName, nil)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to create route")
		}
		lg.Info().Str("route-id", route.ID()).Str("route-name", routeName).Msg("created route")

		if upRouteFlag {
			routeUp(ctx, lg, route)
		}
	},
}

func init() {
	newCmd.Flags().StringP("name", "n", "", "Specify a route name")
	newCmd.Flags().BoolP("up", "u", false, "Bring the route online")

	newCmd.Flags().String("layout-dir", "", "Load layout from a local directory")
	newCmd.Flags().String("layout-git", "", "Load layout from a git repository")
	newCmd.Flags().String("layout-package", "", "Load layout from a cue package")
	newCmd.Flags().String("layout-file", "", "Load layout from a cue or json file")

	newCmd.Flags().String("setup", "auto", "Specify whether to prompt user for initial setup (no|yes|auto)")

	if err := viper.BindPFlags(newCmd.Flags()); err != nil {
		panic(err)
	}
}
