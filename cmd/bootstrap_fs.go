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
	var_path string
	watch    bool
	delete   bool
}

var fsArgs fsFlags

func init() {
	bootstrapCmd.AddCommand(bootstrapFsCmd)
	bootstrapFsCmd.Flags().StringVar(&fsArgs.base_dir, "base-dir", "./", "Path to the base directory")
	bootstrapFsCmd.Flags().StringVar(&fsArgs.path, "path", "**/*.nomad", "glob pattern relative to the base-dir")
	bootstrapFsCmd.Flags().StringVar(&fsArgs.var_path, "var-path", "**/*.vars.yml", "var glob pattern relative to the base-dir")
	bootstrapFsCmd.Flags().BoolVar(&fsArgs.watch, "watch", false, "Enable watch mode")
	bootstrapFsCmd.Flags().BoolVar(&fsArgs.delete, "delete", false, "Enable delete missing jobs")
}

var bootstrapFsCmd = &cobra.Command{
	Use:   "fs [path]",
	Short: "Bootstrap Nomad using a local path",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return reconcile.Run(reconcile.ReconcileOptions{
			Path:    fsArgs.path,
			VarPath: fsArgs.var_path,
			Watch:   fsArgs.watch,
			Delete:  fsArgs.delete,
			Fs: func() (billy.Filesystem, error) {
				fs := osfs.New(fsArgs.base_dir)
				return fs, nil
			}})
	},
}
