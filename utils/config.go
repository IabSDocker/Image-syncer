// package utils

// import (
//     "io/ioutil"
//     "gopkg.in/yaml.v2"
// )

// type Config struct {
//     SourceRepo string   `yaml:"source_repo"`
//     TargetRepo string   `yaml:"target_repo"`
//     Images     []string `yaml:"images"`
// }

// func LoadConfig(path string) (*Config, error) {
//     var config Config
//     data, err := ioutil.ReadFile(path)
//     if err != nil {
//         return nil, err
//     }
//     err = yaml.Unmarshal(data, &config)
//     if err != nil {
//         return nil, err
//     }
//     return &config, nil
// }

package utils

import (
    "gopkg.in/yaml.v2"
    "os"
)

type Config struct {
    SourceRepo string   `yaml:"source_repo"`
    TargetRepo string   `yaml:"target_repo"`
    Namespace  string   `yaml:"namespace"`
    Images     []string `yaml:"images"`
}

func LoadConfig(filePath string) (*Config, error) {
    var config Config
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    decoder := yaml.NewDecoder(file)
    err = decoder.Decode(&config)
    if err != nil {
        return nil, err
    }
    
    return &config, nil
}

