package generate_docs

import (
	"fmt"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//GetDirs takes a string argument and returns a slice of string of directories that containt terraform files
func GetDirs(s string) ([]string, error) {
	var mods []string
	var dirs []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dirs = append(dirs, path)
		return nil
	})
	if err != nil {
		return mods, fmt.Errorf("Unable to parse directory: %v", err)
	}
	for _, v := range dirs {
		check := tfconfig.IsModuleDir(v)
		if check {
			mods = append(mods, v)
		}
	}
	return mods, nil
}

//GetData takes a slice of sting and creates JSON objects from data retrieved such as variables, resources, modules, etc
func GetData(ls []string) (Stats, error) {
	var Final Stats
	for _, v := range ls {
		config, err := tfconfig.LoadModule(v)
		if err != nil {
			return Final, fmt.Errorf("Failed to load module: %v", v)
		}
		newVars := config.Variables
		newResources := config.ManagedResources
		newModules := config.ModuleCalls
		newOutputs := config.Outputs
		newData := config.DataResources
		newProviders := config.ProviderConfigs

		for _, z := range newVars {
			someVars := Var{
				VarName:                z.Name,
				VarType:                z.Type,
				VarDefault:             z.Default,
				VarDescription:         z.Description,
				VarRequired:            z.Required,
				VarSensitive:           z.Sensitive,
				SourcePositionFileName: fmt.Sprintf("./%s", z.Pos.Filename),
				SourcePositionLine:     strconv.Itoa(z.Pos.Line),
			}

			Final.Vars = append(Final.Vars, someVars)
		}

		for _, x := range newResources {
			someResources := Resource{
				Mode:                   x.Mode.String(),
				Type:                   x.Type,
				Name:                   x.Name,
				ProviderName:           x.Provider.Name,
				ProviderAlias:          x.Provider.Alias,
				SourcePositionFileName: fmt.Sprintf("./%s", x.Pos.Filename),
				SourcePositionLine:     strconv.Itoa(x.Pos.Line),
			}
			Final.Resources = append(Final.Resources, someResources)
		}

		for _, y := range newModules {
			someModules := Module{
				Name:                   y.Name,
				ModSource:              y.Source,
				Version:                y.Version,
				SourcePositionFileName: y.Pos.Filename,
				SourcePositionLine:     strconv.Itoa(y.Pos.Line),
			}
			Final.Modules = append(Final.Modules, someModules)
		}

		for _, w := range newOutputs {
			someOutputs := Output{
				Name:                   w.Name,
				Description:            w.Description,
				Sensitive:              w.Sensitive,
				SourcePositionFileName: w.Pos.Filename,
				SourcePositionLine:     strconv.Itoa(w.Pos.Line),
			}
			Final.Outputs = append(Final.Outputs, someOutputs)
		}

		for _, u := range newData {
			someData := Data{
				DataType:               u.Type,
				Name:                   u.Name,
				ProviderName:           u.Provider.Name,
				ProviderAlias:          u.Provider.Alias,
				SourcePositionFileName: u.Pos.Filename,
				SourcePositionLine:     strconv.Itoa(u.Pos.Line),
			}
			Final.Datas = append(Final.Datas, someData)
		}

		for _, u := range newProviders {
			someProvider := Provider{
				Name:  u.Name,
				Alias: u.Alias,
			}
			Final.Providers = append(Final.Providers, someProvider)
		}
	}
	return Final, nil
}

//GetDirData iterates through directories and returns data about each directory
func GetDirData(ls []string) (RepoInfo, error){
	var dirs RepoInfo
	for _, v := range ls {
		read, err := ioutil.ReadDir(v)
		if err != nil {
			return dirs, fmt.Errorf(err.Error())
		}
		for _, z := range read {
			var theDir Dirs
			if z.IsDir() {
				theDir.Name = z.Name()
				theDir.ModificationTime = z.ModTime()
				if tfconfig.IsModuleDir(z.Name()) {
					theDir.IsTerraDir = true
				} else {
					theDir.IsTerraDir = false
				}
				dirs.Directories = append(dirs.Directories, theDir)
			}
		}
	}
	return dirs, nil
}

//GetFileInfo iterates through files and returns data about each file
func GetFileInfo(ls []string) (RepoInfo, error) {
	var files RepoInfo
	for _, v := range ls {
		grep, err := ioutil.ReadDir(v)
		if err != nil {
			return files, fmt.Errorf(err.Error())
		}
		
		for _, x := range grep {
			if !x.IsDir() {
				var file File
				file.Name = x.Name()
				file.ModificationTime = x.ModTime()
				if strings.Contains(x.Name(), ".tf") {
					file.IsTfFile = true
				} else {
					file.IsTfFile = false
				}
				files.Files = append(files.Files, file)
			}
		}
	}
	return files, nil 
}
