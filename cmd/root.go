package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var rootDirectory string
var vendorFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "classifier",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		mapping, err := getMapping(vendorFile)
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		if err := processPDFs(mapping); err != nil {
			fmt.Println("Error Processing: ", err)
			return
		}
	},
}

func processPDFs(mapping Vendors) error {
	files, err := ioutil.ReadDir(filepath.Join(rootDirectory, "Unfiled"))
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), ".pdf") {

				fmt.Println("Processing ", f.Name())
			VendorLoop:
				for _, v := range mapping.Vendor {
					text, err := getContents(f)
					if err != nil {
						fmt.Println("Error reading pdf:", err)
						continue
					}
					if v.KeywordMatch(text) {
						fmt.Println("Doing Regex")
						year := strconv.Itoa(time.Now().Year())
						month := strconv.Itoa(int(time.Now().Month()))

						m := Match{
							Vendor: v,
							File:   f,
							Year:   year,
							Month:  month,
						}
						err := moveFile(m)
						if err != nil {
							fmt.Println("Error moving file:", err)
						}
						break VendorLoop
					}
				}
				fmt.Println("After vendor loop")
			}
		}

	}
	return nil
}
func moveFile(m Match) error {
	fmt.Println("Moving ", m.File.Name(), m.Vendor.Directory)
	err := os.MkdirAll(filepath.Join(rootDirectory, "Filed", m.Vendor.Directory, m.Year, m.Month), 0755)
	if err != nil {
		return err
	}
	err = os.Rename(filepath.Join(rootDirectory, "Unfiled", m.File.Name()), filepath.Join(rootDirectory, "Filed", m.Vendor.Directory, m.Year, m.Month, m.File.Name()))
	return err
}

func getContents(f os.FileInfo) (string, error) {
	cmd := exec.Command("pdf2txt.py", filepath.Join(rootDirectory, "Unfiled", f.Name()))
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdoutStderr))
	}
	return string(stdoutStderr), err
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.classifier.yaml)")
	rootCmd.PersistentFlags().StringVar(&vendorFile, "vendor", "/Users/bketelsen/classifier.toml", "vendor list (default is $HOME/classifier.toml")
	rootCmd.PersistentFlags().StringVarP(&rootDirectory, "directory", "d", "/Users/bketelsen/Documents", "/Users/bketelsen/Documents")
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".classifier" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".classifier")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
