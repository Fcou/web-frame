package command

import "github.com/Fcou/web-frame/framework/cobra"

// AddKernelCommands will add all command/* to root command
func AddKernelCommands(root *cobra.Command) {
	//root.AddCommand(DemoCommand)

	//root.AddCommand(initEnvCommand())
	//root.AddCommand(deployCommand)
	//
	// cron
	root.AddCommand(initCronCommand())
	//// cmd
	//cmdCommand.AddCommand(cmdListCommand)
	//cmdCommand.AddCommand(cmdCreateCommand)
	//root.AddCommand(cmdCommand)

	// build
	root.AddCommand(initBuildCommand())
	//
	// app
	root.AddCommand(initAppCommand())
	// go build
	root.AddCommand(goCommand)
	// npm build
	root.AddCommand(npmCommand)
	//
	//// dev
	//root.AddCommand(initDevCommand())
	//
	//// middleware
	//middlewareCommand.AddCommand(middlewareAllCommand)
	//middlewareCommand.AddCommand(middlewareAddCommand)
	//middlewareCommand.AddCommand(middlewareRemoveCommand)
	//root.AddCommand(middlewareCommand)
	//
	//// swagger
	//swagger.IndexCommand.AddCommand(swagger.InitServeCommand())
	//swagger.IndexCommand.AddCommand(swagger.GenCommand)
	//root.AddCommand(swagger.IndexCommand)
	//
	//// provider
	//providerCommand.AddCommand(providerListCommand)
	//providerCommand.AddCommand(providerCreateCommand)
	//root.AddCommand(providerCommand)
	//
	//// new
	//root.AddCommand(initNewCommand())
}
