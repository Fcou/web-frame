package command

import "github.com/Fcou/web-frame/framework/cobra"

// AddKernelCommands will add all command/* to root command
func AddKernelCommands(root *cobra.Command) {
	// cron 定时
	root.AddCommand(initCronCommand())
	// build 编译前后端
	root.AddCommand(initBuildCommand())
	// app 业务
	root.AddCommand(initAppCommand())
	// go build 编译main.go
	root.AddCommand(goCommand)
	// npm build 编译
	root.AddCommand(npmCommand)
	// env 获取环境变量
	root.AddCommand(initEnvCommand())

	//root.AddCommand(DemoCommand)

	//root.AddCommand(deployCommand)

}
