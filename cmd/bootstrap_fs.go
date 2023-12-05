package cmd

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/spf13/cobra"

	"nomad-gitops-operator/pkg/reconcile"
)

type fsFlags struct {
	base_dir string
	path     string
	watch    bool
}

var fsArgs fsFlags

func init() {
	bootstrapCmd.AddCommand(bootstrapFsCmd)
	bootstrapFsCmd.Flags().StringVar(&fsArgs.base_dir, "base-dir", "./", "Path to the base directory")
	bootstrapFsCmd.Flags().StringVar(&fsArgs.path, "path", "**/*.nomad", "glob pattern relative to the base-dir")
	bootstrapFsCmd.Flags().BoolVar(&fsArgs.watch, "watch", false, "Enable watch mode")
}

var bootstrapFsCmd = &cobra.Command{
	Use:   "fs [path]",
	Short: "Bootstrap Nomad using a local path",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return reconcile.Run(reconcile.ReconcileOptions{
			Path:  fsArgs.path,
			Watch: fsArgs.watch,
			Fs: func() (billy.Filesystem, error) {
				fs := osfs.New(fsArgs.base_dir)
				return fs, nil
			}})
	},
}
