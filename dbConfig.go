package main

import (
    "context"
    "fmt"
    "strings"

    "github.com/spf13/viper"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ssm"
)

const (
    baseRegion  = "AWS_REGION"
    baseRoot    = "AWS_ROOT"
    baseConfig  = "base_config"
    basePath    = "./app_data/"

    withEncrpytion = true
)

var ssmParams = []string{
    "database",
    "host",
    "port",
    "user",
    "password",

}

type awsSSM struct {
    *ssm.SSM
}

// ConfigString returns database connection string based on AWS_ROOT and remote SSM parameters
func ConfigString(ctx context.Context) (string, error) {
	err := loadBaseConfig()
	if err != nil {
		return "", err
    }
    
    awsRegion := viper.GetString(baseRegion)
    ssmRoot := viper.GetString(baseRoot)

    svc := newSSM(awsRegion)
    
    params, err := svc.getParams(ctx, withEncrpytion, ssmRoot, ssmParams)
    if err != nil {
        return "", err
    }

    configString := fmt.Sprintf(
        "database=%s host=%s port=%s user=%s password = %s",
		params["database"],
		params["host"],
		params["port"],
		params["user"],
		params["password"],
    )

    return configString, nil
}

// Pull AWS_ROOT and AWS_REGION from .env file
func loadBaseConfig() error {
    viper.SetConfigName(baseConfig)
    viper.AddConfigPath(basePath)

    err := viper.ReadInConfig()
    if err != nil {
        return err
	}
	
	return nil
}

// newSSM creates a new AWS connection returns a Simple Service Manager session
func newSSM(region string) *awsSSM {
    sess := session.New()

    return &awsSSM{ssm.New(sess, 
        &aws.Config{
            Region: aws.String(region),
    })}

}

// getParams returns map of key:value SSM Parameters as listed in paramsToGet along with any error fectching them
func (svc *awsSSM) getParams(ctx context.Context, encrpyted bool, root string, paramsToGet []string) (map[string]string, error) {
    params := make(map[string]string, len(paramsToGet))
    var paramsToGetPaths []*string

    for _, paramToGet := range paramsToGet {
        paramPath := root + paramToGet
        paramsToGetPaths = append(paramsToGetPaths, &paramPath)
    } 
    
    output, err := svc.GetParametersWithContext(ctx, 
            &ssm.GetParametersInput{
                Names: paramsToGetPaths,
                WithDecryption: aws.Bool(encrpyted),
    })

    for _, param := range output.Parameters {
        key := strings.TrimPrefix(*param.Name, root)
        val := *param.Value
        params[key] = val
    }
    
    return params, err
}