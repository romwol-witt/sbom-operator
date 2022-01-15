package internal

import (
	"github.com/ckotzbauer/sbom-git-operator/internal/git"
	"github.com/ckotzbauer/sbom-git-operator/internal/kubernetes"
	"github.com/ckotzbauer/sbom-git-operator/internal/syft"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RunBackgroundService() {
	workingTree := viper.GetString("git-workingtree")

	gitAccount := git.New(viper.GetString("git-access-token"), viper.GetString("git-author-name"), viper.GetString("git-author-email"))
	gitAccount.Clone(viper.GetString("git-repository"), workingTree, viper.GetString("git-branch"))

	client := kubernetes.NewClient()
	pods := client.ListPods("monitoring")
	logrus.Debugf("Discovered %v pods", len(pods))
	digests := client.GetContainerDigests(pods)

	for _, d := range digests {
		syft.ExecuteSyft(d, workingTree)
	}

	gitAccount.CommitAll(workingTree, "Created new SBOMs")
}
