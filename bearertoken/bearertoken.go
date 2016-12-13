package bearertoken

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type key struct {
	Key string `json:"key"`
}

func GetToken() (key, error) {
	svc := s3.New(session.New(), &aws.Config{Region: aws.String("us-west-2")})

	response, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("elasticbeanstalk-us-west-2-194925301021"),
		Key:    aws.String("bearer-token.json"),
	})

	if err != nil {
		return key{}, err
	}

	defer response.Body.Close()

	result := key{}
	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)
	if err != nil {
		return key{}, err
	}

	return result, nil
}
